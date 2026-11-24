---
icon: material/rocket-launch-outline
---

# Shadowrocket

## نصب با یک کلیک

روی آیفون یا آیپدی که Shadowrocket روی آن نصب است، روی دکمه بزنید. Shadowrocket باز می‌شود، کانفیگ را
دانلود می‌کند و به آن تغییر می‌دهد.

[:material-rocket-launch: نصب در Shadowrocket](shadowrocket://config/add/https://github.com/split-route/routing-rules/releases/latest/download/shadowrocket.conf){ .md-button .md-button--primary }

سپس در تب **Home** سرور خود را اضافه کنید و **Global Routing** را روی **Config** بگذارید.

!!! note "نکته"
    این دکمه فقط روی دستگاهی کار می‌کند که Shadowrocket روی آن نصب است. کانفیگ در حال استفاده را جایگزین
    می‌کند؛ برای نگه داشتن کانفیگ خودتان از [ماژول](#keep-your-own-config-module) استفاده کنید. اگر اتفاقی
    نیفتاد، مراحل دستی زیر را انجام دهید.

## ایمپورت دستی

1. لینک کانفیگ را کپی کنید:

    ```
    https://github.com/split-route/routing-rules/releases/latest/download/shadowrocket.conf
    ```

2. Shadowrocket ←‏ تب **Config** ←‏ **+** ←‏ لینک را پیست کنید ←‏ **Download**.
3. روی کانفیگ جدید بزنید ←‏ **Use Config**. تیک نشان‌دهنده کانفیگ در حال استفاده است.

    ![تب Config: روی + بزنید، لینک را پیست کنید، Download را بزنید، سپس روی کانفیگ بزنید و Use Config را انتخاب کنید](../assets/img/shadowrocket-import.webp){ width="720" loading=lazy }

4. در تب **Home** سرور خود را اضافه کنید و **Global Routing** را روی **Config** بگذارید.

    ![تب Home: گزینه Global Routing روی Config](../assets/img/shadowrocket-routing.webp){ width="480" loading=lazy }

### آپدیت نگه داشتن

لیست‌ها هر روز تغییر می‌کنند. کانفیگ `update-url` دارد و همه لیست‌ها `RULE-SET` راه دور هستند، پس
لازم نیست دوباره ایمپورت کنید. برای آپدیت فوری، در تب **Config** انگشتتان را چند لحظه روی کانفیگ نگه دارید ←‏ **Update**.

![انگشتتان را روی کانفیگ نگه دارید و Update را بزنید](../assets/img/shadowrocket-update.webp){ width="300" loading=lazy }

## کانفیگ خودتان را نگه دارید (ماژول) {#keep-your-own-config-module}

اگر کانفیگ دلخواه خودتان را دارید (سرورها، قوانین یا تنظیمات خودتان)، به‌جای ایمپورت کانفیگ بالا این
قوانین را به‌صورت **ماژول** اضافه کنید. ماژول همان قوانین را بدون قانون پایانی دارد، پس کانفیگ شما سرورها و
`FINAL` خودش را نگه می‌دارد.

[:material-puzzle-outline: نصب ماژول در Shadowrocket](shadowrocket://install?module=https://github.com/split-route/routing-rules/releases/latest/download/routing-rules.sgmodule){ .md-button }

یا آن را دستی اضافه کنید:

1. لینک ماژول را کپی کنید:

    ```
    https://github.com/split-route/routing-rules/releases/latest/download/routing-rules.sgmodule
    ```

2. Shadowrocket ←‏ تب **Config** ←‏ **Modules** ←‏ **+** ←‏ لینک را پیست کنید ←‏ **Download**.
3. مطمئن شوید ماژول فعال است.

قوانین برای سایت‌هایی که IP ایران را مسدود می‌کنند از policy با نام `PROXY` (پراکسی داخلی Shadowrocket)
استفاده می‌کنند. ماژول‌ها هم مثل کانفیگ، خودکار آپدیت می‌شوند؛ برای آپدیت فوری، ماژول را از صفحه
**Modules** آپدیت کنید.

## ترتیب قوانین

| # | قانون | عمل |
|---|-------|-----|
| ۱ | `ir-injected.list` (آدرس‌های صفحه مسدودسازی DNS ایران) | REJECT |
| ۲ | `private.list` (نام‌های محلی مانند `.lan`، `.local`، `localhost`، صفحه روترها؛ بازه‌های IP خصوصی) | DIRECT |
| ۳ | `malware.list`، `phishing.list`، `cryptominers.list`، `ads.list` | REJECT |
| ۴ | `sanctioned.list` | PROXY |
| ۵ | `DOMAIN-SUFFIX,ir`، `ir-domains.list` | DIRECT |
| ۶ | `ir-cidr.list` | DIRECT |
| ۷ | `GEOIP,IR` | DIRECT |
| ۸ | `FINAL` | PROXY |

## افزودن لیست اختیاری

کانفیگ را ویرایش کنید (**Config** ←‏ انگشتتان را روی کانفیگ نگه دارید ←‏ **Edit Config** ←‏ **Rule**) و یک خط بالای
`FINAL` اضافه کنید. برای مثال:

```
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/nsfw.list,REJECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/telegram.list,PROXY
```

!!! warning "توجه"
    ویرایش محلی با آپدیت کانفیگ از `update-url` از بین می‌رود. برای حفظ تغییرات، کانفیگ خودتان را
    بسازید و فقط با خطوط `RULE-SET` به این لیست‌ها ارجاع دهید.

لینک همه لیست‌ها: [لینک فایل‌ها](urls.md#shadowrocket-rule-lists).

## لیست‌های سبک‌تر

iOS حافظه برنامه‌های VPN را محدود می‌کند و `ads.list` بزرگ است. دو گزینه سبک‌تر:

- **DOMAIN-SET**: همان دامنه‌ها در قالبی که Shadowrocket و Surge بهینه‌تر لود می‌کنند. خط `RULE-SET`
  یک لیست دامنه را با خط `DOMAIN-SET` جایگزین کنید، مثلاً
  `DOMAIN-SET,https://github.com/split-route/routing-rules/releases/latest/download/ads-domainset.txt,REJECT`. فایل‌های DOMAIN-SET فقط دامنه دارند، پس برای لیست‌هایی که
  بازه IP هم دارند (`malware`، `phishing`) همان `RULE-SET` را نگه دارید.
- **ads-lite**: حدود ۵ هزار قانون از منابع کوچک و دست‌چین‌شده (شبکه‌های تبلیغاتی و ردیاب‌های ایرانی، لیست
  Peter Lowe، v2fly) به‌جای حدود ۱۷۶ هزار. به‌جای `ads.list` از `ads-lite.list` (یا `ads-lite-domainset.txt`)
  استفاده کنید.

!!! note "نکته"
    کانفیگ پیش‌فرض `RULE-SET` با `ads.list` کامل را نگه می‌دارد. این گزینه‌ها روی همه دستگاه‌ها تست
    نشده‌اند؛ اگر برنامه مشکل حافظه دارد یا مدام دوباره وصل می‌شود، از آن‌ها استفاده کنید.
