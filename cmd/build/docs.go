// Generated docs pages and the domain checker's lists.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/split-route/routing-rules/internal/emit"
	"github.com/split-route/routing-rules/internal/merge"
	"github.com/split-route/routing-rules/internal/parse"
)

// docLangs are the docs languages with generated pages (docs/<lang>/); each
// non-English one has internal/emit/locales/<lang>.json and Persian templates.
var docLangs = []string{"en", "fa"}

// writeDocs generates the URL and sources pages for every docs language, plus
// the checker's lists (English folder only; other languages load them from there).
func writeDocs(dir string, cfg Config, g geo, files map[string][]byte, health map[string]outcome, sets map[string]catRules) error {
	// Which sources contributed each rule: the checker's .prov files, and
	// each source's rule counts on the sources page.
	prov := map[string][]byte{}
	stats := map[string]srcStat{}
	for _, c := range cfg.Categories {
		if _, ok := files[c.Name+".list"]; !ok {
			continue
		}
		var srcs []outcome
		for _, s := range cfg.Sources {
			if o := health[s.Name]; s.Category == c.Name && (o.status == "ok" || o.status == "stale") {
				srcs = append(srcs, o)
			}
		}
		prov[c.Name], stats = provenance(srcs, sets[c.Name], stats)
	}
	for _, lang := range docLangs {
		l, err := emit.LoadLocale(lang)
		if err != nil {
			return err
		}
		if err := writeDocPages(filepath.Join(dir, lang), l, cfg, g, files, health, stats); err != nil {
			return err
		}
		if m := l.Missing(); len(m) > 0 {
			log.Printf("::warning::docs: %d strings missing from locales/%s.json: %q", len(m), lang, m)
		}
	}
	en := filepath.Join(dir, "en")
	// Download buttons link to release assets (GitHub sends them as attachments),
	// so docs/en/dl only holds the checker's lists; clear any older file copies.
	if err := syncDir(filepath.Join(en, "dl"), nil); err != nil {
		return err
	}
	// Lists the "Check a domain" page loads. Same origin: GitHub's download
	// servers are not known to allow cross-origin fetches.
	checkLists := map[string][]byte{}
	for _, c := range cfg.Categories {
		if b, ok := files[c.Name+".list"]; ok && c.Name != "ads-lite" {
			checkLists[c.Name+".list"] = b
			checkLists[c.Name+".prov"] = prov[c.Name]
		}
	}
	return syncDir(filepath.Join(en, "dl", "lists"), checkLists)
}

