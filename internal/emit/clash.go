package emit

import (
	"bytes"
	"fmt"
	"net/netip"
	"strings"

	"github.com/split-route/routing-rules/internal/parse"
)

// ClashDomainText renders a mihomo "domain" behavior text rule-set (the input
// for convert-ruleset → .mrs): "+.x" matches x and its subdomains, "x" only x.
// The domain behavior has no keyword type; skipped keywords are counted.
func ClashDomainText(domains []parse.Domain) (b []byte, skipped int) {
	var buf bytes.Buffer
	for _, d := range domains {
		switch d.Kind {
		case parse.Suffix:
			buf.WriteString("+." + d.Name + "\n")
		case parse.Exact:
			buf.WriteString(d.Name + "\n")
		default:
			skipped++
		}
	}
	return buf.Bytes(), skipped
}

// ClashIPText renders a mihomo "ipcidr" behavior text rule-set.
func ClashIPText(cidrs []netip.Prefix) []byte {
	var buf bytes.Buffer
	for _, p := range cidrs {
		buf.WriteString(p.String() + "\n")
	}
	return buf.Bytes()
}

// ClashConfig renders mihomo rule-providers and rules (same order as the other
// configs) to merge into a user's config. It expects a proxy group named PROXY.
func ClashConfig(base string, d RouteData) []byte {
	var b strings.Builder
	b.WriteString("# routing-rules: rules for mihomo (Clash.Meta). Merge into your config.\n")
	b.WriteString("# Requires a proxy group named PROXY. Rule-sets refresh daily.\n")
	b.WriteString("rule-providers:\n")
	for _, t := range d.All {
		behavior := "domain"
		if strings.HasPrefix(t, "geoip-") {
			behavior = "ipcidr"
		}
		fmt.Fprintf(&b, "  %s:\n    type: http\n    behavior: %s\n    format: mrs\n    url: %s/%s.mrs\n    path: ./rules/%s.mrs\n    interval: 86400\n",
			t, behavior, base, t, t)
	}
	b.WriteString("\nrules:\n")
	rule := func(tag, target string) {
		suffix := ""
		if strings.HasPrefix(tag, "geoip-") && target != "DIRECT" {
			suffix = ",no-resolve" // never resolve a domain just to block or proxy it by IP
		}
		fmt.Fprintf(&b, "  - RULE-SET,%s,%s%s\n", tag, target, suffix)
	}
	for _, t := range d.Early {
		rule(t, "REJECT")
	}
	for _, t := range d.Private {
		if strings.HasPrefix(t, "geoip-") {
			fmt.Fprintf(&b, "  - RULE-SET,%s,DIRECT,no-resolve\n", t)
		} else {
			fmt.Fprintf(&b, "  - RULE-SET,%s,DIRECT\n", t)
		}
	}
	for _, t := range d.Reject {
		rule(t, "REJECT")
	}
	for _, t := range d.Proxy {
		rule(t, "PROXY")
	}
	for _, t := range d.Direct {
		rule(t, "DIRECT") // geoip-ir resolves, like Shadowrocket's GEOIP,IR
	}
	b.WriteString("  - MATCH,PROXY\n")
	return []byte(b.String())
}
