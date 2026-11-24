// Command build fetches all sources in sources.yaml and writes Shadowrocket
// rule lists, client configs and rule-sets into dist/, binary geo files into
// out/, and the generated docs pages into docs/.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/split-route/routing-rules/internal/emit"
	"github.com/split-route/routing-rules/internal/fetch"
	"github.com/split-route/routing-rules/internal/merge"
	"github.com/split-route/routing-rules/internal/parse"
)

func main() {
	cfgPath := flag.String("config", "sources.yaml", "source registry")
	dist := flag.String("dist", "dist", "output directory")
	cache := flag.String("cache", ".cache", "HTTP + last-known-good cache directory")
	allowPath := flag.String("allowlist", "data/allowlist.txt", "domains never to REJECT")
	canaryPath := flag.String("canaries", "data/canaries.txt", "known answers checked before publishing")
	out := flag.String("out", "out", "directory for binary release assets (geosite.dat, geoip.dat); not committed")
	docs := flag.String("docs", "docs", "docs root (docs/<lang>/) for the generated URL and sources pages and the checker's lists")
	flag.Parse()
	log.SetFlags(0)
	if err := run(*cfgPath, *dist, *cache, *allowPath, *out, *docs, *canaryPath); err != nil {
		log.Printf("::error::%v", err)
		os.Exit(1)
	}
}

