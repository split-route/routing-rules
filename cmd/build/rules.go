// Merged rules: config rule order, canaries and v2ray/sing-box geo tags.

package main

import (
	"errors"
	"fmt"
	"log"
	"net/netip"
	"os"
	"slices"
	"strings"

	"github.com/split-route/routing-rules/internal/emit"
	"github.com/split-route/routing-rules/internal/merge"
	"github.com/split-route/routing-rules/internal/parse"
)

// dnsBlocklist marks optional categories that also make sense as DNS blocklists.
var dnsBlocklist = map[string]bool{"ads-lite": true, "tracking": true, "nsfw": true, "social": true}

// ruleOrder lists the merged lists in the default configs' rule order:
// block page, private, REJECT lists, PROXY lists, the .ir TLDs, DIRECT lists.
func ruleOrder(cats []Category, sets map[string]catRules) []merge.Ruleset {
	verdict := map[string]string{"reject": "block", "proxy": "proxy", "direct": "direct"}
	add := func(out []merge.Ruleset, c Category) []merge.Ruleset {
		if r, ok := sets[c.Name]; ok {
			out = append(out, merge.Ruleset{Name: c.Name, Action: verdict[c.Action], Domains: r.domains, CIDRs: r.cidrs})
		}
		return out
	}
	var out []merge.Ruleset
	for _, name := range []string{"ir-injected", "private"} {
		for _, c := range cats {
			if c.Name == name {
				out = add(out, c)
			}
		}
	}
	for _, phase := range []string{"reject", "proxy", "direct"} {
		if phase == "direct" { // DOMAIN-SUFFIX,ir / .ایران come before the DIRECT lists
			out = append(out, merge.Ruleset{Name: "ir TLD", Action: "direct", Domains: []parse.Domain{
				{Kind: parse.Suffix, Name: "ir"}, {Kind: parse.Suffix, Name: "xn--mgba3a4f16a"}}})
		}
		for _, c := range cats {
			if c.Action == phase && c.Name != "ir-injected" && c.Name != "private" {
				out = add(out, c)
			}
		}
	}
	return out
}

// canary is one known answer and the list that decided it ("" for the
// default rule).
type canary struct{ q, want, list string }

// checkCanaries evaluates the known answers in path against the merged lists
// and returns them with one message per mismatch.
func checkCanaries(path string, cats []Category, sets map[string]catRules) (cs []canary, bad []string) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, []string{err.Error()}
	}
	order := ruleOrder(cats, sets)
	for i, line := range strings.Split(string(b), "\n") {
		line, _, _ = strings.Cut(line, "#")
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		if len(f) != 2 || (f[1] != "block" && f[1] != "direct" && f[1] != "proxy") {
			bad = append(bad, fmt.Sprintf("%s:%d: want \"<domain or IP> <block|direct|proxy>\"", path, i+1))
			continue
		}
		got, list := merge.Verdict(order, f[0])
		cs = append(cs, canary{f[0], f[1], list})
		if got != f[1] {
			if list == "" {
				list = "no list, default rule"
			}
			bad = append(bad, fmt.Sprintf("canary %s: want %s, got %s (%s)", f[0], f[1], got, list))
		}
	}
	log.Printf("canaries: %d checked, %d failed", len(cs), len(bad))
	return cs, bad
}

// toolCheck is one known answer for cmd/verify to test against a compiled
// file with the clients' own tools: File is a rule-set name ("geosite-ir")
// or "proxy.pac"; Match says whether it must match Q (for the PAC file:
// whether Q must go DIRECT).
type toolCheck struct {
	File  string `json:"file"`
	Q     string `json:"q"`
	Match bool   `json:"match"`
}

// toolChecks turns the canaries into checks of the compiled rule-sets and
// the PAC file, so data/canaries.txt is the one list of known answers.
func toolChecks(cs []canary, g geo) []toolCheck {
	have := map[string]bool{}
	for _, t := range g.sites {
		have["geosite-"+t.Name] = true
	}
	for _, t := range g.ips {
		have["geoip-"+t.Name] = true
	}
	var out []toolCheck
	for _, c := range cs {
		_, err := netip.ParseAddr(c.q)
		kind := "geosite-"
		if err == nil {
			kind = "geoip-"
		}
		switch c.list {
		case "": // default rule: no rule-set the configs use may match
			for _, f := range g.route.All {
				if strings.HasPrefix(f, kind) {
					out = append(out, toolCheck{f, c.q, false})
				}
			}
			if kind == "geosite-" {
				out = append(out, toolCheck{"proxy.pac", c.q, false})
			}
			continue
		case "ir TLD", "ir-domains": // geosite-ir and the PAC file include the TLDs
			out = append(out, toolCheck{"proxy.pac", c.q, true})
		case "private": // the PAC file only knows plain names and private IPv4 ranges
			if err == nil || !strings.Contains(c.q, ".") {
				out = append(out, toolCheck{"proxy.pac", c.q, true})
			}
		}
		tag := "ir"
		if c.list != "ir TLD" {
			tag = geoTag(c.list)
		}
		if have[kind+tag] {
			out = append(out, toolCheck{kind + tag, c.q, true})
		}
	}
	return out
}

