---
icon: material/shield-search
---

# AdGuard Home and Pi-hole

The block lists are also published as DNS blocklists, for your own DNS filter at home or on a server.

| File | Contents | Works in |
|------|----------|----------|
| `adguard-default.txt` | Everything the default configs block: malware, phishing, cryptominers, ads | AdGuard Home, Pi-hole v6 |
| `adguard-<list>.txt` | One list, e.g. `adguard-ads.txt`, `adguard-tracking.txt`, `adguard-nsfw.txt` | AdGuard Home, Pi-hole v6 |
| `hosts-default.txt` | Same as `adguard-default.txt` in hosts format | Pi-hole v5, hosts files |

## AdGuard Home

**Filters → DNS blocklists → Add blocklist → Add a custom list**, then paste:

```
https://github.com/split-route/routing-rules/releases/latest/download/adguard-default.txt
```

## Pi-hole

**Adlists** (v6: **Lists**) → add the URL, then run **Update Gravity**:

=== "Pi-hole v6"

    ```
    https://github.com/split-route/routing-rules/releases/latest/download/adguard-default.txt
    ```

=== "Pi-hole v5"

    ```
    https://github.com/split-route/routing-rules/releases/latest/download/hosts-default.txt
    ```

!!! note
    Entries use `||domain^`, which blocks the domain and its subdomains. The hosts format can't express
    subdomains, so `hosts-default.txt` only blocks the listed domains themselves.

Every per-list file: [File URLs](urls.md#adguard-home-and-pi-hole).
