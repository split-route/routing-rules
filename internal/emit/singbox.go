package emit

import (
	"encoding/json"
	"net/netip"

	"github.com/split-route/routing-rules/internal/parse"
)

// SingBoxVersion is the source rule-set format version, which also sets the
// compiled .srs version. Version 1 covers every rule type used here and is
// read by all sing-box releases since 1.8; version 2 (sing-box 1.10+) only
// stores domain_suffix more compactly, and Karing's bundled core rejects it.
const SingBoxVersion = 1

type sbRule struct {
	Domain        []string `json:"domain,omitempty"`
	DomainSuffix  []string `json:"domain_suffix,omitempty"`
	DomainKeyword []string `json:"domain_keyword,omitempty"`
	IPCIDR        []string `json:"ip_cidr,omitempty"`
}

type sbRuleSet struct {
	Version int      `json:"version"`
	Rules   []sbRule `json:"rules"`
}

// SingBoxDomains renders a geosite-style source rule-set. A suffix rule is a
// single domain_suffix "x": sing-box matches it against x and its subdomains on
// a label boundary (not "xx"), both in the legacy expansion used for format
// version 1 and in the root-label form (sing/common/domain.NewMatcher).
func SingBoxDomains(domains []parse.Domain) ([]byte, error) {
	var r sbRule
	for _, d := range domains {
		switch d.Kind {
		case parse.Exact:
			r.Domain = append(r.Domain, d.Name)
		case parse.Suffix:
			r.DomainSuffix = append(r.DomainSuffix, d.Name)
		case parse.Keyword:
			r.DomainKeyword = append(r.DomainKeyword, d.Name)
		}
	}
	return render(r)
}

// SingBoxCIDRs renders a geoip-style source rule-set.
func SingBoxCIDRs(cidrs []netip.Prefix) ([]byte, error) {
	var r sbRule
	for _, p := range cidrs {
		r.IPCIDR = append(r.IPCIDR, p.String())
	}
	return render(r)
}

func render(r sbRule) ([]byte, error) {
	// Compact: these files are machine-read and no longer committed, so
	// indentation would only add size (it roughly doubled the larger files).
	b, err := json.Marshal(sbRuleSet{SingBoxVersion, []sbRule{r}})
	return append(b, '\n'), err
}
