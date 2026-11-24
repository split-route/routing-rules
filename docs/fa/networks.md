---
icon: material/router-network
---

# WireGuard، OpenVPN و روترها

برای VPNهایی که برنامه پراکسی نیستند: ترافیک ایران را بر اساس بازه IP بیرون از تونل نگه دارید.

## WireGuard (تونل جداشده)

WireGuard مقصدی را که با `AllowedIPs` تطبیق پیدا کند از تونل عبور می‌دهد. این فایل‌ها **همه آدرس‌ها به‌جز
بازه‌های ایران و شبکه‌های خصوصی** را دارند، پس سایت‌های ایرانی مستقیم می‌مانند و بقیه از VPN استفاده
می‌کنند. محتوای فایل را بعد از `AllowedIPs =` در بخش `[Peer]` پیست کنید.

| فایل | تعداد | مناسب برای |
|------|-------|------------|
| `wireguard-allowedips.txt` | حدود ۱۶ هزار، IPv4 و IPv6 | کلاینت‌های دسکتاپ |
| `wireguard-allowedips-lite.txt` | حدود ۱٫۷ هزار IPv4، به‌علاوه `::/0` | گوشی‌ها. فقط بلوک‌های ایرانی `/20` یا بزرگ‌تر را مستثنا می‌کند (حدود ۹۰٪ آدرس‌های IPv4 ایران) و همه IPv6 را از تونل عبور می‌دهد. |

[:material-download: دانلود wireguard-allowedips-lite.txt](https://github.com/split-route/routing-rules/releases/latest/download/wireguard-allowedips-lite.txt){ .md-button .md-button--primary } [:material-download: دانلود wireguard-allowedips.txt](https://github.com/split-route/routing-rules/releases/latest/download/wireguard-allowedips.txt){ .md-button }

!!! warning "DNS"
    تونل جداشده بر اساس IP تعیین نمی‌کند DNS کجا برود. `DNS =` در WireGuard را روی یک سرور DNS داخل تونل
    بگذارید، وگرنه DNS ایران ممکن است برای سایت‌های مسدود پاسخ دستکاری‌شده برگرداند.

## OpenVPN (تونل جداشده)

OpenVPN برعکس عمل می‌کند: تونل کامل سرور را نگه دارید و مسیرهایی اضافه کنید که بازه‌های ایران را **بیرون**
از تونل و از طریق اتصال عادی شما بفرستند. خطوط فایل را در فایل `.ovpn` کلاینت (هر جایی بعد از `client`)
پیست کنید:

```
route 2.144.0.0 255.252.0.0 net_gateway
```

این مسیرها دقیق‌تر از مسیر پیش‌فرض سرور هستند، پس اولویت دارند. سرور به هیچ تغییری نیاز ندارد.

| فایل | تعداد مسیر | مناسب برای |
|------|-----------|------------|
| `openvpn-ir-routes.txt` | حدود ۲ هزار، IPv4 | کلاینت‌های دسکتاپ |
| `openvpn-ir-routes-lite.txt` | حدود ۳۵۰، IPv4 | گوشی‌ها و روترها. فقط بلوک‌های ایرانی `/20` یا بزرگ‌تر (حدود ۹۰٪ آدرس‌های IPv4 ایران). |

[:material-download: دانلود openvpn-ir-routes-lite.txt](https://github.com/split-route/routing-rules/releases/latest/download/openvpn-ir-routes-lite.txt){ .md-button .md-button--primary } [:material-download: دانلود openvpn-ir-routes.txt](https://github.com/split-route/routing-rules/releases/latest/download/openvpn-ir-routes.txt){ .md-button }

!!! note "نکته"
    مسیرها فقط IPv4 هستند، پس IPv6 در تونل می‌ماند. `net_gateway` نحو استاندارد OpenVPN 2.x است (دسکتاپ،
    OpenVPN for Android)؛ اگر مسیرها نادیده گرفته می‌شوند، پشتیبانی برنامه خود را بررسی کنید. مانند WireGuard،
    DNS را داخل تونل نگه دارید.

## MikroTik (RouterOS)

`mikrotik-ir.rsc` لیست آدرس `IR` (IPv4 و IPv6) را می‌سازد. اجرای دوباره آن لیست را جایگزین می‌کند.

[:material-download: دانلود mikrotik-ir.rsc](https://github.com/split-route/routing-rules/releases/latest/download/mikrotik-ir.rsc){ .md-button .md-button--primary }

```
/tool fetch url="https://github.com/split-route/routing-rules/releases/latest/download/mikrotik-ir.rsc" dst-path=mikrotik-ir.rsc
/import file-name=mikrotik-ir.rsc
```

سپس در قوانین mangle یا مسیریابی از `dst-address-list=IR` استفاده کنید؛ مثلاً هر چیزی را که در `IR` **نیست**
برای جدول مسیریابی VPN علامت بزنید. یک scheduler اضافه کنید تا هر روز آپدیت شود.

## OpenWrt، ipset و nftables

`ir-cidr4.txt` و `ir-cidr6.txt` در هر خط یک بازه دارند:

```sh
ipset create ir hash:net
curl -sL https://github.com/split-route/routing-rules/releases/latest/download/ir-cidr4.txt | while read -r n; do ipset add ir "$n"; done
```

```sh
nft add set inet fw4 ir '{ type ipv4_addr; flags interval; }'
nft add element inet fw4 ir "{ $(curl -sL https://github.com/split-route/routing-rules/releases/latest/download/ir-cidr4.txt | paste -sd, -) }"
```

برنامه‌هایی روی OpenWrt که از v2ray/Xray استفاده می‌کنند (PassWall، ShadowSocksR Plus+) می‌توانند به‌جای این
فایل‌ها از `geosite.dat` / `geoip.dat` استفاده کنند؛ ببینید [راهنمای v2ray / Xray](v2ray.md).