// run builds every output; nothing is written if a safety check fails.
func run(cfgPath, dist, cacheDir, allowPath, outDir, docsDir, canaryPath string) error {
	raw, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	var cfg Config
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return fmt.Errorf("%s: %w", cfgPath, err)
	}
	cats := map[string]Category{}
	for _, c := range cfg.Categories {
		cats[c.Name] = c
	}
	for _, s := range cfg.Sources {
		if _, ok := cats[s.Category]; !ok {
			return fmt.Errorf("source %s: unknown category %q", s.Name, s.Category)
		}
		if parse.Formats[s.Format] == nil {
			return fmt.Errorf("source %s: unknown format %q", s.Name, s.Format)
		}
	}
	allow, err := readAllowlist(allowPath)
	if err != nil {
		return err
	}

	st := cfg.Settings
	f := &fetch.Fetcher{
		Client:   &http.Client{Timeout: st.Timeout, Transport: transport()},
		CacheDir: cacheDir,
		UA:       st.UserAgent,
		MaxBytes: st.MaxBytes,
		Retries:  st.Retries,
		Backoff:  2 * time.Second,
	}
	outcomes := fetchAll(f, cfg, cats, cacheDir)

	// Per-category merge.
	type built struct {
		cat     Category
		domains []parse.Domain
		cidrs   []netip.Prefix
		body    []byte
	}
	byCat := map[string]*built{}
	for _, c := range cfg.Categories {
		if c.Enabled {
			byCat[c.Name] = &built{cat: c}
		}
	}
	var stale []string
	for _, o := range outcomes {
		b := byCat[o.src.Category]
		if b == nil {
			continue
		}
		b.domains = append(b.domains, o.res.Domains...)
		b.cidrs = append(b.cidrs, o.res.CIDRs...)
		if o.status == "stale" {
			stale = append(stale, o.src.Name)
		}
	}
	for _, b := range byCat {
		b.domains = merge.Domains(b.domains)
		var broad []netip.Prefix
		b.cidrs = slices.DeleteFunc(merge.CIDRs(b.cidrs), func(p netip.Prefix) bool {
			// The private list legitimately holds broad blocks (224.0.0.0/4, fc00::/7).
			if b.cat.Name != "private" && p.Bits() < minBits[p.Addr().Is4()] {
				broad = append(broad, p)
				return true
			}
			return false
		})
		for _, p := range broad {
			log.Printf("::warning::%s: dropped over-broad prefix %s", b.cat.Name, p)
		}
	}

	// REJECT lists never block allowlisted domains. Rules that also match
	// ir-domains are kept: the configs check REJECT lists first, so they win.
	var warns []string
	for _, b := range byCat {
		if b.cat.Action != "reject" {
			continue
		}
		var rm []parse.Domain
		b.domains, rm = merge.RemoveAllowed(b.domains, allow)
		for _, d := range rm {
			warns = append(warns, fmt.Sprintf("%s: allowlist blocked %s", b.cat.Name, ruleString(d)))
		}
	}
	slices.Sort(warns)
	for _, w := range warns {
		log.Printf("::warning::ALLOWLIST HIT — %s (removed from REJECT list)", w)
	}

	// Render. An enabled list with no rules means every source failed or a
	// parser broke; refuse to publish it. Big swings are allowed.
	var fatal []string
	for _, c := range cfg.Categories {
		b := byCat[c.Name]
		if b == nil {
			continue
		}
		b.body = emit.List(c.Name, b.domains, b.cidrs)
		n := len(b.domains) + len(b.cidrs)
		if n == 0 {
			fatal = append(fatal, c.Name+": list is empty")
		}
		if st.ListCapWarn > 0 && n > st.ListCapWarn {
			log.Printf("::warning::%s: %d rules exceeds list_cap_warn %d (iOS network extension memory)", c.Name, n, st.ListCapWarn)
		}
	}

	printSummary(outcomes, cfg.Categories, func(name string) (int, int) {
		b := byCat[name]
		if b == nil {
			return 0, 0
		}
		return len(b.domains) + len(b.cidrs), len(b.body)
	})
	if len(stale) > 0 {
		log.Printf("::warning::using last-known-good data for: %s", strings.Join(stale, ", "))
	}
	sets := map[string]catRules{}
	for name, b := range byCat {
		sets[name] = catRules{b.cat, b.domains, b.cidrs}
	}
	canaries, bad := checkCanaries(canaryPath, cfg.Categories, sets)
	fatal = append(fatal, bad...)
	if len(fatal) > 0 {
		return fmt.Errorf("build refused, nothing written:\n  %s", strings.Join(fatal, "\n  "))
	}

	// Render every output into one map, then make dist/ match it exactly
	// (stale files from disabled categories are removed).
	if st.BaseURL == "" {
		return errors.New("settings.base_url must be set")
	}
	base := strings.TrimRight(st.BaseURL, "/")
	files := map[string][]byte{}
	for _, c := range cfg.Categories {
		if b := byCat[c.Name]; b != nil {
			files[c.Name+".list"] = b.body
		}
	}
	g := buildGeo(cfg.Categories, sets)
	for _, t := range g.sites {
		b, err := emit.SingBoxDomains(t.Domains)
		if err != nil {
			return err
		}
		files["geosite-"+t.Name+".json"] = b
	}
	for _, t := range g.ips {
		b, err := emit.SingBoxCIDRs(t.CIDRs)
		if err != nil {
			return err
		}
		files["geoip-"+t.Name+".json"] = b
	}
	xr, err := emit.XrayRouting(g.route)
	if err != nil {
		return err
	}
	files["xray-routing.json"] = xr
	vn, err := emit.V2rayNRouting(g.route)
	if err != nil {
		return err
	}
	files["v2rayn-routing.json"] = vn
	// Files that embed rule-set URLs.
	cd := confData(cfg.Categories, byCat)
	cd.Repo, cd.BaseURL, cd.Self = st.RepoURL, base, "shadowrocket.conf"
	conf, err := emit.Conf(cd)
	if err != nil {
		return err
	}
	files[cd.Self] = conf
	files["routing-rules.sgmodule"] = emit.Module(conf)
	rd := g.route
	rd.BaseURL = base
	route, err := emit.SingBoxRoute(rd)
	if err != nil {
		return err
	}
	files["sing-box-route.json"] = route
	// Karing has no PROXY group here: unmatched traffic already uses the selected proxy.
	karing, err := emit.Karing(base, g.route.Reject, slices.Concat(g.route.Private, g.route.Direct))
	if err != nil {
		return err
	}
	files["karing.json"] = karing
	files["clash.yaml"] = emit.ClashConfig(base, g.route)

	// DNS blocks for the cores (merge with the matching route/rules file).
	files["sing-box-dns.json"] = emit.SingBoxDNS()
	files["xray-dns.json"] = emit.XrayDNS()
	files["clash-dns.yaml"] = emit.ClashDNS()

	// Networks: WireGuard split tunnel, routers, ipset/nftables.
	var private []netip.Prefix
	if p := byCat["private"]; p != nil {
		private = p.cidrs
	}
	if ir := byCat["ir-cidr"]; ir != nil {
		files["wireguard-allowedips.txt"] = emit.WireGuardAllowedIPs(merge.Complement(append(slices.Clone(ir.cidrs), private...)))
		// Lite: IPv4 only, excluding Iranian blocks of /20 or larger (~1.7k entries,
		// ~90% of Iranian IPv4 space still direct), all IPv6 through the tunnel.
		// ponytail: fixed /20 cut-off; make it configurable if mobile clients need fewer entries.
		var big []netip.Prefix
		for _, p := range ir.cidrs {
			if p.Addr().Is4() && p.Bits() <= 20 {
				big = append(big, p)
			}
		}
		var lite []netip.Prefix
		for _, p := range merge.Complement(append(big, private...)) {
			if p.Addr().Is4() {
				lite = append(lite, p)
			}
		}
		files["wireguard-allowedips-lite.txt"] = emit.WireGuardAllowedIPs(append(lite, netip.MustParsePrefix("::/0")))
		files["openvpn-ir-routes.txt"] = emit.OpenVPN(ir.cidrs)
		files["openvpn-ir-routes-lite.txt"] = emit.OpenVPN(big) // same /20 cut-off as the WireGuard lite file
		files["mikrotik-ir.rsc"] = emit.MikroTik("IR", ir.cidrs)
		files["ir-cidr4.txt"] = emit.CIDRLines(ir.cidrs, true)
		files["ir-cidr6.txt"] = emit.CIDRLines(ir.cidrs, false)
	}
	// Browsers: PAC with Iranian domains (including the TLDs) direct.
	for _, t := range g.sites {
		if t.Name == "ir" {
			files["proxy.pac"] = emit.PAC(t.Domains)
		}
	}
	// Surge/Shadowrocket DOMAIN-SET, AdGuard Home / Pi-hole blocklists.
	var defaultBlock []parse.Domain
	for _, c := range cfg.Categories {
		b := byCat[c.Name]
		if b == nil || len(b.domains) == 0 {
			continue
		}
		files[c.Name+"-domainset.txt"] = emit.DomainSet(b.domains)
		if c.Action == "reject" || dnsBlocklist[c.Name] {
			files["adguard-"+c.Name+".txt"] = emit.AdGuard(c.Name, b.domains)
		}
		if c.Action == "reject" {
			defaultBlock = append(defaultBlock, b.domains...)
		}
	}
	defaultBlock = merge.Domains(defaultBlock)
	files["adguard-default.txt"] = emit.AdGuard("default (malware, phishing, cryptominers, ads)", defaultBlock)
	files["hosts-default.txt"] = emit.Hosts("default (malware, phishing, cryptominers, ads)", defaultBlock)
	if err := syncDir(dist, files); err != nil {
		return err
	}

	// Binary outputs, published as release assets by CI.
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	if err := writeIfChanged(filepath.Join(outDir, "geosite.dat"), emit.GeoSiteDat(g.sites)); err != nil {
		return err
	}
	if err := writeIfChanged(filepath.Join(outDir, "geoip.dat"), emit.GeoIPDat(g.ips)); err != nil {
		return err
	}
	// Known answers for cmd/verify, which tests the compiled files with the clients' tools.
	checks, err := json.MarshalIndent(toolChecks(canaries, g), "", "  ")
	if err != nil {
		return err
	}
	if err := writeIfChanged(filepath.Join(outDir, "checks.json"), checks); err != nil {
		return err
	}
	// mihomo text rule-sets; CI compiles them to .mrs with the official binary.
	clash := map[string][]byte{}
	for _, t := range g.sites {
		b, skipped := emit.ClashDomainText(t.Domains)
		if skipped > 0 {
			log.Printf("::warning::geosite-%s: %d DOMAIN-KEYWORD rules skipped for mihomo (domain behavior has no keyword type)", t.Name, skipped)
		}
		if len(b) > 0 {
			clash["geosite-"+t.Name+".txt"] = b
		}
	}
	for _, t := range g.ips {
		clash["geoip-"+t.Name+".txt"] = emit.ClashIPText(t.CIDRs)
	}
	if err := syncDir(filepath.Join(outDir, "clash"), clash); err != nil {
		return err
	}
	health := map[string]outcome{}
	for _, o := range outcomes {
		health[o.src.Name] = o
	}
	return writeDocs(docsDir, cfg, g, files, health, sets)
}
