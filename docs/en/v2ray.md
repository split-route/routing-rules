---
icon: material/alpha-x-circle-outline
---

# v2ray / Xray

`geosite.dat` and `geoip.dat` use the standard v2ray format, so they work in Xray, v2ray, v2rayN, v2rayNG
and other Xray-based apps.

Set these as the custom geo file URLs in your app:

`geosite.dat`

```
https://github.com/split-route/routing-rules/releases/latest/download/geosite.dat
```

`geoip.dat`

```
https://github.com/split-route/routing-rules/releases/latest/download/geoip.dat
```

Opening these URLs in a browser downloads the files.

## Tags

| `geosite:` | `geoip:` |
|------------|----------|
| `ads`, `malware`, `phishing`, `cryptominers`, `sanctioned`, `ir`, `private`, `social`, `nsfw`, `tracking` | `ir`, `malware`, `phishing`, `private`, `cloudflare`, `cloudfront`, `fastly`, `gcore`, `telegram`, `tor`, `netflix` |

!!! warning "These files only contain the tags above"
    If you replace your app's stock geo files with these, rules that use other tags (`geosite:google`,
    `geoip:cn`, `geosite:category-ads-all`) stop working. Either update those rules, or keep the stock files
    and load these under another name (Xray only):

    ```
    ext:iran-geosite.dat:ads
    ext:iran-geoip.dat:ir
    ```

## Example routing

`xray-routing.json` follows the same order as the Shadowrocket config:

[:material-download: Download xray-routing.json](https://github.com/split-route/routing-rules/releases/latest/download/xray-routing.json){ .md-button .md-button--primary }

```json
{
  "routing": {
    "domainStrategy": "IPIfNonMatch",
    "rules": [
      { "type": "field", "ip": ["geoip:ir-injected"], "outboundTag": "block" },
      { "type": "field", "domain": ["geosite:private"], "outboundTag": "direct" },
      { "type": "field", "ip": ["geoip:private"], "outboundTag": "direct" },
      { "type": "field", "domain": ["geosite:malware", "geosite:phishing", "geosite:cryptominers", "geosite:ads"], "outboundTag": "block" },
      { "type": "field", "ip": ["geoip:ir-injected", "geoip:malware", "geoip:phishing"], "outboundTag": "block" },
      { "type": "field", "domain": ["geosite:sanctioned"], "outboundTag": "proxy" },
      { "type": "field", "domain": ["geosite:ir"], "outboundTag": "direct" },
      { "type": "field", "ip": ["geoip:ir"], "outboundTag": "direct" },
      { "type": "field", "network": "tcp,udp", "outboundTag": "proxy" }
    ]
  }
}
```

It expects outbounds tagged `proxy`, `direct` and `block` (blackhole). The first rule blocks Iran's DNS
block-page addresses before the private-IP rule, because they sit inside `10.0.0.0/8`. `IPIfNonMatch` resolves domains that
no domain rule matched, so `geoip:ir` also catches Iranian sites missing from `geosite:ir`.

## DNS

`xray-dns.json` is a `dns` block: Iranian domains (`geosite:ir`) use the system resolver, everything else
DNS-over-HTTPS (`https://1.1.1.1/dns-query`), which the routing above sends through the proxy.

[:material-download: Download xray-dns.json](https://github.com/split-route/routing-rules/releases/latest/download/xray-dns.json){ .md-button .md-button--primary }

## v2rayN / v2rayNG

These apps also need custom routing rules that use the tags. A ready-made rules file and import steps are
in the [v2rayN / v2rayNG guide](v2rayn.md).
