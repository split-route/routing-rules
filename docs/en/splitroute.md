---
icon: material/call-split
---

# macOS & Linux (splitroute)

`splitroute` is a small command-line tool that keeps Iranian IP ranges **outside your VPN**. It adds system
routes that send them through your normal connection (your Wi-Fi router), while everything else stays in
the tunnel. It works next to VPN apps like NordVPN, OpenVPN or WireGuard, because a route for a specific
range always wins over the VPN's catch-all route.

It can also **block** IP lists (for example `malware`), and it keeps the routes correct when you switch
networks, wake from sleep or the VPN reconnects.

!!! warning "Before you install: test that your VPN allows it"
    Some VPN apps force all traffic into the tunnel (often as part of a kill switch), and then no
    route-based tool works. With your VPN connected, run:

    ```sh
    gw=$(ipconfig getoption en0 router)        # Linux: gw=$(ip -4 route show default | awk '{print $3; exit}')
    sudo route -n add -net 2.144.0.0/14 "$gw"  # Linux: sudo ip route add 2.144.0.0/14 via "$gw"
    route -n get 2.144.0.1                     # Linux: ip route get 2.144.0.1
    sudo route -n delete -net 2.144.0.0/14     # Linux: sudo ip route del 2.144.0.0/14
    ```

    If the third command shows your router and `en0` (not `utun…`), splitroute will work.

## Install

=== "Homebrew (macOS)"

    ```sh
    brew install --cask split-route/tap/splitroute
    ```

=== "Binary"

    Download the archive for your system (macOS or Linux, amd64 or arm64) from the
    [releases page](https://github.com/split-route/splitroute/releases) and copy `splitroute` to `/usr/local/bin`.

=== "Go"

    ```sh
    go install github.com/split-route/splitroute@latest
    sudo cp "$(go env GOPATH)/bin/splitroute" /usr/local/bin/
    ```

Source code and issues: [github.com/split-route/splitroute](https://github.com/split-route/splitroute).

## Use

```sh
sudo splitroute up        # download the lists and add the routes
splitroute status         # gateway, VPN, applied routes
sudo splitroute down      # remove every route it added
```

To keep it running across network changes and reboots, install it as a system service (launchd on macOS,
systemd on Linux):

```sh
sudo splitroute install   # starts now and at boot
sudo splitroute uninstall # stops it and removes the routes
```

Add `--dry-run` to any command to print the route changes instead of making them.

## Choose lists

Without a config file, splitroute sends `ir-cidr` (Iranian IP ranges) direct. To choose lists yourself,
create `/etc/splitroute.conf`:

```ini
ir-cidr  = direct   # through your normal connection, outside the VPN
malware  = block    # dropped
telegram = vpn      # left alone: follows the VPN
```

`splitroute lists` shows every list in the latest release. Only a list's **IPv4 ranges** are routed:
domain-only lists such as `ads` or `nsfw` have no effect yet. After editing the file, run `sudo splitroute up`
again (the service picks changes up with its daily refresh, or restart it).

## Safety

- Lists are downloaded over HTTPS and checked against the release's `SHA256SUMS`. If a download fails, the
  last verified copy is used.
- A list with a range broader than `/8` is refused, so a broken or tampered list can't send most of your
  traffic outside the VPN.
- Private and special ranges (your LAN, `10.0.0.0/8`, `100.64.0.0/10`, …) are never touched.
- `down` removes exactly the routes splitroute added.

## Limits

- **VPN kill switches** that force all traffic into the tunnel override these routes; see the test above.
- **IPv4 only.** IPv6 traffic stays in the tunnel.
- **DNS is not changed.** Lookups go wherever your VPN sends them.
- **Matching is by IP.** Iranian sites hosted abroad (for example behind a foreign CDN) still use the VPN.
