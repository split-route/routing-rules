---
icon: material/shield-search
---

# AdGuard Home و Pi-hole

لیست‌های مسدودسازی به شکل لیست مسدودسازی DNS هم منتشر می‌شوند، برای فیلتر DNS خودتان در خانه یا روی
سرور.

| فایل | محتوا | قابل استفاده در |
|------|-------|-----------------|
| `adguard-default.txt` | هر چیزی که کانفیگ‌های پیش‌فرض مسدود می‌کنند: بدافزار، فیشینگ، استخراج رمزارز، تبلیغات | AdGuard Home، Pi-hole v6 |
| `adguard-<list>.txt` | یک لیست، مثلاً `adguard-ads.txt`، `adguard-tracking.txt`، `adguard-nsfw.txt` | AdGuard Home، Pi-hole v6 |
| `hosts-default.txt` | همان `adguard-default.txt` در قالب hosts | Pi-hole v5، فایل‌های hosts |

## AdGuard Home

**Filters ←‏ DNS blocklists ←‏ Add blocklist ←‏ Add a custom list** و سپس این را پیست کنید:

```
https://github.com/split-route/routing-rules/releases/latest/download/adguard-default.txt
```

## Pi-hole

**Adlists** (در v6: **Lists**) ←‏ لینک را اضافه کنید و سپس **Update Gravity** را اجرا کنید:

=== "Pi-hole v6"

    ```
    https://github.com/split-route/routing-rules/releases/latest/download/adguard-default.txt
    ```

=== "Pi-hole v5"

    ```
    https://github.com/split-route/routing-rules/releases/latest/download/hosts-default.txt
    ```

!!! note "نکته"
    ورودی‌ها از `||domain^` استفاده می‌کنند که دامنه و زیردامنه‌هایش را مسدود می‌کند. قالب hosts نمی‌تواند
    زیردامنه‌ها را بیان کند، پس `hosts-default.txt` فقط خود دامنه‌های داخل لیست را مسدود می‌کند.

همه فایل‌های هر لیست: [لینک فایل‌ها](urls.md#adguard-home-and-pi-hole).
