---
icon: material/vector-polyline
---

# Karing

[Karing](https://karing.app) is built on sing-box and has its own format for custom diversion (routing) rules.
`karing.json` is a ready-made rules file in that format, using this project's sing-box rule-sets.

## Get the rules file

[:material-download: Download karing.json](https://github.com/split-route/routing-rules/releases/latest/download/karing.json){ .md-button .md-button--primary }

## Import it

1. Download the file with a button above. Don't copy the text from a browser tab: a partial selection
   breaks the file (Karing reports `FormatException: Unterminated string`).
2. In Karing, tap **Diversion** → **Diversion Rules** → the pencil icon (top right). Make sure
   **Custom Diversion Group** is switched on.

    ![Karing: Diversion, then Diversion Rules, then the pencil icon](../assets/img/karing-open.webp){ width="720" loading=lazy }

3. On the **Custom Diversion Group** screen, tap **⋮** → **Import** and pick the downloaded file. The
   **🚫 Block** and **✅ Direct** groups appear in the list.

    ![Karing: tap the menu, choose Import, and the Block and Direct groups appear](../assets/img/karing-import.webp){ width="720" loading=lazy }

4. Make sure both groups are switched on.

Karing downloads the rule-sets itself and refreshes them, so you don't need to re-import when the lists change.

## What it contains

| Group | Action | Rule-sets |
|-------|--------|-----------|
| 🚫 Block | `block` | `geoip-ir-injected` (DNS block-page addresses), `geosite-malware`, `geoip-malware`, `geosite-phishing`, `geoip-phishing`, `geosite-cryptominers`, `geosite-ads` |
| ✅ Direct | `direct` | `geosite-private`, `geoip-private` (local names and private ranges), `geosite-ir`, `geoip-ir` |

Everything else goes through the proxy you select in Karing. `geosite-ir` already includes the `.ir` and
`.ایران` TLDs, so you don't need a separate `.ir` rule.

!!! tip "Adding your own rules"
    You can add your own domains (for example `domain_suffix` or `domain_regex`) to either group in Karing.
    Importing a new copy of this file may replace those edits, so keep your own rules in a separate group.

To use more rule-sets (`geosite-nsfw`, `geoip-telegram`, …), add their `.srs` URLs to a group; see
[File URLs](urls.md#sing-box-rule-sets-binary).
