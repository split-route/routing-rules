---
icon: material/waves
---

# Surge and Loon

Surge and Loon read the same `RULE-SET` list format as Shadowrocket (`DOMAIN-SUFFIX,…`, `IP-CIDR,…,no-resolve`),
so they use the Shadowrocket `.list` files directly. There is no ready-made config for them; add the rules
below to your own config.

!!! warning "Not tested on device"
    These snippets follow each app's documented syntax but haven't been tested in Surge or Loon.

## Surge

The quickest way is the module, which adds all the rules below on top of your config:

```
https://github.com/split-route/routing-rules/releases/latest/download/routing-rules.sgmodule
```

Install it from **Modules** → **Install New Module** with that URL. The rules send sites that block Iranian
IPs to a policy named `PROXY`; create a policy group with that name, or edit the module's copy if yours is
called something else.

To add the rules by hand instead, put these lines in the `[Rule]` section, above `FINAL`:

```
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/ir-injected.list,REJECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/malware.list,REJECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/phishing.list,REJECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/cryptominers.list,REJECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/ads.list,REJECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/sanctioned.list,PROXY
DOMAIN-SUFFIX,ir,DIRECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/ir-domains.list,DIRECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/ir-cidr.list,DIRECT
GEOIP,IR,DIRECT
FINAL,PROXY
```

Replace `PROXY` with the name of your proxy or policy group.

## Loon

Add these lines to the `[Remote Rule]` section:

```
https://github.com/split-route/routing-rules/releases/latest/download/ir-injected.list, policy=REJECT, tag=ir-injected, enabled=true
https://github.com/split-route/routing-rules/releases/latest/download/malware.list, policy=REJECT, tag=malware, enabled=true
https://github.com/split-route/routing-rules/releases/latest/download/phishing.list, policy=REJECT, tag=phishing, enabled=true
https://github.com/split-route/routing-rules/releases/latest/download/cryptominers.list, policy=REJECT, tag=cryptominers, enabled=true
https://github.com/split-route/routing-rules/releases/latest/download/ads.list, policy=REJECT, tag=ads, enabled=true
https://github.com/split-route/routing-rules/releases/latest/download/sanctioned.list, policy=PROXY, tag=sanctioned, enabled=true
https://github.com/split-route/routing-rules/releases/latest/download/ir-domains.list, policy=DIRECT, tag=ir-domains, enabled=true
https://github.com/split-route/routing-rules/releases/latest/download/ir-cidr.list, policy=DIRECT, tag=ir-cidr, enabled=true
```

Then add `DOMAIN-SUFFIX,ir,DIRECT` and `GEOIP,IR,DIRECT` in `[Rule]`, and set `FINAL` to your proxy policy.
Replace `PROXY` with your policy group's name.

## Lighter lists

Surge can load domain lists as `DOMAIN-SET`, which uses less memory than `RULE-SET`:

```
DOMAIN-SET,https://github.com/split-route/routing-rules/releases/latest/download/ads-domainset.txt,REJECT
```

Every domain list has a `-domainset.txt` version; see [File URLs](urls.md#surge-shadowrocket-domain-set).
`ads-lite.list` is a ~5k-rule alternative to the full `ads.list`.
