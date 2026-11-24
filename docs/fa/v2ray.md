---
icon: material/alpha-x-circle-outline
---

# v2ray / Xray

`geosite.dat` و `geoip.dat` قالب استاندارد v2ray را دارند، پس در Xray، v2ray، v2rayN، v2rayNG و دیگر
برنامه‌های مبتنی بر Xray کار می‌کنند.

این لینک‌ها را به‌عنوان لینک سفارشی فایل‌های geo در برنامه خود تنظیم کنید:

`geosite.dat`

```
https://github.com/split-route/routing-rules/releases/latest/download/geosite.dat
```

`geoip.dat`

```
https://github.com/split-route/routing-rules/releases/latest/download/geoip.dat
```

باز کردن این لینک‌ها در مرورگر، فایل‌ها را دانلود می‌کند.

## برچسب‌ها

| `geosite:` | `geoip:` |
|------------|----------|
| `ads`، `malware`، `phishing`، `cryptominers`، `sanctioned`، `ir`، `private`، `social`، `nsfw`، `tracking` | `ir`، `malware`، `phishing`، `private`، `cloudflare`، `cloudfront`، `fastly`، `gcore`، `telegram`، `tor`، `netflix` |

!!! warning "این فایل‌ها فقط برچسب‌های بالا را دارند"
    اگر فایل‌های geo پیش‌فرض برنامه را با این‌ها جایگزین کنید، قوانینی که از برچسب‌های دیگر استفاده می‌کنند
    (`geosite:google`، `geoip:cn`، `geosite:category-ads-all`) دیگر کار نمی‌کنند. یا آن قوانین را عوض کنید،
    یا فایل‌های پیش‌فرض را نگه دارید و این‌ها را با نام دیگری لود کنید (فقط Xray):

    ```
    ext:iran-geosite.dat:ads
    ext:iran-geoip.dat:ir
    ```

## نمونه مسیریابی

`xray-routing.json` همان ترتیب کانفیگ Shadowrocket را دارد:

[:material-download: دانلود xray-routing.json](https://github.com/split-route/routing-rules/releases/latest/download/xray-routing.json){ .md-button .md-button--primary }

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

به outboundهایی با برچسب `proxy`، `direct` و `block` (blackhole) نیاز دارد. قانون اول آدرس‌های صفحه
مسدودسازی DNS ایران را قبل از قانون IP خصوصی مسدود می‌کند، چون در بازه `10.0.0.0/8` قرار دارند.
`IPIfNonMatch` دامنه‌هایی را که با هیچ قانون دامنه‌ای تطبیق پیدا نکردند resolve می‌کند، پس `geoip:ir` سایت‌های
ایرانی‌ای را که در `geosite:ir` نیستند هم پوشش می‌دهد.

## DNS

`xray-dns.json` یک بلوک `dns` است: دامنه‌های ایرانی (`geosite:ir`) از DNS سیستم و بقیه از DNS-over-HTTPS
(`https://1.1.1.1/dns-query`) استفاده می‌کنند که مسیریابی بالا آن را از پراکسی عبور می‌دهد.

[:material-download: دانلود xray-dns.json](https://github.com/split-route/routing-rules/releases/latest/download/xray-dns.json){ .md-button .md-button--primary }

## v2rayN / v2rayNG

این برنامه‌ها به قوانین مسیریابی سفارشی هم نیاز دارند که از این برچسب‌ها استفاده کنند. فایل قوانین آماده و
مراحل ایمپورت در [راهنمای v2rayN / v2rayNG](v2rayn.md) است.
