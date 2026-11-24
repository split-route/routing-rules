---
icon: material/cube-outline
---

# sing-box

Rule-sets use format version 1, so they work with sing-box 1.8+ and apps built on older cores (Karing,
Hiddify, NekoBox…). The ready-made route block below uses rule actions and needs sing-box 1.11+. Rule-sets
are named like [Chocolate4U/Iran-sing-box-rules](https://github.com/Chocolate4U/Iran-sing-box-rules):
`geosite-<tag>` and `geoip-<tag>`.

## Use the ready route

`sing-box-route.json` is a complete `route` block, in the same order as the Shadowrocket config:

[:material-download: Download sing-box-route.json](https://github.com/split-route/routing-rules/releases/latest/download/sing-box-route.json){ .md-button .md-button--primary }

Copy its `route` object into your config. It expects outbounds tagged `proxy` and `direct`. Rule-sets are
downloaded through `proxy` (`download_detour`) and refreshed automatically.

## DNS

`sing-box-dns.json` is a `dns` block (sing-box 1.12+ server format) that pairs with the route above:
Iranian domains (`geosite-ir`) resolve through the device's own DNS, so Iranian sites get their domestic
addresses, and everything else uses DNS-over-HTTPS (`1.1.1.1`) through the proxy, so lookups are neither
leaked nor tampered with.

[:material-download: Download sing-box-dns.json](https://github.com/split-route/routing-rules/releases/latest/download/sing-box-dns.json){ .md-button .md-button--primary }

Merge it into your config: it adds a `dns` block and sets `route.default_domain_resolver` (required since
sing-box 1.14) so proxy server hostnames resolve with the device's DNS. It uses the `geosite-ir` rule-set the
route defines and the `proxy` outbound.

## Or reference rule-sets yourself

```json
{
  "route": {
    "rule_set": [
      {
        "tag": "geosite-ads",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/split-route/routing-rules/releases/latest/download/geosite-ads.srs",
        "download_detour": "proxy"
      },
      {
        "tag": "geosite-ir",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/split-route/routing-rules/releases/latest/download/geosite-ir.srs",
        "download_detour": "proxy"
      }
    ],
    "rules": [
      { "rule_set": "geosite-ads", "action": "reject" },
      { "rule_set": "geosite-ir", "outbound": "direct" }
    ],
    "final": "proxy"
  }
}
```

Every `.srs` file also exists as `.json` for `"format": "source"`.

## Tags

| Domains (`geosite-…`) | IP ranges (`geoip-…`) |
|------------------------|------------------------|
| `ads`, `malware`, `phishing`, `cryptominers`, `sanctioned`, `ir`, `private`, `social`, `nsfw`, `tracking` | `ir`, `malware`, `phishing`, `private`, `cloudflare`, `cloudfront`, `fastly`, `gcore`, `telegram`, `tor`, `netflix` |

`geosite-ir` already includes the `.ir` and `.ایران` TLDs.

## Matching Iranian IPs

`geoip-ir` only matches connections that go to an IP address, like `no-resolve` in Shadowrocket. If you
also want domains resolved and matched by IP, add `{ "action": "resolve" }` before the `geoip-ir` rule. That
sends those DNS lookups through your configured DNS server.

All rule-set URLs: [File URLs](urls.md#sing-box-rule-sets-binary).