// writeDocPages writes urls.md and sources.md for one language.
func writeDocPages(dir string, l *emit.Locale, cfg Config, g geo, files map[string][]byte, health map[string]outcome, stats map[string]srcStat) error {
	st := cfg.Settings
	T, Tf := l.T, l.Tf
	row := func(name, desc string) emit.DocRow {
		return emit.DocRow{Name: name, Desc: desc, URL: strings.TrimRight(st.BaseURL, "/") + "/" + name}
	}
	var sr, lists, sb, sbSrc, v2 emit.DocSection
	sr = emit.DocSection{ID: emit.Slug("Shadowrocket config"), Title: T("Shadowrocket config"), Note: T("Import this URL in Shadowrocket.")}
	sr.Rows = append(sr.Rows,
		row("shadowrocket.conf", T("Full config: REJECT, PROXY, DIRECT, GEOIP, FINAL")),
		row("routing-rules.sgmodule", T("Shadowrocket / Surge module: the same rules on top of your own config")))
	lists = emit.DocSection{ID: emit.Slug("Shadowrocket rule lists"), Title: T("Shadowrocket rule lists"), Note: T("`RULE-SET` text lists in Shadowrocket (Surge-style) syntax.")}
	for _, c := range cfg.Categories {
		if _, ok := files[c.Name+".list"]; ok {
			lists.Rows = append(lists.Rows, row(c.Name+".list", fmt.Sprintf("%s (%s)", c.Name, T(actionDesc(c.Action)))))
		}
	}
	sb = emit.DocSection{ID: emit.Slug("sing-box rule-sets (binary)"), Title: T("sing-box rule-sets (binary)"), Note: T("Use with `\"type\": \"remote\", \"format\": \"binary\"`. Rule-set format version 1, readable by sing-box 1.8+ and apps built on it.")}
	sbSrc = emit.DocSection{ID: emit.Slug("sing-box rule-sets (JSON source)"), Title: T("sing-box rule-sets (JSON source)"), Note: T("Same rules as JSON, for `\"format\": \"source\"`.")}
	sb.Rows = append(sb.Rows, row("sing-box-route.json", T("Example `route` block (sing-box 1.11+)")))
	for _, t := range g.sites {
		sb.Rows = append(sb.Rows, row("geosite-"+t.Name+".srs", Tf("domains: %s", t.Name)))
		sbSrc.Rows = append(sbSrc.Rows, row("geosite-"+t.Name+".json", Tf("domains: %s", t.Name)))
	}
	for _, t := range g.ips {
		sb.Rows = append(sb.Rows, row("geoip-"+t.Name+".srs", Tf("IP ranges: %s", t.Name)))
		sbSrc.Rows = append(sbSrc.Rows, row("geoip-"+t.Name+".json", Tf("IP ranges: %s", t.Name)))
	}
	apps := emit.DocSection{ID: emit.Slug("App configs"), Title: T("App configs"), Note: T("Ready-made rule files for apps with their own format.")}
	apps.Rows = append(apps.Rows,
		row("karing.json", T("Karing custom diversion rules (Block, Direct)")),
		row("v2rayn-routing.json", T("v2rayN / v2rayNG custom routing rules (needs this project's geo files)")))
	cl := emit.DocSection{ID: emit.Slug("Clash / mihomo rule-sets"), Title: T("Clash / mihomo rule-sets"), Note: T("`.mrs` rule-sets for mihomo (`format: mrs`; `behavior: domain` for geosite, `ipcidr` for geoip).")}
	cl.Rows = append(cl.Rows, row("clash.yaml", T("`rule-providers` and `rules` to merge into a mihomo config (needs a `PROXY` group)")))
	for _, t := range g.sites {
		cl.Rows = append(cl.Rows, row("geosite-"+t.Name+".mrs", Tf("domains: %s", t.Name)))
	}
	for _, t := range g.ips {
		cl.Rows = append(cl.Rows, row("geoip-"+t.Name+".mrs", Tf("IP ranges: %s", t.Name)))
	}
	dns := emit.DocSection{ID: emit.Slug("DNS"), Title: T("DNS"), Note: T("DNS blocks to merge into a config: Iranian domains use the device's DNS, everything else DNS-over-HTTPS through the proxy.")}
	dns.Rows = append(dns.Rows,
		row("sing-box-dns.json", T("sing-box `dns` block (1.12+), merge with `sing-box-route.json`")),
		row("xray-dns.json", T("Xray `dns` block, merge with `xray-routing.json`")),
		row("clash-dns.yaml", T("mihomo `dns` block (fake-ip), merge with `clash.yaml`")))
	net := emit.DocSection{ID: emit.Slug("Networks and routers"), Title: T("Networks and routers"), Note: T("Iran IP ranges for split tunneling outside proxy apps.")}
	net.Rows = append(net.Rows,
		row("wireguard-allowedips.txt", T("WireGuard `AllowedIPs`: everything except Iran and private ranges (~16k entries, desktop)")),
		row("wireguard-allowedips-lite.txt", T("Same, IPv4 only, larger Iranian blocks only (~1.7k entries, phones); all IPv6 through the tunnel")),
		row("openvpn-ir-routes.txt", T("OpenVPN client routes that keep Iranian IPv4 ranges outside the tunnel")),
		row("openvpn-ir-routes-lite.txt", T("Same, larger Iranian blocks only (fewer routes, for phones)")),
		row("mikrotik-ir.rsc", T("MikroTik RouterOS script: address list `IR`")),
		row("ir-cidr4.txt", T("Iran IPv4 ranges, one per line (ipset, nftables)")),
		row("ir-cidr6.txt", T("Iran IPv6 ranges, one per line")),
		row("proxy.pac", T("Browser PAC: Iranian domains direct, everything else through a local proxy")))
	dnsf := emit.DocSection{ID: emit.Slug("AdGuard Home and Pi-hole"), Title: T("AdGuard Home and Pi-hole"), Note: T("DNS blocklists. `adguard-*` work in AdGuard Home and Pi-hole v6; `hosts-*` in Pi-hole v5 and hosts files.")}
	dnsf.Rows = append(dnsf.Rows,
		row("adguard-default.txt", T("Everything the default configs block (malware, phishing, cryptominers, ads)")),
		row("hosts-default.txt", T("Same, hosts format (subdomains not blocked)")))
	ds := emit.DocSection{ID: emit.Slug("Surge / Shadowrocket DOMAIN-SET"), Title: T("Surge / Shadowrocket DOMAIN-SET"), Note: T("Domain-only lists in DOMAIN-SET format, lighter than the same rules as RULE-SET.")}
	for _, c := range cfg.Categories {
		if _, ok := files["adguard-"+c.Name+".txt"]; ok {
			dnsf.Rows = append(dnsf.Rows, row("adguard-"+c.Name+".txt", c.Name))
		}
		if _, ok := files[c.Name+"-domainset.txt"]; ok {
			ds.Rows = append(ds.Rows, row(c.Name+"-domainset.txt", c.Name))
		}
	}
	v2 = emit.DocSection{ID: emit.Slug("v2ray / Xray"), Title: T("v2ray / Xray"), Note: T("Tags: `geosite:<name>` and `geoip:<name>` as listed in the sing-box tables above.")}
	v2.Rows = append(v2.Rows,
		row("geosite.dat", T("Domain tags")),
		row("geoip.dat", T("IP tags, including `private`")),
		row("xray-routing.json", T("Example `routing` block")))

	urls, err := emit.URLsPage(l.Lang, emit.URLDocs{RepoURL: st.RepoURL, BaseURL: strings.TrimRight(st.BaseURL, "/"), Sections: []emit.DocSection{sr, lists, ds, apps, sb, sbSrc, cl, v2, dns, net, dnsf}})
	if err != nil {
		return err
	}
	var srcs []emit.DocSource
	for _, s := range cfg.Sources {
		o := health[s.Name]
		day := func(t time.Time) string {
			if t.IsZero() {
				return "—"
			}
			return t.UTC().Format("2006-01-02")
		}
		label := map[string]string{"ok": "✅ " + T("healthy"), "stale": "⚠️ " + T("stale (last good copy)"),
			"FAILED": "❌ " + T("failing"), "disabled": T("disabled"), "skipped": T("disabled")}[o.status]
		srcs = append(srcs, emit.DocSource{Name: s.Name, URL: sourceLink(st.RepoURL, s.URL), Category: s.Category, Format: s.Format, License: s.License,
			Enabled: s.Enabled && catEnabled(cfg, s.Category), Health: label, Fetched: day(o.fetched), Changed: day(o.changed), Rules: "—", Unique: "—"})
		if n, ok := stats[s.Name]; ok {
			srcs[len(srcs)-1].Rules, srcs[len(srcs)-1].Unique = strconv.Itoa(n.rules), strconv.Itoa(n.unique)
		}
	}
	sources, err := emit.SourcesPage(l.Lang, srcs)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := writeIfChanged(filepath.Join(dir, "urls.md"), urls); err != nil {
		return err
	}
	return writeIfChanged(filepath.Join(dir, "sources.md"), sources)
}

