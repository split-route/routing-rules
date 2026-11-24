package emit

import (
	"bytes"
	"encoding/json"
	"net/netip"
	"strings"
	"testing"

	"github.com/split-route/routing-rules/internal/parse"
	"gopkg.in/yaml.v3"
)

func TestGeoSiteDat(t *testing.T) {
	got := GeoSiteDat([]Tag{{Name: "ads", Domains: []parse.Domain{
		{Kind: parse.Suffix, Name: "a.io"}, {Kind: parse.Exact, Name: "b.io"}, {Kind: parse.Keyword, Name: "ad"},
	}}})
	want := []byte{
		0x0a, 0x1f, // GeoSiteList.entry, len 31
		0x0a, 0x03, 'A', 'D', 'S', // country_code
		0x12, 0x08, 0x08, 0x02, 0x12, 0x04, 'a', '.', 'i', 'o', // RootDomain a.io
		0x12, 0x08, 0x08, 0x03, 0x12, 0x04, 'b', '.', 'i', 'o', // Full b.io
		0x12, 0x04, 0x12, 0x02, 'a', 'd', // Plain (type omitted) ad
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got  % x\nwant % x", got, want)
	}
}

func TestGeoIPDat(t *testing.T) {
	got := GeoIPDat([]Tag{{Name: "ir", CIDRs: []netip.Prefix{netip.MustParsePrefix("2.144.0.0/14"), netip.MustParsePrefix("::/0")}}})
	want := []byte{
		0x0a, 0x22, // GeoIPList.entry, len 34
		0x0a, 0x02, 'I', 'R',
		0x12, 0x08, 0x0a, 0x04, 2, 144, 0, 0, 0x10, 14, // v4, prefix 14
		0x12, 0x12, 0x0a, 0x10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ::/0, prefix omitted
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got  % x\nwant % x", got, want)
	}
}

func TestSingBox(t *testing.T) {
	b, err := SingBoxDomains([]parse.Domain{{Kind: parse.Suffix, Name: "z.io"}, {Kind: parse.Exact, Name: "a.io"}, {Kind: parse.Keyword, Name: "ad"}})
	if err != nil {
		t.Fatal(err)
	}
	var rs struct {
		Version int `json:"version"`
		Rules   []map[string][]string
	}
	if err := json.Unmarshal(b, &rs); err != nil {
		t.Fatal(err)
	}
	r := rs.Rules[0]
	if rs.Version != SingBoxVersion || strings.Join(r["domain"], ",") != "a.io" || strings.Join(r["domain_suffix"], ",") != "z.io" || strings.Join(r["domain_keyword"], ",") != "ad" {
		t.Fatalf("unexpected rule-set: %s", b)
	}
}

func TestRouteTemplates(t *testing.T) {
	d := RouteData{BaseURL: "https://x/release",
		Reject: []string{"geosite-ads", "geoip-malware"}, Proxy: []string{"geosite-sanctioned"}, Direct: []string{"geosite-ir", "geoip-ir"},
		All:        []string{"geosite-ads", "geoip-malware", "geosite-sanctioned", "geosite-ir", "geoip-ir"},
		RejectSite: []string{"ads"}, RejectIP: []string{"malware"}, ProxySite: []string{"sanctioned"}}
	for name, fn := range map[string]func(RouteData) ([]byte, error){"sing-box": SingBoxRoute, "xray": XrayRouting} {
		b, err := fn(d)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !strings.Contains(string(b), "ads") || !strings.Contains(string(b), "sanctioned") {
			t.Fatalf("%s: missing tags:\n%s", name, b)
		}
	}
	d.Proxy, d.ProxySite, d.RejectIP = nil, nil, nil // optional sections must still render valid JSON
	for _, fn := range []func(RouteData) ([]byte, error){SingBoxRoute, XrayRouting} {
		if _, err := fn(d); err != nil {
			t.Fatal(err)
		}
	}
}

