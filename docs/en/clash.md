---
icon: material/cat
---

# Clash / mihomo

The original Clash core was discontinued in 2023. Its maintained successor is
[mihomo](https://github.com/MetaCubeX/mihomo) (formerly Clash.Meta), which powers Clash Verge Rev, FlClash,
Clash Meta for Android and similar apps. The rule-sets here use mihomo's compact binary format, `.mrs`.

## Use the ready rules

`clash.yaml` contains `rule-providers` and `rules` in the same order as the Shadowrocket config:

[:material-download: Download clash.yaml](https://github.com/split-route/routing-rules/releases/latest/download/clash.yaml){ .md-button .md-button--primary }

Merge both sections into your config. In Clash Verge Rev you can paste them into the profile's **Merge**
(extend config) file, so your subscription stays untouched. The rules need a proxy group named **`PROXY`**;
rename it in the file if your group has a different name.

## DNS

`clash-dns.yaml` is a `dns` block for mihomo: fake-ip for everything except Iranian domains, which resolve
through the system DNS (`nameserver-policy`), and DNS-over-HTTPS through the `PROXY` group for the rest.
Merge it together with `clash.yaml`, which defines the `geosite-ir` rule-provider it uses.

[:material-download: Download clash-dns.yaml](https://github.com/split-route/routing-rules/releases/latest/download/clash-dns.yaml){ .md-button .md-button--primary }

## Or add rule-sets yourself

```yaml
rule-providers:
  geosite-ads:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/split-route/routing-rules/releases/latest/download/geosite-ads.mrs
    path: ./rules/geosite-ads.mrs
    interval: 86400
  geoip-ir:
    type: http
    behavior: ipcidr
    format: mrs
    url: https://github.com/split-route/routing-rules/releases/latest/download/geoip-ir.mrs
    path: ./rules/geoip-ir.mrs
    interval: 86400

rules:
  - RULE-SET,geosite-ads,REJECT
  - RULE-SET,geoip-ir,DIRECT
  - MATCH,PROXY
```

Use `behavior: domain` for `geosite-*` and `behavior: ipcidr` for `geoip-*`. The tags are the same as for
sing-box; see the [sing-box guide](sing-box.md#tags).

## IP rules and DNS

Block and proxy rules that match by IP use `no-resolve`, so a domain is never looked up just to check it
against a blocklist. The `geoip-ir` DIRECT rule does resolve, like Shadowrocket's `GEOIP,IR`. That catches
Iranian sites missing from `geosite-ir`, but sends those DNS lookups through your configured DNS.

!!! note "Other Clash-compatible apps"
    `.mrs` is mihomo-only. Apps with their own Clash implementation (Stash, for example) may not read it.
    The Shadowrocket `.list` files use the classic `DOMAIN-SUFFIX,…` / `IP-CIDR,…` syntax that Clash's
    `behavior: classical` providers read, but they haven't been tested in those apps.

All URLs: [File URLs](urls.md#clash-mihomo-rule-sets).
