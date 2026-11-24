package emit

import (
	"bytes"
	"fmt"
	"net/netip"
	"strings"

	"github.com/split-route/routing-rules/internal/parse"
)

// WireGuardAllowedIPs renders a comma-separated AllowedIPs value: every
// address except the excluded prefixes (callers pass the complement).
func WireGuardAllowedIPs(allowed []netip.Prefix) []byte {
	s := make([]string, 0, len(allowed))
	for _, p := range allowed {
		s = append(s, p.String())
	}
	return []byte(strings.Join(s, ", ") + "\n")
}

// MikroTik renders a RouterOS script that replaces the address list `list`
// (IPv4 and IPv6) with the given prefixes.
func MikroTik(list string, cidrs []netip.Prefix) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# routing-rules: Iran IP ranges as the RouterOS address list %q.\n# Re-running replaces the list.\n", list)
	for _, v4 := range []bool{true, false} {
		menu := "/ipv6 firewall address-list"
		if v4 {
			menu = "/ip firewall address-list"
		}
		fmt.Fprintf(&b, "%s remove [find list=%q]\n%s\n", menu, list, menu)
		for _, p := range cidrs {
			if p.Addr().Is4() == v4 {
				fmt.Fprintf(&b, "add list=%s address=%s\n", list, p)
			}
		}
	}
	return b.Bytes()
}

// CIDRLines renders one prefix per line for one address family (ipset, nftables).
func CIDRLines(cidrs []netip.Prefix, v4 bool) []byte {
	var b bytes.Buffer
	for _, p := range cidrs {
		if p.Addr().Is4() == v4 {
			b.WriteString(p.String() + "\n")
		}
	}
	return b.Bytes()
}

// DomainSet renders a Surge-style DOMAIN-SET: ".x" matches x and its
// subdomains, "x" only x. Keywords have no DOMAIN-SET form and are skipped.
func DomainSet(domains []parse.Domain) []byte {
	var b bytes.Buffer
	for _, d := range domains {
		switch d.Kind {
		case parse.Suffix:
			b.WriteString("." + d.Name + "\n")
		case parse.Exact:
			b.WriteString(d.Name + "\n")
		}
	}
	return b.Bytes()
}

// AdGuard renders an AdGuard Home / Pi-hole v6 blocklist. Every entry uses
// "||x^" (x and its subdomains), the form both understand; exact entries are
// widened to their subdomains, as DNS blocklists usually are.
func AdGuard(title string, domains []parse.Domain) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "! Title: routing-rules %s\n! Homepage: https://github.com/split-route/routing-rules\n", title)
	for _, d := range domains {
		if d.Kind != parse.Keyword {
			b.WriteString("||" + d.Name + "^\n")
		}
	}
	return b.Bytes()
}

// Hosts renders a hosts-file blocklist (Pi-hole v5, /etc/hosts). A hosts file
// cannot block subdomains, so suffix entries only block the domain itself.
func Hosts(title string, domains []parse.Domain) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# Title: routing-rules %s (hosts format; subdomains of listed domains are not blocked)\n", title)
	for _, d := range domains {
		if d.Kind != parse.Keyword {
			b.WriteString("0.0.0.0 " + d.Name + "\n")
		}
	}
	return b.Bytes()
}

