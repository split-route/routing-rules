---
icon: material/vector-polyline
---

# Karing

[Karing](https://karing.app) بر پایه sing-box ساخته شده و برای قوانین مسیریابی سفارشی (diversion) قالب
خودش را دارد. `karing.json` یک فایل قوانین آماده در همین قالب است که از rule-setهای sing-box این پروژه
استفاده می‌کند.

## دریافت فایل قوانین

[:material-download: دانلود karing.json](https://github.com/split-route/routing-rules/releases/latest/download/karing.json){ .md-button .md-button--primary }

## ایمپورت

1. فایل را با دکمه بالا دانلود کنید. متن را از تب مرورگر کپی نکنید: انتخاب ناقص فایل را خراب می‌کند
   (Karing خطای `FormatException: Unterminated string` می‌دهد).
2. در Karing روی **Diversion** ←‏ **Diversion Rules** ←‏ نماد مداد (بالا سمت راست) بزنید. مطمئن شوید
   **Custom Diversion Group** فعال است.

    ![Karing: گزینه Diversion، سپس Diversion Rules، سپس نماد مداد](../assets/img/karing-open.webp){ width="720" loading=lazy }

3. در صفحه **Custom Diversion Group** روی **⋮** ←‏ **Import** بزنید و فایل دانلودشده را انتخاب کنید. گروه‌های
   **🚫 Block** و **✅ Direct** به لیست اضافه می‌شوند.

    ![Karing: منو را باز کنید، Import را بزنید؛ گروه‌های Block و Direct اضافه می‌شوند](../assets/img/karing-import.webp){ width="720" loading=lazy }

4. مطمئن شوید هر دو گروه فعال هستند.

Karing خودش rule-setها را دانلود و آپدیت می‌کند، پس با تغییر لیست‌ها لازم نیست دوباره ایمپورت کنید.

## محتوا

| گروه | عمل | rule-setها |
|------|-----|-----------|
| 🚫 Block | `block` | `geoip-ir-injected` (آدرس‌های صفحه مسدودسازی DNS)، `geosite-malware`، `geoip-malware`، `geosite-phishing`، `geoip-phishing`، `geosite-cryptominers`، `geosite-ads` |
| ✅ Direct | `direct` | `geosite-private`، `geoip-private` (نام‌های محلی و بازه‌های خصوصی)، `geosite-ir`، `geoip-ir` |

بقیه ترافیک از پراکسی‌ای که در Karing انتخاب کرده‌اید عبور می‌کند. `geosite-ir` دامنه‌های `.ir` و
`.ایران` را هم شامل می‌شود، پس قانون جداگانه‌ای برای `.ir` لازم نیست.

!!! tip "افزودن قوانین خودتان"
    می‌توانید دامنه‌های خودتان (مثلاً `domain_suffix` یا `domain_regex`) را در Karing به هر کدام از گروه‌ها
    اضافه کنید. ایمپورت نسخه جدید این فایل ممکن است آن تغییرات را پاک کند، پس قوانین خودتان را در گروهی
    جداگانه نگه دارید.

برای rule-setهای بیشتر (`geosite-nsfw`، `geoip-telegram` و …) لینک `.srs` آن‌ها را به یک گروه اضافه
کنید؛ ببینید [لینک فایل‌ها](urls.md#sing-box-rule-sets-binary).
