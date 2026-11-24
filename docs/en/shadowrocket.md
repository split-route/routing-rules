---
icon: material/rocket-launch-outline
---

# Shadowrocket

## Install with one tap

On your iPhone or iPad with Shadowrocket installed, tap the button. Shadowrocket opens, downloads the config
and switches to it.

[:material-rocket-launch: Install in Shadowrocket](shadowrocket://config/add/https://github.com/split-route/routing-rules/releases/latest/download/shadowrocket.conf){ .md-button .md-button--primary }

Then, on the **Home** tab, add your server and set **Global Routing** to **Config**.

!!! note
    The button works only on the device where Shadowrocket is installed. It replaces the config in use; to keep
    your own config, use the [module](#keep-your-own-config-module) instead. If nothing happens, follow the
    manual steps below.

## Import manually

1. Copy the config URL:

    ```
    https://github.com/split-route/routing-rules/releases/latest/download/shadowrocket.conf
    ```

2. Shadowrocket → **Config** tab → **+** → paste the URL → **Download**.
3. Tap the new config → **Use Config**. A checkmark marks the config in use.

    ![Config tab: tap +, paste the config URL, tap Download, then tap the config and Use Config](../assets/img/shadowrocket-import.webp){ width="720" loading=lazy }

4. On the **Home** tab, add your server and set **Global Routing** to **Config**.

    ![Home tab: Global Routing set to Config](../assets/img/shadowrocket-routing.webp){ width="480" loading=lazy }

### Keep it updated

The lists change daily. The config includes `update-url` and every rule list is a remote `RULE-SET`, so you
don't need to re-import. To refresh now, long-press the config in the **Config** tab → **Update**.

![Long-press the config and tap Update](../assets/img/shadowrocket-update.webp){ width="300" loading=lazy }

## Keep your own config (module)

If you already have a config you like (your own servers, rules or settings), add these rules as a
**module** instead of importing the config above. The module holds the same rules without the final
catch-all, so your config keeps its servers and `FINAL`.

[:material-puzzle-outline: Install module in Shadowrocket](shadowrocket://install?module=https://github.com/split-route/routing-rules/releases/latest/download/routing-rules.sgmodule){ .md-button }

Or add it manually:

1. Copy the module URL:

    ```
    https://github.com/split-route/routing-rules/releases/latest/download/routing-rules.sgmodule
    ```

2. Shadowrocket → **Config** tab → **Modules** → **+** → paste the URL → **Download**.
3. Make sure the module is switched on.

The rules use the policy `PROXY` (Shadowrocket's built-in proxy policy) for sites that block Iranian IPs.
Modules update on their own like the config; to refresh now, update the module from the **Modules**
screen.

## Rule order

| # | Rule | Action |
|---|------|--------|
| 1 | `ir-injected.list` (Iran's DNS block-page addresses) | REJECT |
| 2 | `private.list` (local names such as `.lan`, `.local`, `localhost`, router pages; private IP ranges) | DIRECT |
| 3 | `malware.list`, `phishing.list`, `cryptominers.list`, `ads.list` | REJECT |
| 4 | `sanctioned.list` | PROXY |
| 5 | `DOMAIN-SUFFIX,ir`, `ir-domains.list` | DIRECT |
| 6 | `ir-cidr.list` | DIRECT |
| 7 | `GEOIP,IR` | DIRECT |
| 8 | `FINAL` | PROXY |

## Add an optional list

Edit the config (**Config** → long-press the config → **Edit Config** → **Rule**) and add a line above `FINAL`. For
example:

```
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/nsfw.list,REJECT
RULE-SET,https://github.com/split-route/routing-rules/releases/latest/download/telegram.list,PROXY
```

!!! warning
    A local edit is lost when the config updates from `update-url`. To keep your changes, make your own
    config and only reference these lists with `RULE-SET` lines.

All list URLs: [File URLs](urls.md#shadowrocket-rule-lists).

## Lighter lists

iOS limits how much memory a VPN app may use, and `ads.list` is large. Two lighter options:

- **DOMAIN-SET**: the same domains in a format Shadowrocket and Surge load more efficiently. Replace a
  domain list's `RULE-SET` line with a `DOMAIN-SET` line, e.g.
  `DOMAIN-SET,https://github.com/split-route/routing-rules/releases/latest/download/ads-domainset.txt,REJECT`. DOMAIN-SET files hold only domains, so keep the `RULE-SET`
  lines for lists that also contain IP ranges (`malware`, `phishing`).
- **ads-lite**: about 5k rules from small curated sources (Iranian ad networks and trackers, Peter Lowe's
  list, v2fly) instead of about 176k. Use `ads-lite.list` (or `ads-lite-domainset.txt`) in place of `ads.list`.

!!! note
    The default config keeps `RULE-SET` with the full `ads.list`. These options haven't been tested on every
    device; switch if the app reports memory problems or reconnects often.
