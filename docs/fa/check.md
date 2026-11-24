---
icon: material/magnify
---

# بررسی دامنه

ببینید کانفیگ‌های پیش‌فرض با یک دامنه، لینک یا آدرس IP چه می‌کنند: مسدود، مستقیم یا از طریق پراکسی.
همه چیز در مرورگر شما و با لیست‌های منتشرشده انجام می‌شود.

<form id="vr-check-form" class="vr-check" autocomplete="off" data-repo="https://github.com/split-route/routing-rules" data-lists="../../dl/lists/">
  <input id="vr-check-input" class="vr-check__input" type="text" dir="ltr" placeholder="digikala.com, https://example.com/page or 2.144.0.1" aria-label="دامنه، لینک یا آدرس IP">
  <button class="md-button md-button--primary" type="submit">بررسی</button>
</form>

<div id="vr-check-result" class="vr-check__result" aria-live="polite"></div>

!!! note "چگونه تصمیم می‌گیرد"
    قوانین با همان ترتیب کانفیگ‌های پیش‌فرض بررسی می‌شوند: آدرس‌های صفحه مسدودسازی DNS ایران مسدود می‌شوند؛ نام‌های محلی و بازه‌های خصوصی مستقیم می‌روند؛ بدافزار،
    فیشینگ، استخراج رمزارز و تبلیغات مسدود می‌شوند؛ سایت‌های تحریمی از پراکسی عبور می‌کنند؛ دامنه‌های `.ir`،
    سایت‌های ایرانی و بازه‌های IP ایران مستقیم می‌روند؛ بقیه از پراکسی عبور می‌کنند. اولین بار لیست‌ها (چند
    مگابایت) دانلود می‌شوند، پس اولین بررسی کمی طول می‌کشد.

می‌توانید به نتیجه لینک بدهید: `https://split-route.github.io/routing-rules/fa/check/?q=example.com`.
اگر نتیجه‌ای اشتباه به نظر می‌رسد، [گزارش دهید](https://github.com/split-route/routing-rules/issues/new/choose).