type catRules struct {
	cat     Category
	domains []parse.Domain
	cidrs   []netip.Prefix
}

type geo struct {
	sites, ips []emit.Tag
	route      emit.RouteData // BaseURL is filled in by the caller
}

// geoTag maps a category to its v2ray/sing-box tag (Chocolate4U naming).
func geoTag(cat string) string {
	switch cat {
	case "ir-domains", "ir-cidr":
		return "ir"
	case "tor-exit":
		return "tor"
	}
	return cat
}

// buildGeo derives v2ray/sing-box tags and the routing template data.
func buildGeo(cats []Category, sets map[string]catRules) geo {
	var g geo
	rd := &g.route
	for _, c := range cats {
		r, ok := sets[c.Name]
		if !ok {
			continue
		}
		tag := geoTag(c.Name)
		if len(r.domains) > 0 {
			d := r.domains
			if tag == "ir" { // the Shadowrocket config adds these TLDs inline instead
				d = merge.Domains(append(slices.Clone(d), parse.Domain{Kind: parse.Suffix, Name: "ir"}, parse.Domain{Kind: parse.Suffix, Name: "xn--mgba3a4f16a"}))
			}
			g.sites = append(g.sites, emit.Tag{Name: tag, Domains: d})
			if c.Name == "private" {
				rd.Private = append(rd.Private, "geosite-"+tag)
				rd.PrivateSite = append(rd.PrivateSite, tag)
				goto ips
			}
			switch c.Action {
			case "reject":
				rd.Reject = append(rd.Reject, "geosite-"+tag)
				rd.RejectSite = append(rd.RejectSite, tag)
			case "proxy":
				rd.Proxy = append(rd.Proxy, "geosite-"+tag)
				rd.ProxySite = append(rd.ProxySite, tag)
			case "direct":
				rd.Direct = append(rd.Direct, "geosite-"+tag)
			}
		}
	ips:
		if len(r.cidrs) > 0 {
			g.ips = append(g.ips, emit.Tag{Name: tag, CIDRs: r.cidrs})
			if c.Name == "private" {
				rd.Private = append(rd.Private, "geoip-"+tag)
				rd.PrivateIP = append(rd.PrivateIP, tag)
				continue
			}
			if c.Name == "ir-injected" {
				rd.Early = append(rd.Early, "geoip-"+tag)
				rd.EarlyIP = append(rd.EarlyIP, tag)
			}
			switch c.Action {
			case "reject":
				rd.Reject = append(rd.Reject, "geoip-"+tag)
				rd.RejectIP = append(rd.RejectIP, tag)
			case "proxy":
				rd.Proxy = append(rd.Proxy, "geoip-"+tag)
			case "direct":
				rd.Direct = append(rd.Direct, "geoip-"+tag)
			}
		}
	}
	// The route only declares the rule-sets its rules reference.
	rd.All = slices.Concat(rd.Private, rd.Reject, rd.Proxy, rd.Direct)
	byName := func(a, b emit.Tag) int { return strings.Compare(a.Name, b.Name) }
	slices.SortFunc(g.sites, byName)
	slices.SortFunc(g.ips, byName)
	return g
}

func confData[T any](cats []Category, on map[string]T) emit.ConfData {
	d := emit.ConfData{On: map[string]bool{}}
	for _, c := range cats {
		if _, ok := on[c.Name]; !ok {
			continue
		}
		switch c.Name {
		case "sanctioned":
			d.On["sanctioned"] = true
		case "private":
			d.On["private"] = true
		case "ir-injected": // Iran's block-page range is inside 10.0.0.0/8
			d.Early = append(d.Early, c.Name)
		case "ir-domains":
			d.On["irdomains"] = true
		case "ir-cidr":
			d.On["ircidr"] = true
		default:
			if c.Action == "reject" {
				d.Reject = append(d.Reject, c.Name)
			}
		}
	}
	return d
}

func ruleString(d parse.Domain) string {
	return [...]string{"DOMAIN-SUFFIX", "DOMAIN", "DOMAIN-KEYWORD"}[d.Kind] + "," + d.Name
}
