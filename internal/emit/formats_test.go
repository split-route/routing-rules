package emit

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/split-route/routing-rules/internal/parse"
)

var sample = []parse.Domain{{Kind: parse.Suffix, Name: "a.ir"}, {Kind: parse.Exact, Name: "b.com"}, {Kind: parse.Keyword, Name: "ad"}}

func TestTextFormats(t *testing.T) {
	cases := map[string]struct{ got, want string }{
		"domainset": {string(DomainSet(sample)), ".a.ir\nb.com\n"},
		"adguard":   {string(AdGuard("x", sample)), "! Title: routing-rules x\n! Homepage: https://github.com/split-route/routing-rules\n||a.ir^\n||b.com^\n"},
		"hosts":     {string(Hosts("x", sample)), "# Title: routing-rules x (hosts format; subdomains of listed domains are not blocked)\n0.0.0.0 a.ir\n0.0.0.0 b.com\n"},
		"wireguard": {string(WireGuardAllowedIPs([]netip.Prefix{netip.MustParsePrefix("1.0.0.0/8"), netip.MustParsePrefix("::/0")})), "1.0.0.0/8, ::/0\n"},
		"openvpn": {string(OpenVPN([]netip.Prefix{netip.MustParsePrefix("2.144.0.0/14"), netip.MustParsePrefix("5.22.0.0/23"), netip.MustParsePrefix("2a01::/32")})),
			"# routing-rules: send Iranian IPv4 ranges outside the OpenVPN tunnel.\n# Paste into the client .ovpn file. IPv6 is not covered and stays in the tunnel.\n" +
				"route 2.144.0.0 255.252.0.0 net_gateway\nroute 5.22.0.0 255.255.254.0 net_gateway\n"},
		"cidr4": {string(CIDRLines([]netip.Prefix{netip.MustParsePrefix("2.144.0.0/14"), netip.MustParsePrefix("2a01::/32")}, true)), "2.144.0.0/14\n"},
	}
	for name, c := range cases {
		if c.got != c.want {
			t.Errorf("%s:\n got %q\nwant %q", name, c.got, c.want)
		}
	}
	mt := string(MikroTik("IR", []netip.Prefix{netip.MustParsePrefix("2.144.0.0/14"), netip.MustParsePrefix("2a01::/32")}))
	for _, want := range []string{"/ip firewall address-list remove [find list=\"IR\"]\n/ip firewall address-list\nadd list=IR address=2.144.0.0/14\n",
		"/ipv6 firewall address-list remove [find list=\"IR\"]\n/ipv6 firewall address-list\nadd list=IR address=2a01::/32\n"} {
		if !strings.Contains(mt, want) {
			t.Errorf("mikrotik missing %q in\n%s", want, mt)
		}
	}
	pac := string(PAC(sample))
	if !strings.Contains(pac, `var SUFFIX = {"a.ir":1};`) || !strings.Contains(pac, `var EXACT = {"b.com":1};`) || !strings.Contains(pac, "function FindProxyForURL") {
		t.Errorf("pac:\n%s", pac[:300])
	}
}

func TestEarlyRejectBeforePrivate(t *testing.T) {
	d := RouteData{BaseURL: "https://x", Early: []string{"geoip-ir-injected"}, EarlyIP: []string{"ir-injected"},
		Private: []string{"geosite-private", "geoip-private"}, PrivateSite: []string{"private"}, PrivateIP: []string{"private"},
		Reject: []string{"geosite-ads"}, RejectSite: []string{"ads"}, Direct: []string{"geosite-ir"},
		All: []string{"geoip-ir-injected", "geosite-ads", "geosite-ir"}}
	order := func(name, s, early, private string) {
		e, p := strings.Index(s, early), strings.Index(s, private)
		if e < 0 || p < 0 || e > p {
			t.Errorf("%s: %q must come before %q:\n%s", name, early, private, s)
		}
	}
	sb, _ := SingBoxRoute(d)
	order("sing-box", string(sb), `"rule_set": ["geoip-ir-injected"], "action": "reject"`, `"ip_is_private"`)
	order("sing-box private before ads", string(sb), `"rule_set": ["geosite-private", "geoip-private"], "outbound": "direct"`, `"geosite-ads"`)
	xr, _ := XrayRouting(d)
	order("xray", string(xr), `"geoip:ir-injected"`, `"geosite:private"`)
	order("xray private before ads", string(xr), `"geosite:private"`, `"geosite:ads"`)
	vn, _ := V2rayNRouting(d)
	order("v2rayn", string(vn), `"geoip:ir-injected"`, `"geosite:private"`)
	order("clash", string(ClashConfig("https://x", d)), "RULE-SET,geoip-ir-injected,REJECT", "RULE-SET,geoip-private,DIRECT")
}
