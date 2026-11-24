---
icon: material/cube-outline
---

# sing-box

rule-setها با نسخه ۱ قالب ساخته می‌شوند، پس با sing-box نسخه ۱.۸ به بعد و برنامه‌هایی با هسته‌های قدیمی‌تر
(Karing، Hiddify، NekoBox و …) کار می‌کنند. بلوک route آماده پایین از rule action استفاده می‌کند و به
sing-box نسخه ۱.۱۱ به بعد نیاز دارد. نام rule-setها مانند
[Chocolate4U/Iran-sing-box-rules](https://github.com/Chocolate4U/Iran-sing-box-rules) است:
`geosite-<tag>` و `geoip-<tag>`.

## استفاده از route آماده

`sing-box-route.json` یک بلوک `route` کامل با همان ترتیب کانفیگ Shadowrocket است:

[:material-download: دانلود sing-box-route.json](https://github.com/split-route/routing-rules/releases/latest/download/sing-box-route.json){ .md-button .md-button--primary }

شیء `route` آن را در کانفیگ خود کپی کنید. این بلوک به outboundهایی با برچسب `proxy` و `direct` نیاز
دارد. rule-setها از طریق `proxy` دانلود (`download_detour`) و خودکار آپدیت می‌شوند.

## DNS

`sing-box-dns.json` یک بلوک `dns` (قالب سرور sing-box نسخه ۱.۱۲ به بعد) است که با route بالا کار می‌کند:
دامنه‌های ایرانی (`geosite-ir`) با DNS خود دستگاه resolve می‌شوند تا سایت‌های ایرانی آدرس داخلی بگیرند، و بقیه
از DNS-over-HTTPS (`1.1.1.1`) از طریق پراکسی استفاده می‌کنند تا درخواست‌ها نه نشت کنند و نه دستکاری شوند.

[:material-download: دانلود sing-box-dns.json](https://github.com/split-route/routing-rules/releases/latest/download/sing-box-dns.json){ .md-button .md-button--primary }

آن را با کانفیگ خود ادغام کنید: یک بلوک `dns` اضافه می‌کند و `route.default_domain_resolver` را تنظیم
می‌کند (از sing-box ۱.۱۴ الزامی است) تا نام دامنه سرورهای پراکسی با DNS دستگاه resolve شود. از rule-set
`geosite-ir` که route تعریف می‌کند و outbound `proxy` استفاده می‌کند.

## یا rule-setها را خودتان اضافه کنید

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

هر فایل `.srs` به شکل `.json` هم برای `"format": "source"` موجود است.

## برچسب‌ها { #tags }

| دامنه‌ها (`geosite-…`) | بازه‌های IP (`geoip-…`) |
|------------------------|------------------------|
| `ads`، `malware`، `phishing`، `cryptominers`، `sanctioned`، `ir`، `private`، `social`، `nsfw`، `tracking` | `ir`، `malware`، `phishing`، `private`، `cloudflare`، `cloudfront`، `fastly`، `gcore`، `telegram`، `tor`، `netflix` |

`geosite-ir` دامنه‌های `.ir` و `.ایران` را هم شامل می‌شود.

## تطبیق IPهای ایران

`geoip-ir` فقط با اتصال‌هایی که مستقیم به یک آدرس IP می‌روند تطبیق پیدا می‌کند، مانند `no-resolve` در
Shadowrocket. اگر می‌خواهید دامنه‌ها هم به IP resolve و بر اساس آن تطبیق داده شوند، قبل از قانون `geoip-ir`
قانون `{ "action": "resolve" }` را اضافه کنید؛ این کار آن درخواست‌های DNS را به سرور DNS تنظیم‌شده شما
می‌فرستد.

لینک همه rule-setها: [لینک فایل‌ها](urls.md#sing-box-rule-sets-binary).
