---
icon: material/alpha-v-circle-outline
---

# v2rayN / v2rayNG

`v2rayn-routing.json` مجموعه‌ای از قوانین مسیریابی سفارشی در قالب JSON است که v2rayN و v2rayNG ایمپورت
می‌کنند. از برچسب‌های `geosite:` / `geoip:` این پروژه استفاده می‌کند.

[:material-download: دانلود v2rayn-routing.json](https://github.com/split-route/routing-rules/releases/latest/download/v2rayn-routing.json){ .md-button .md-button--primary }

برای ایمپورت با لینک، از این استفاده کنید:

```
https://github.com/split-route/routing-rules/releases/latest/download/v2rayn-routing.json
```

## ۱. اول فایل‌های geo را نصب کنید

قوانین به برچسب‌هایی مانند `geosite:ads` و `geoip:ir` ارجاع می‌دهند که فقط در `geosite.dat` / `geoip.dat`
این پروژه وجود دارند. قبل از ایمپورت قوانین، آن‌ها را طبق [راهنمای v2ray / Xray](v2ray.md) تنظیم کنید؛ در
غیر این صورت قوانین لود نمی‌شوند یا با هیچ چیز تطبیق نمی‌یابند.

## ۲. ایمپورت قوانین

تنظیمات مسیریابی برنامه را باز کنید، یک مجموعه قانون سفارشی جدید بسازید و فایل دانلودشده (یا لینک بالا) را
ایمپورت کنید. سپس آن مجموعه را به‌عنوان مسیریابی فعال انتخاب کنید.

!!! note "نکته"
    نام منوها در نسخه‌ها و سیستم‌عامل‌های مختلف فرق دارد. دنبال **Routing setting** و گزینه **Import rules**
    بگردید.

## محتوا

| # | تطبیق | خروجی |
|---|-------|-------|
| ۱ | `geoip:ir-injected` (آدرس‌های صفحه مسدودسازی DNS ایران) | block |
| ۲ | `geosite:private` (نام‌های محلی: `.lan`، `.local`، `localhost`، صفحه روترها) | direct |
| ۳ | `geoip:private` | direct |
| ۴ | `geosite:malware`، `geosite:phishing`، `geosite:cryptominers`، `geosite:ads` | block |
| ۵ | `geoip:malware`، `geoip:phishing` | block |
| ۶ | `geosite:sanctioned` | proxy |
| ۷ | `geosite:ir` | direct |
| ۸ | `geoip:ir` | direct |
| ۹ | بقیه | proxy |

برچسب‌های خروجی همان `proxy`، `direct` و `block` استاندارد هر دو برنامه هستند.
