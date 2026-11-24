---
icon: material/call-split
---

# macOS و Linux (splitroute)

`splitroute` یک ابزار خط فرمان کوچک است که بازه‌های IP ایران را **بیرون از VPN** نگه می‌دارد. route سیستمی
اضافه می‌کند تا این بازه‌ها از اتصال عادی شما (روتر Wi-Fi) بروند و بقیه ترافیک داخل تونل بماند. کنار برنامه‌های
VPN مثل NordVPN، OpenVPN یا WireGuard کار می‌کند، چون route یک بازه مشخص همیشه به route کلی VPN اولویت دارد.

می‌تواند لیست‌های IP را **مسدود** هم کند (مثلاً `malware`) و وقتی شبکه را عوض می‌کنید، سیستم از خواب بیدار
می‌شود یا VPN دوباره وصل می‌شود، routeها را درست نگه می‌دارد.

!!! warning "قبل از نصب: مطمئن شوید VPN شما اجازه می‌دهد"
    بعضی برنامه‌های VPN همه ترافیک را به‌زور داخل تونل می‌فرستند (معمولاً به‌عنوان kill switch) و آن‌وقت هیچ
    ابزار مبتنی بر route کار نمی‌کند. با VPN وصل، این را اجرا کنید:

    ```sh
    gw=$(ipconfig getoption en0 router)        # Linux: gw=$(ip -4 route show default | awk '{print $3; exit}')
    sudo route -n add -net 2.144.0.0/14 "$gw"  # Linux: sudo ip route add 2.144.0.0/14 via "$gw"
    route -n get 2.144.0.1                     # Linux: ip route get 2.144.0.1
    sudo route -n delete -net 2.144.0.0/14     # Linux: sudo ip route del 2.144.0.0/14
    ```

    اگر دستور سوم روتر شما و `en0` را نشان داد (نه `utun…`)، splitroute کار می‌کند.

## نصب

=== "Homebrew (macOS)"

    ```sh
    brew install --cask split-route/tap/splitroute
    ```

=== "فایل آماده"

    فایل مناسب سیستم خود (macOS یا Linux، amd64 یا arm64) را از
    [صفحه releaseها](https://github.com/split-route/splitroute/releases) دانلود کنید و `splitroute` را در
    `/usr/local/bin` کپی کنید.

=== "Go"

    ```sh
    go install github.com/split-route/splitroute@latest
    sudo cp "$(go env GOPATH)/bin/splitroute" /usr/local/bin/
    ```

کد و گزارش مشکل: [github.com/split-route/splitroute](https://github.com/split-route/splitroute).

## استفاده

```sh
sudo splitroute up        # download the lists and add the routes
splitroute status         # gateway, VPN, applied routes
sudo splitroute down      # remove every route it added
```

برای این‌که با عوض شدن شبکه و ری‌استارت هم کار کند، آن را به‌صورت سرویس سیستم نصب کنید (launchd در macOS و
systemd در Linux):

```sh
sudo splitroute install   # starts now and at boot
sudo splitroute uninstall # stops it and removes the routes
```

با افزودن `--dry-run` به هر دستور، تغییرات route فقط نمایش داده می‌شوند و اعمال نمی‌شوند.

## انتخاب لیست‌ها

بدون فایل کانفیگ، splitroute لیست `ir-cidr` (بازه‌های IP ایران) را مستقیم می‌فرستد. برای انتخاب لیست‌ها،
فایل `/etc/splitroute.conf` را بسازید:

```ini
ir-cidr  = direct   # through your normal connection, outside the VPN
malware  = block    # dropped
telegram = vpn      # left alone: follows the VPN
```

`splitroute lists` همه لیست‌های آخرین نسخه را نشان می‌دهد. فقط **بازه‌های IPv4** هر لیست route می‌شوند:
لیست‌هایی که فقط دامنه دارند، مثل `ads` یا `nsfw`، فعلاً اثری ندارند. بعد از ویرایش فایل، دوباره
`sudo splitroute up` را اجرا کنید (سرویس تغییرات را در آپدیت روزانه می‌گیرد، یا آن را ری‌استارت کنید).

## امنیت

- لیست‌ها با HTTPS دانلود و با `SHA256SUMS` همان نسخه چک می‌شوند. اگر دانلود نشد، آخرین نسخه چک‌شده استفاده
  می‌شود.
- لیستی که بازه‌ای بزرگ‌تر از `/8` داشته باشد رد می‌شود، تا یک لیست خراب یا دستکاری‌شده نتواند بیشتر ترافیک
  شما را بیرون از VPN بفرستد.
- بازه‌های خصوصی و خاص (شبکه محلی، `10.0.0.0/8`، `100.64.0.0/10` و …) هیچ‌وقت تغییر نمی‌کنند.
- `down` دقیقاً همان routeهایی را که splitroute اضافه کرده پاک می‌کند.

## محدودیت‌ها

- **kill switch در VPN** که همه ترافیک را به‌زور داخل تونل می‌فرستد این routeها را بی‌اثر می‌کند؛ تست بالا را
  ببینید.
- **فقط IPv4.** ترافیک IPv6 داخل تونل می‌ماند.
- **DNS تغییر نمی‌کند.** درخواست‌های DNS همان‌جایی می‌روند که VPN شما می‌فرستد.
- **تطبیق بر اساس IP است.** سایت‌های ایرانی که خارج از ایران میزبانی می‌شوند (مثلاً پشت یک CDN خارجی) همچنان
  از VPN می‌روند.
