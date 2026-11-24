---
icon: material/web
---

# Browsers (PAC)

A proxy auto-config (PAC) file tells a browser which sites to send through a proxy. `proxy.pac` sends
Iranian domains (including all `.ir` domains) direct and everything else through your proxy client's
local proxy (v2rayN, Clash, sing-box and similar).

[:material-download: Download proxy.pac](https://github.com/split-route/routing-rules/releases/latest/download/proxy.pac){ .md-button .md-button--primary }

## PAC URL

macOS, Windows and Firefox take a PAC **URL**. Copy this one:

```
https://github.com/split-route/routing-rules/releases/latest/download/proxy.pac
```

It always serves the newest file, so the lists stay current without re-downloading. It uses the default
proxy address below; if your client listens on a different port, download the file, edit it, and point the
setting at your copy instead.

!!! warning "If the URL can't be reached"
    GitHub redirects the download to `release-assets.githubusercontent.com`, which is filtered on some
    Iranian networks. If the system can't fetch the PAC file, it connects directly, so nothing goes through
    the proxy. If sites suddenly open without the proxy, use a downloaded copy instead.

## Set the proxy address

Open the file and set the first line to your client's local proxy:

```js
var PROXY = "SOCKS5 127.0.0.1:10808; SOCKS 127.0.0.1:10808";
```

The defaults match v2rayN (SOCKS on `10808`; HTTP is `10809`, written as `PROXY 127.0.0.1:10809`).

!!! note "No direct fallback"
    The file deliberately has no `DIRECT` fallback. If your proxy client stops, sites fail to load instead of
    silently opening without protection.

## Use it

- **SwitchyOmega / ZeroOmega**: new profile → **PAC Profile** → paste the file's content.
- **macOS**: System Settings → **Network** → your connection (Wi-Fi or Ethernet) → **Details** → **Proxies** →
  turn on **Automatic proxy configuration** → paste the [PAC URL](#pac-url) → **OK**.
- **Windows**: Settings → **Network & internet** → **Proxy** → **Use setup script** → paste the
  [PAC URL](#pac-url) → **Save**.
- **Firefox**: Settings → Network Settings → **Automatic proxy configuration URL** → paste the
  [PAC URL](#pac-url), or a `file://` path to an edited copy.

A PAC file can't block ads, and it only matches domain names: Iranian sites opened by IP address go through
the proxy.