// PAC renders a proxy auto-config file: the listed domains (and their
// subdomains, for suffix rules) go DIRECT, everything else to the proxy.
func PAC(direct []parse.Domain) []byte {
	var suffix, exact []string
	for _, d := range direct {
		switch d.Kind {
		case parse.Suffix:
			suffix = append(suffix, fmt.Sprintf("%q:1", d.Name))
		case parse.Exact:
			exact = append(exact, fmt.Sprintf("%q:1", d.Name))
		}
	}
	var b bytes.Buffer
	b.WriteString(`// routing-rules PAC: Iranian domains go direct, everything else through your proxy.
// Change PROXY to your client's local proxy (v2rayN defaults: SOCKS 10808, HTTP 10809).
// No DIRECT fallback on purpose: if the proxy stops, browsing fails instead of silently
// going out unprotected.
var PROXY = "SOCKS5 127.0.0.1:10808; SOCKS 127.0.0.1:10808";

`)
	fmt.Fprintf(&b, "var SUFFIX = {%s};\nvar EXACT = {%s};\n", strings.Join(suffix, ","), strings.Join(exact, ","))
	b.WriteString(`
function FindProxyForURL(url, host) {
  host = host.toLowerCase().replace(/\.$/, "");
  if (isPlainHostName(host) || /^(127\.|10\.|192\.168\.|172\.(1[6-9]|2[0-9]|3[01])\.|169\.254\.)/.test(host)) {
    return "DIRECT";
  }
  if (Object.prototype.hasOwnProperty.call(EXACT, host)) {
    return "DIRECT";
  }
  for (var h = host; ; ) {
    if (Object.prototype.hasOwnProperty.call(SUFFIX, h)) {
      return "DIRECT";
    }
    var i = h.indexOf(".");
    if (i < 0) {
      break;
    }
    h = h.substring(i + 1);
  }
  return PROXY;
}
`)
	return b.Bytes()
}

// SingBoxDNS renders a "dns" block (sing-box 1.12+ server format): local names
// (geosite-private) and Iranian domains use the device's DNS, everything else DNS-over-HTTPS through the proxy.
// route.default_domain_resolver (required since 1.14) resolves outbound server
// hostnames with the device's DNS; the DoH server is only reachable through the
// proxy, so it cannot resolve the proxy's own address.
func SingBoxDNS() []byte {
	return []byte(`{
  "route": {
    "default_domain_resolver": "local"
  },
  "dns": {
    "servers": [
      { "type": "https", "tag": "remote", "server": "1.1.1.1", "detour": "proxy" },
      { "type": "local", "tag": "local" }
    ],
    "rules": [
      { "rule_set": ["geosite-private", "geosite-ir"], "server": "local" }
    ],
    "final": "remote"
  }
}
`)
}

// XrayDNS renders a "dns" block: Iranian domains use the system resolver,
// everything else DNS-over-HTTPS (routed through the proxy by the routing rules).
func XrayDNS() []byte {
	return []byte(`{
  "dns": {
    "servers": [
      "https://1.1.1.1/dns-query",
      { "address": "localhost", "domains": ["geosite:private", "geosite:ir"], "skipFallback": true }
    ],
    "queryStrategy": "UseIP"
  }
}
`)
}

// ClashDNS renders a mihomo "dns" block: fake-ip for everything except
// Iranian domains, which resolve through the system DNS; other lookups use
// DNS-over-HTTPS through the PROXY group.
func ClashDNS() []byte {
	return []byte(`# routing-rules DNS for mihomo (Clash.Meta). Merge together with clash.yaml,
# which defines the geosite-ir rule-provider. Requires a proxy group named PROXY.
dns:
  enable: true
  ipv6: true
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  fake-ip-filter:
    - "rule-set:geosite-private"
    - "rule-set:geosite-ir"
    - "+.lan"
    - "+.local"
  nameserver:
    - "https://1.1.1.1/dns-query#PROXY"
    - "https://8.8.8.8/dns-query#PROXY"
  nameserver-policy:
    "rule-set:geosite-private": system
    "rule-set:geosite-ir": system
  proxy-server-nameserver:
    - system
`)
}

// OpenVPN renders client-config lines that route the given IPv4 prefixes
// outside the tunnel ("route <net> <mask> net_gateway"). They are more
// specific than the server's redirect-gateway routes, so they take priority.
func OpenVPN(cidrs []netip.Prefix) []byte {
	var b bytes.Buffer
	b.WriteString("# routing-rules: send Iranian IPv4 ranges outside the OpenVPN tunnel.\n# Paste into the client .ovpn file. IPv6 is not covered and stays in the tunnel.\n")
	for _, p := range cidrs {
		if !p.Addr().Is4() {
			continue
		}
		mask := ^uint32(0) << (32 - p.Bits())
		if p.Bits() == 0 {
			mask = 0
		}
		fmt.Fprintf(&b, "route %s %d.%d.%d.%d net_gateway\n", p.Addr(), mask>>24, mask>>16&0xff, mask>>8&0xff, mask&0xff)
	}
	return b.Bytes()
}