// provenance renders which sources contributed each rule of a list, for the
// domain checker: a JSON array of source names, then one hex bitmask per rule
// in the list's order (domains, then IP ranges). A domain rule comes from the
// sources that contained exactly that rule; a merged IP range from every
// source with an overlapping range.
// It also counts, per source, the final rules it contributed to and those no
// other source contributed to (added to stats, keyed by source name); for IP
// ranges, "unique" is the prefixes covering addresses no other source has.
// ponytail: overlap check is O(final × source ranges), fine for a few thousand ranges.
func provenance(srcs []outcome, final catRules, stats map[string]srcStat) ([]byte, map[string]srcStat) {
	if len(srcs) > 64 { // one bit per source in a uint64
		log.Printf("::warning::%s: %d sources; the checker only names the first 64 as a rule's source", final.cat.Name, len(srcs))
		srcs = srcs[:64]
	}
	names := make([]string, len(srcs))
	byRule := map[parse.Domain]uint64{}
	for i, o := range srcs {
		names[i] = o.src.Name
		for _, d := range o.res.Domains {
			byRule[d] |= 1 << i
		}
	}
	var b bytes.Buffer
	nb, _ := json.Marshal(names)
	b.Write(nb)
	b.WriteByte('\n')
	count := func(mask uint64, unique bool) {
		for i := range srcs {
			if mask&(1<<i) == 0 {
				continue
			}
			n := stats[names[i]]
			n.rules++
			if unique && mask&(mask-1) == 0 { // only this source
				n.unique++
			}
			stats[names[i]] = n
		}
	}
	for _, d := range final.domains {
		fmt.Fprintf(&b, "%x\n", byRule[d])
		count(byRule[d], true)
	}
	for _, p := range final.cidrs {
		var mask uint64
		for i, o := range srcs {
			for _, q := range o.res.CIDRs {
				if p.Overlaps(q) {
					mask |= 1 << i
					break
				}
			}
		}
		fmt.Fprintf(&b, "%x\n", mask)
		count(mask, false)
	}
	// A merged range overlaps several sources, so its mask says little about
	// what a source adds. Count instead the prefixes covering the final list's
	// addresses that no other source has: what dropping the source would lose.
	if len(final.cidrs) > 0 {
		outside := merge.Complement(final.cidrs)
		for i, o := range srcs {
			if len(o.res.CIDRs) == 0 {
				continue
			}
			drop := slices.Clone(outside)
			for j, p := range srcs {
				if j != i {
					drop = append(drop, p.res.CIDRs...)
				}
			}
			n := stats[names[i]]
			n.unique += len(merge.Subtract(o.res.CIDRs, drop))
			stats[names[i]] = n
		}
	}
	for _, n := range names { // a source whose rules were all filtered out still gets a row
		stats[n] = stats[n]
	}
	return b.Bytes(), stats
}

// srcStat counts a source's rules in its category's final list.
type srcStat struct{ rules, unique int }

// sourceLink turns a file:/// source (a file in this repo) into a link to it
// on the repo's default branch; other URLs are returned as is.
func sourceLink(repo, u string) string {
	p, ok := strings.CutPrefix(u, "file:///")
	if !ok || repo == "" {
		return u
	}
	return strings.TrimRight(repo, "/") + "/blob/HEAD/" + p
}

func actionDesc(a string) string {
	switch a {
	case "reject":
		return "REJECT"
	case "proxy":
		return "PROXY"
	case "direct":
		return "DIRECT"
	}
	return "optional, not in the config"
}
