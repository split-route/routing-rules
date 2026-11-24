package emit

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/split-route/routing-rules/internal/parse"
)

func TestList(t *testing.T) {
	got := string(List("x", []parse.Domain{{Kind: parse.Suffix, Name: "a.com"}, {Kind: parse.Exact, Name: "b.com"}, {Kind: parse.Keyword, Name: "ad"}},
		[]netip.Prefix{netip.MustParsePrefix("1.2.3.0/24"), netip.MustParsePrefix("2001:db8::/32")}))
	want := "# NAME: x\n# TOTAL: 5\nDOMAIN-SUFFIX,a.com\nDOMAIN,b.com\nDOMAIN-KEYWORD,ad\nIP-CIDR,1.2.3.0/24,no-resolve\nIP-CIDR6,2001:db8::/32,no-resolve\n"
	if got != want {
		t.Fatalf("got\n%s", got)
	}
	if n := CountRules([]byte(got)); n != 5 {
		t.Fatalf("count %d", n)
	}
}

func TestConfOrder(t *testing.T) {
	b, err := Conf(ConfData{BaseURL: "https://x/dist", Self: "shadowrocket.conf", Reject: []string{"malware", "ads"}, On: map[string]bool{"sanctioned": true, "irdomains": true, "ircidr": true}})
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	order := []string{"[General]", "malware.list,REJECT", "ads.list,REJECT", "sanctioned.list,PROXY", "DOMAIN-SUFFIX,ir,DIRECT", "ir-domains.list,DIRECT", "ir-cidr.list,DIRECT", "GEOIP,IR,DIRECT", "FINAL,PROXY", "[Host]"}
	last := -1
	for _, o := range order {
		i := strings.Index(s, o)
		if i <= last {
			t.Fatalf("%q out of order or missing:\n%s", o, s)
		}
		last = i
	}
}
