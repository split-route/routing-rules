package merge

import (
	"net/netip"
	"strings"

	"github.com/split-route/routing-rules/internal/parse"
)

// Ruleset is one list in rule order, with what matching it does.
type Ruleset struct {
	Name    string
	Action  string // reject | direct | proxy
	Domains []parse.Domain
	CIDRs   []netip.Prefix
}

// Verdict returns the action and list of the first ruleset matching q, or
// ("proxy", "") when nothing matches (the configs' final rule). q is a domain
// or an IP address; domains are not resolved, mirroring no-resolve rules.
func Verdict(order []Ruleset, q string) (action, list string) {
	q = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(q)), ".")
	ip, ipErr := netip.ParseAddr(q)
	for _, r := range order {
		if ipErr == nil {
			for _, p := range r.CIDRs {
				if p.Contains(ip.Unmap()) {
					return r.Action, r.Name
				}
			}
			continue
		}
		for _, d := range r.Domains {
			if Matches(d, q) {
				return r.Action, r.Name
			}
		}
	}
	return "proxy", ""
}
