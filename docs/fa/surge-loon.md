---
icon: material/waves
---

# Surge و Loon

Surge و Loon همان قالب لیست `RULE-SET` مربوط به Shadowrocket را می‌خوانند (`DOMAIN-SUFFIX,…`،
`IP-CIDR,…,no-resolve`)، پس مستقیم از فایل‌های `.list` مربوط به Shadowrocket استفاده می‌کنند. کانفیگ
آماده‌ای برای آن‌ها وجود ندارد؛ قوانین زیر را به کانفیگ خودتان اضافه کنید.

!!! warning "روی دستگاه تست نشده"
    این قطعه‌ها طبق نحو مستند هر برنامه نوشته شده‌اند، اما در Surge یا Loon تست نشده‌اند.

## Surge

سریع‌ترین راه، ماژول است که همه قوانین زیر را روی کانفیگ شما اضافه می‌کند:

```
https://github.com/split-route/routing-rules/releases/latest/download/routing-rules.sgmodule
```

آن را از **Modules** ←‏ **Install New Module** با همین لینک نصب کنید. قوانین، سایت‌هایی را که IP ایران را
مسدود می‌کنند به policy با نام `PROXY` می‌فرستند؛ یک policy group با همین نام بسازید، یا اگر نام گروه شما فرق
دارد، نسخه‌ای از ماژول را ویرایش کنید.

برای افزودن دستی قوانین، این خطوط را در بخش `[Rule]`، بالای `FINAL`، اضافه کنید:

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

به‌جای `PROXY` نام پراکسی یا policy group خود را بگذارید.

## Loon

این خطوط را در بخش `[Remote Rule]` اضافه کنید:

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

سپس `DOMAIN-SUFFIX,ir,DIRECT` و `GEOIP,IR,DIRECT` را در `[Rule]` اضافه کنید و `FINAL` را روی policy پراکسی
خود بگذارید. به‌جای `PROXY` نام policy group خود را بگذارید.

## لیست‌های سبک‌تر

Surge می‌تواند لیست‌های دامنه را به شکل `DOMAIN-SET` لود کند که حافظه کمتری از `RULE-SET` مصرف می‌کند:

```
DOMAIN-SET,https://github.com/split-route/routing-rules/releases/latest/download/ads-domainset.txt,REJECT
```

هر لیست دامنه یک نسخه `-domainset.txt` دارد؛ ببینید [لینک فایل‌ها](urls.md#surge-shadowrocket-domain-set).
`ads-lite.list` جایگزینی حدود ۵ هزار قانونی برای `ads.list` کامل است.
