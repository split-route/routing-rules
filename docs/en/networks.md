---
icon: material/router-network
---

# WireGuard, OpenVPN and routers

For VPNs that aren't proxy apps: keep Iranian traffic outside the tunnel by IP range.

## WireGuard (split tunnel)

WireGuard sends a destination through the tunnel if it matches `AllowedIPs`. These files list **every
address except Iran's ranges and private networks**, so Iranian sites stay direct and everything else uses
the VPN. Paste the file's content after `AllowedIPs =` in the `[Peer]` section.

| File | Entries | Use it for |
|------|---------|------------|
| `wireguard-allowedips.txt` | ~16k, IPv4 and IPv6 | Desktop clients |
| `wireguard-allowedips-lite.txt` | ~1.7k IPv4, plus `::/0` | Phones. Excludes only Iranian blocks of `/20` or larger (about 90% of Iranian IPv4 addresses) and sends all IPv6 through the tunnel. |

[:material-download: Download wireguard-allowedips-lite.txt](https://github.com/split-route/routing-rules/releases/latest/download/wireguard-allowedips-lite.txt){ .md-button .md-button--primary } [:material-download: Download wireguard-allowedips.txt](https://github.com/split-route/routing-rules/releases/latest/download/wireguard-allowedips.txt){ .md-button }

!!! warning "DNS"
    Split tunneling by IP doesn't decide where DNS goes. Set the WireGuard `DNS =` to a resolver inside the
    tunnel, or Iranian DNS may return tampered answers for blocked sites.

## OpenVPN (split tunnel)

OpenVPN works the other way round: keep the server's full tunnel and add routes that send Iranian ranges
**outside** it, through your normal connection. Paste the file's lines into your client `.ovpn` file (anywhere
after `client`):

```
route 2.144.0.0 255.252.0.0 net_gateway
```

These routes are more specific than the server's default route, so they take priority. The server needs no
changes.

| File | Routes | Use it for |
|------|--------|------------|
| `openvpn-ir-routes.txt` | ~2k, IPv4 | Desktop clients |
| `openvpn-ir-routes-lite.txt` | ~350, IPv4 | Phones and routers. Only Iranian blocks of `/20` or larger (about 90% of Iranian IPv4 addresses). |

[:material-download: Download openvpn-ir-routes-lite.txt](https://github.com/split-route/routing-rules/releases/latest/download/openvpn-ir-routes-lite.txt){ .md-button .md-button--primary } [:material-download: Download openvpn-ir-routes.txt](https://github.com/split-route/routing-rules/releases/latest/download/openvpn-ir-routes.txt){ .md-button }

!!! note
    The routes are IPv4 only, so IPv6 stays in the tunnel. `net_gateway` is standard OpenVPN 2.x syntax
    (desktop, OpenVPN for Android); check that your app supports it if routes seem ignored. As with WireGuard,
    keep DNS inside the tunnel.

## MikroTik (RouterOS)

`mikrotik-ir.rsc` creates the address list `IR` (IPv4 and IPv6). Running it again replaces the list.

[:material-download: Download mikrotik-ir.rsc](https://github.com/split-route/routing-rules/releases/latest/download/mikrotik-ir.rsc){ .md-button .md-button--primary }

```
/tool fetch url="https://github.com/split-route/routing-rules/releases/latest/download/mikrotik-ir.rsc" dst-path=mikrotik-ir.rsc
/import file-name=mikrotik-ir.rsc
```

Then use `dst-address-list=IR` in your mangle or routing rules, for example to mark everything **not** in
`IR` for your VPN routing table. Add a scheduler entry to refresh it daily.

## OpenWrt, ipset and nftables

`ir-cidr4.txt` and `ir-cidr6.txt` hold one range per line:

```sh
ipset create ir hash:net
curl -sL https://github.com/split-route/routing-rules/releases/latest/download/ir-cidr4.txt | while read -r n; do ipset add ir "$n"; done
```

```sh
nft add set inet fw4 ir '{ type ipv4_addr; flags interval; }'
nft add element inet fw4 ir "{ $(curl -sL https://github.com/split-route/routing-rules/releases/latest/download/ir-cidr4.txt | paste -sd, -) }"
```

Apps on OpenWrt that use v2ray/Xray (PassWall, ShadowSocksR Plus+) can use `geosite.dat` / `geoip.dat`
instead; see the [v2ray / Xray guide](v2ray.md).
