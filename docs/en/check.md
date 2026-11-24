---
icon: material/magnify
---

# Check a domain

See what the default configs do with a domain, URL or IP address: blocked, direct, or through the proxy.
Everything runs in your browser using the published lists.

<form id="vr-check-form" class="vr-check" autocomplete="off" data-repo="https://github.com/split-route/routing-rules">
  <input id="vr-check-input" class="vr-check__input" type="text" placeholder="digikala.com, https://example.com/page or 2.144.0.1" aria-label="Domain, URL or IP address">
  <button class="md-button md-button--primary" type="submit">Check</button>
</form>

<div id="vr-check-result" class="vr-check__result" aria-live="polite"></div>

!!! note "How it decides"
    The rules are checked in the same order as the default configs: Iran's DNS block-page addresses
    are blocked; local names and private ranges go direct; malware, phishing, cryptominers and ads are blocked; sanctioned sites go through the proxy; `.ir`
    domains, Iranian sites and Iranian IP ranges go direct; everything else goes through the tunnel.
    The first visit downloads the lists (a few MB), so the first check takes a moment.

You can link to a result: `https://split-route.github.io/routing-rules/check/?q=example.com`.
If a result looks wrong, [report it](https://github.com/split-route/routing-rules/issues/new/choose).