func TestKaring(t *testing.T) {
	b, err := Karing("https://x/release", []string{"geosite-ads", "geoip-malware"}, []string{"geosite-ir", "geoip-ir", "geoip-private"})
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Rules []struct {
			Outbound string   `json:"outbound"`
			Switch   bool     `json:"switch"`
			Or       bool     `json:"or"`
			RuleSet  []string `json:"rule_set"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	r := doc.Rules
	if len(r) != 2 || r[0].Outbound != "block" || r[1].Outbound != "direct" || !r[0].Switch || !r[0].Or ||
		r[0].RuleSet[0] != "https://x/release/geosite-ads.srs" || r[1].RuleSet[2] != "https://x/release/geoip-private.srs" {
		t.Fatalf("unexpected karing rules:\n%s", b)
	}
}

func TestClash(t *testing.T) {
	txt, skipped := ClashDomainText([]parse.Domain{{Kind: parse.Suffix, Name: "ir"}, {Kind: parse.Exact, Name: "a.io"}, {Kind: parse.Keyword, Name: "ad"}})
	if string(txt) != "+.ir\na.io\n" || skipped != 1 {
		t.Fatalf("domain text %q skipped %d", txt, skipped)
	}
	if got := string(ClashIPText([]netip.Prefix{netip.MustParsePrefix("2.144.0.0/14")})); got != "2.144.0.0/14\n" {
		t.Fatalf("ip text %q", got)
	}
	conf := string(ClashConfig("https://x/release", RouteData{
		Private: []string{"geosite-private", "geoip-private"},
		Reject:  []string{"geosite-ads", "geoip-malware"}, Proxy: []string{"geosite-sanctioned"}, Direct: []string{"geosite-ir", "geoip-ir"},
		All: []string{"geosite-ads", "geoip-malware", "geosite-sanctioned", "geosite-ir", "geoip-ir", "geoip-private"}}))
	var doc struct {
		Providers map[string]struct {
			Type, Behavior, Format, URL string
		} `yaml:"rule-providers"`
		Rules []string `yaml:"rules"`
	}
	if err := yaml.Unmarshal([]byte(conf), &doc); err != nil {
		t.Fatalf("invalid YAML: %v\n%s", err, conf)
	}
	want := []string{"RULE-SET,geosite-private,DIRECT", "RULE-SET,geoip-private,DIRECT,no-resolve", "RULE-SET,geosite-ads,REJECT", "RULE-SET,geoip-malware,REJECT,no-resolve",
		"RULE-SET,geosite-sanctioned,PROXY", "RULE-SET,geosite-ir,DIRECT", "RULE-SET,geoip-ir,DIRECT", "MATCH,PROXY"}
	if strings.Join(doc.Rules, "|") != strings.Join(want, "|") {
		t.Fatalf("rules\n got %v\nwant %v", doc.Rules, want)
	}
	if p := doc.Providers["geoip-ir"]; p.Behavior != "ipcidr" || p.Format != "mrs" || p.URL != "https://x/release/geoip-ir.mrs" {
		t.Fatalf("provider %+v", p)
	}
	if doc.Providers["geosite-ir"].Behavior != "domain" {
		t.Fatal("geosite-ir must use domain behavior")
	}
}

func TestV2rayNRouting(t *testing.T) {
	b, err := V2rayNRouting(RouteData{RejectSite: []string{"ads"}, RejectIP: []string{"malware"}, ProxySite: []string{"sanctioned"}})
	if err != nil {
		t.Fatal(err)
	}
	var rules []struct {
		OutboundTag string   `json:"outboundTag"`
		Domain      []string `json:"domain"`
		IP          []string `json:"ip"`
		Port        string   `json:"port"`
	}
	if err := json.Unmarshal(b, &rules); err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range rules {
		got = append(got, r.OutboundTag+":"+strings.Join(append(append(r.Domain, r.IP...), r.Port), ","))
	}
	want := "direct:geoip:private,|block:geosite:ads,|block:geoip:malware,|proxy:geosite:sanctioned,|direct:geosite:ir,|direct:geoip:ir,|proxy:0-65535"
	if strings.Join(got, "|") != want {
		t.Fatalf("got  %s\nwant %s", strings.Join(got, "|"), want)
	}
}
