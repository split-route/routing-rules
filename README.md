# Routing Rules

Routing rules and ready-made configs for proxy and VPN clients, rebuilt daily by GitHub Actions from public
upstream lists. Currently built for **Iran**; the pipeline is country-agnostic.
Inspired by [Chocolate4U/Iran-v2ray-rules](https://github.com/Chocolate4U/Iran-v2ray-rules) and
[Chocolate4U/Iran-sing-box-rules](https://github.com/Chocolate4U/Iran-sing-box-rules).

- **Block** ads and trackers (including Iranian ad networks and trackers), malware, phishing, cryptominers,
  and Iran's DNS block-page addresses.
- **Proxy** sites that refuse Iranian IPs.
- **Direct** for Iranian domains and IP ranges; everything else goes through the tunnel.

**Outputs**, published as GitHub release assets:

| For | Files |
|-----|-------|
| sing-box | `.srs` and JSON rule-sets, route and DNS blocks |
| Clash / mihomo | `.mrs` rule-sets, rules and DNS blocks |
| v2ray / Xray | `geosite.dat`, `geoip.dat`, routing and DNS blocks |
| Shadowrocket, Surge, Loon | Full config, `RULE-SET` lists, `DOMAIN-SET` lists |
| Karing, v2rayN / v2rayNG | Ready-made rule files |
| WireGuard, OpenVPN, routers | Split-tunnel `AllowedIPs` and routes, MikroTik script, CIDR lists |
| macOS, Linux | [`splitroute`](https://github.com/split-route/splitroute): keeps Iranian IPs outside any VPN with system routes |
| Browsers | PAC file |
| AdGuard Home, Pi-hole | DNS blocklists |

## Documentation

Setup guides for each client, every file URL, and the full source list:
**https://split-route.github.io/routing-rules/**

## How it works

`cmd/build` reads [`sources.yaml`](sources.yaml), then fetches, parses, normalizes and merges every source.
It runs safety checks before writing anything: no empty lists, the allowlist, and known answers
(`data/canaries.txt`). A source that fails falls back to its last good copy. Outputs are never committed to `main`. CI
validates them with sing-box, Xray and mihomo, then publishes them as a new GitHub release, but only
when something changed since the last one. `releases/latest/download/<file>` always serves the newest copy.

The `release` branch holds the whole release set plus every upstream geosite category and geoip country as
`.srs` (`rule-set-geosite/`, `rule-set-geoip/`, from [SagerNet/sing-geosite](https://github.com/SagerNet/sing-geosite) and
[sing-geoip](https://github.com/SagerNet/sing-geoip)). It is one commit, force-pushed on each change (no history), with no
size or file-count limit like a release has. `v2ray/` holds the full upstream `geosite.dat` ([v2fly/domain-list-community](https://github.com/v2fly/domain-list-community)) and `geoip.dat` ([v2fly/geoip](https://github.com/v2fly/geoip)), checked against their published sha256 (the root `geosite.dat` and `geoip.dat` are ours): <https://raw.githubusercontent.com/split-route/routing-rules/release/SHA256SUMS>.
The macOS app downloads its lists from there.

```sh
go test ./...
go run ./cmd/build                                        # also regenerates docs/en/urls.md and docs/en/sources.md
pip install -r requirements-docs.txt && mkdocs serve      # preview the docs (English and Persian)
```

To add a source, add an entry to `sources.yaml`; the format is documented at the top of that file.
Hand-maintained lists live in [`data/`](data): [`allowlist.txt`](data/allowlist.txt) (domains that must never
be blocked), [`ir-ads.txt`](data/ir-ads.txt) (Iranian ad networks to always block),
[`ir-injected.txt`](data/ir-injected.txt) (Iran's DNS block-page addresses) and
[`private.txt`](data/private.txt) (private and special-use IP ranges).

## Verifying files

Every release includes `SHA256SUMS`, and each file carries a signed build provenance attestation from
this repository's GitHub Actions workflow:

```sh
gh attestation verify ads.list --repo split-route/routing-rules
```

Tool versions used for validation (sing-box, Xray, mihomo), GitHub Actions, Go modules and the docs
dependencies are kept up to date by [Renovate](renovate.json).

## License

The code is licensed under [GPL-3.0](LICENSE).

The generated lists are derived from third-party data and keep their upstream licenses. Some of those are
GPL-3.0 (hagezi, AdGuard) and AGPL-3.0 (PersianBlocker). Attribution:

- IP geolocation data: [IP Geolocation by DB-IP](https://db-ip.com) (CC BY 4.0).
- Country IP blocks: [ipverse](https://github.com/ipverse/country-ip-blocks) (CC0), [ipdeny.com](https://www.ipdeny.com), [RIPEstat](https://stat.ripe.net).
- Domain lists: [bootmortis/iran-hosted-domains](https://github.com/bootmortis/iran-hosted-domains), [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community), and the other sources listed on the [sources page](https://split-route.github.io/routing-rules/sources/).
