---
icon: material/web
---

# مرورگرها (PAC)

فایل PAC (Proxy Auto-Config) به مرورگر می‌گوید کدام سایت‌ها را از پراکسی عبور دهد. `proxy.pac`
دامنه‌های ایرانی (از جمله همه دامنه‌های `.ir`) را مستقیم و بقیه را از طریق پراکسی محلی کلاینت شما (v2rayN،
Clash، sing-box و مشابه) می‌فرستد.

[:material-download: دانلود proxy.pac](https://github.com/split-route/routing-rules/releases/latest/download/proxy.pac){ .md-button .md-button--primary }

## لینک PAC {#pac-url}

macOS، Windows و Firefox **لینک** فایل PAC را می‌گیرند. این لینک را کپی کنید:

```
https://github.com/split-route/routing-rules/releases/latest/download/proxy.pac
```

این لینک همیشه جدیدترین فایل را می‌دهد، پس لیست‌ها بدون دانلود دوباره آپدیت می‌مانند. از آدرس پراکسی پیش‌فرض
پایین استفاده می‌کند؛ اگر کلاینت شما روی پورت دیگری است، فایل را دانلود و ویرایش کنید و تنظیمات را روی نسخه
خودتان بگذارید.

!!! warning "اگر لینک باز نشود"
    گیت‌هاب دانلود را به `release-assets.githubusercontent.com` می‌فرستد که در بعضی شبکه‌های ایران فیلتر است.
    اگر سیستم نتواند فایل PAC را بگیرد، مستقیم وصل می‌شود و هیچ چیز از پراکسی رد نمی‌شود. اگر یک‌دفعه سایت‌ها
    بدون پراکسی باز شدند، از نسخه دانلودشده استفاده کنید.

## تنظیم آدرس پراکسی

فایل را باز کنید و خط اول را روی پراکسی محلی کلاینت خود تنظیم کنید:

```js
var PROXY = "SOCKS5 127.0.0.1:10808; SOCKS 127.0.0.1:10808";
```

مقدار پیش‌فرض با v2rayN هماهنگ است (SOCKS روی `10808`؛ HTTP روی `10809` که به شکل
`PROXY 127.0.0.1:10809` نوشته می‌شود).

!!! note "بدون بازگشت به اتصال مستقیم"
    این فایل عمداً `DIRECT` را به‌عنوان جایگزین ندارد. اگر کلاینت پراکسی شما متوقف شود، سایت‌ها باز نمی‌شوند
    به‌جای این‌که بی‌صدا بدون محافظت باز شوند.

## استفاده

- **SwitchyOmega / ZeroOmega**: پروفایل جدید ←‏ **PAC Profile** ←‏ محتوای فایل را پیست کنید.
- **macOS**: **System Settings** ←‏ **Network** ←‏ اتصال شما (Wi-Fi یا Ethernet) ←‏ **Details** ←‏
  **Proxies** ←‏ **Automatic proxy configuration** را فعال کنید ←‏ [لینک PAC](#pac-url) را پیست کنید ←‏ **OK**.
- **Windows**: **Settings** ←‏ **Network & internet** ←‏ **Proxy** ←‏ **Use setup script** ←‏
  [لینک PAC](#pac-url) را پیست کنید ←‏ **Save**.
- **Firefox**: تنظیمات ←‏ Network Settings ←‏ **Automatic proxy configuration URL** ←‏ [لینک PAC](#pac-url) را
  پیست کنید، یا مسیر `file://` نسخه ویرایش‌شده.

فایل PAC نمی‌تواند تبلیغات را مسدود کند و فقط با نام دامنه تطبیق می‌دهد: سایت‌های ایرانی که با آدرس IP باز
شوند از پراکسی عبور می‌کنند.
