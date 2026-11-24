---
icon: material/cat
---

# Clash / mihomo

هسته اصلی Clash در سال ۲۰۲۳ متوقف شد. جانشین فعال آن [mihomo](https://github.com/MetaCubeX/mihomo)
(Clash.Meta سابق) است که در Clash Verge Rev، FlClash، Clash Meta for Android و برنامه‌های مشابه استفاده
می‌شود. rule-setهای این پروژه از قالب باینری فشرده mihomo یعنی `.mrs` استفاده می‌کنند.

## استفاده از قوانین آماده

`clash.yaml` بخش‌های `rule-providers` و `rules` را با همان ترتیب کانفیگ Shadowrocket دارد:

[:material-download: دانلود clash.yaml](https://github.com/split-route/routing-rules/releases/latest/download/clash.yaml){ .md-button .md-button--primary }

هر دو بخش را با کانفیگ خود ادغام کنید. در Clash Verge Rev می‌توانید آن‌ها را در فایل **Merge** (کانفیگ
تکمیلی) پروفایل پیست کنید تا اشتراک شما دست‌نخورده بماند. قوانین به یک proxy group با نام **`PROXY`** نیاز
دارند؛ اگر گروه شما نام دیگری دارد، آن را در فایل تغییر دهید.

## DNS

`clash-dns.yaml` یک بلوک `dns` برای mihomo است: fake-ip برای همه دامنه‌ها به‌جز دامنه‌های ایرانی که با DNS
سیستم resolve می‌شوند (`nameserver-policy`)، و DNS-over-HTTPS از طریق گروه `PROXY` برای بقیه. آن را همراه
`clash.yaml` ادغام کنید که rule-provider مربوط به `geosite-ir` را تعریف می‌کند.

[:material-download: دانلود clash-dns.yaml](https://github.com/split-route/routing-rules/releases/latest/download/clash-dns.yaml){ .md-button .md-button--primary }

## یا rule-setها را خودتان اضافه کنید

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

برای `geosite-*` از `behavior: domain` و برای `geoip-*` از `behavior: ipcidr` استفاده کنید. برچسب‌ها همان
برچسب‌های sing-box هستند؛ ببینید [راهنمای sing-box](sing-box.md#tags).

## قوانین IP و DNS

قوانین مسدودسازی و پراکسی که بر اساس IP تطبیق می‌دهند `no-resolve` دارند، پس هیچ دامنه‌ای فقط برای مقایسه با
لیست مسدودسازی resolve نمی‌شود. قانون DIRECT برای `geoip-ir` resolve را انجام می‌دهد، مانند `GEOIP,IR` در
Shadowrocket؛ این کار سایت‌های ایرانی را که در `geosite-ir` نیستند هم پوشش می‌دهد، اما آن درخواست‌های DNS
را به DNS تنظیم‌شده شما می‌فرستد.

!!! note "دیگر برنامه‌های سازگار با Clash"
    `.mrs` مخصوص mihomo است. برنامه‌هایی که پیاده‌سازی Clash خودشان را دارند (مثلاً Stash) ممکن است آن را
    نخوانند. فایل‌های `.list` مربوط به Shadowrocket از نحو کلاسیک `DOMAIN-SUFFIX,…` / `IP-CIDR,…` استفاده
    می‌کنند که providerهای `behavior: classical` در Clash می‌خوانند، اما در آن برنامه‌ها تست نشده‌اند.

همه لینک‌ها: [لینک فایل‌ها](urls.md#clash-mihomo-rule-sets).
