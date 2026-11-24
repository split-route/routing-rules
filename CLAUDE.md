# CLAUDE.md

Routing rules and ready-made configs for proxy/VPN clients, rebuilt daily from public lists (currently Iran).
Go build pipeline, MkDocs site in English and Persian, published as GitHub release assets.

## Architecture

Data flows one way: `sources.yaml` → fetch → parse → merge → safety checks → emit → `dist/`, `out/`, `docs/`.

| Path | Role |
|------|------|
| `sources.yaml` | Source registry: categories, actions, upstream URLs and formats. Format documented at the top. |
| `data/` | Hand-maintained lists (allowlist, Iranian ad networks, DNS block-page IPs, private ranges) and `canaries.txt` (known answers). |
| `cmd/build` | Orchestration and all I/O: config, fetching with last-known-good fallback, rule order, canaries, writing outputs and generated docs. |
| `cmd/verify` | Checks compiled rule-sets and the PAC file against canaries using the clients' own tools (sing-box, mihomo, Node.js). |
| `internal/fetch` | HTTP with ETag cache. |
| `internal/parse` | Upstream formats (hosts, domains, ABP, v2fly, CIDR, JSON, DB-IP) → normalized `parse.Result`. |
| `internal/merge` | Dedup, allowlist, verdicts across categories. |
| `internal/emit` | One file per client format (`singbox.go`, `clash.go`, `v2ray.go`, `karing.go`, …), text templates in `templates/`, doc-page locales in `locales/`. |
| `scripts/` | CI steps: install tools, compile, check configs, publish, release notes. |
| `docs/en`, `docs/fa` | MkDocs pages. `urls.md` and `sources.md` are generated and gitignored. |

Rules:

- `internal/*` packages are pure: bytes in, bytes out, no network or filesystem side effects (except `fetch`). I/O lives in `cmd/`.
- Dependencies point downward: `cmd` → `emit`/`merge` → `parse`. `internal` packages never import `cmd`.
- Rule order is the contract: early REJECT (DNS block pages) → private DIRECT → REJECT → PROXY → Iran DIRECT → final PROXY. Every client format must keep the same order; `ruleOrder` in `cmd/build/rules.go` is the single source.
- Outputs must be byte-deterministic (sorted, no timestamps). A release is published only when checksums change, so non-determinism means a release every day.
- Safety checks run before anything is written: no empty lists, allowlist, canaries. A failing source falls back to its last good copy; never weaken these checks to make a build pass.
- Outputs are never committed to `main`. The release set plus every upstream geosite/geoip `.srs` (`rule-set-geosite/`, `rule-set-geoip/`, from SagerNet) and the full v2fly `.dat` files (`v2ray/`) is force-pushed as one orphan commit to the `release` branch (`scripts/publish-branch.sh`, only when its `SHA256SUMS` changes; a failed upstream fetch leaves the branch alone); that bot-only branch is the one force-push exception. Releases keep their current file set and stay under GitHub's 1000-asset cap: large generated sets go to `release` only. The app's default base URL is the branch's raw address, so root files there must keep the release's flat names.

### Adding a client format

1. Emitter in `internal/emit/<client>.go` (or a template in `templates/`), driven by `RouteData`/`ConfData`, not by re-deriving rule order.
2. A test in `internal/emit` that checks rule order and at least one REJECT/DIRECT line.
3. Wire it in `cmd/build`, add its row to the File URLs data, and add any new strings to `internal/emit/locales/fa.json`.
4. If the client has a CLI validator, add it to `scripts/check-configs.sh`.
5. Docs page in both `docs/en` and `docs/fa`, nav entry in `mkdocs.yml`, row in the README outputs table.

## Code

- Go, stdlib first. A new dependency needs a reason the stdlib can't cover.
- Small functions with one job; name things after what they do in this domain (`ruleOrder`, `checkCanaries`), not generic names (`process`, `handle`).
- No speculative abstractions: no interfaces with one implementation, no config knobs nobody asked for.
- Errors: wrap with context (`fmt.Errorf("parse %s: %w", name, err)`); never swallow. A build that would publish bad data must fail.
- Comments explain why, not what. Doc comments on exported identifiers.
- Tests: table-driven, inputs in `testdata/`, no network. Every bug fix gets a test that fails without it.
- English only in code, comments, identifiers, commits and logs. Persian text lives only in `docs/fa/`, `*.fa.md.tmpl`, `internal/emit/locales/fa.json` and `docs/assets/check-i18n.js`.

## Verify before reporting done

```sh
gofmt -l . && go vet ./... && go test ./...
go run ./cmd/build                 # full build; needs network
mkdocs build --strict              # docs; deps pinned in requirements-docs.txt (use a venv)
```

`scripts/check-configs.sh` and `cmd/verify` need sing-box, Xray and mihomo (`scripts/install-tools.sh`); CI runs them.

## Docs

- `splitroute` lives in its own repo ([split-route/splitroute](https://github.com/split-route/splitroute)); only its docs page (`splitroute.md`) is here.

- Every page change goes into both `docs/en` and `docs/fa`. English is the source; Persian is rewritten naturally, not translated word by word.
- `mkdocs build --strict` must pass; anchors referenced from other pages must stay stable (use `{#id}` on Persian headings that are linked).
- Screenshots: `docs/assets/img/<client>-<step>.webp`, `cwebp -q 90 -m 6`, width = 2× the display width set on the page (`{ width="720" }` → 1440 px). Shared by both languages via `../assets/img/`.
- Deep links (`shadowrocket://…`) are fine as buttons; always keep the manual steps as a fallback.

## Persian style

Write the way Iranian users actually talk about these apps: short, conversational sentences, everyday tech words.

| Use | Not |
|-----|-----|
| کانفیگ | پیکربندی |
| تب | زبانه |
| لینک (URL) | نشانی، نشانی وب |
| لیست | فهرست |
| پیست کنید | بچسبانید |
| انگشتتان را چند لحظه روی … نگه دارید | لمس کنید و نگه دارید |
| فعال است (toggle) | روشن است |
| آپدیت | به‌روزرسانی، به‌روز |
| ایمپورت | وارد کردن (for import; «وارد کنید» stays for typing into a field) |
| دانلود / لود | بارگیری |
| تست نشده | آزموده نشده |
| نکته (admonition title) | یادداشت |
| است / شود / کنید | می‌باشد / گردد / نمایید |

Undecided, keep as they are for now: قوانین, مسیریابی, مسدود/مسدودسازی, مستقیم, دامنه.

- App UI labels stay in English and bold, exactly as the app shows them: **Use Config**, **Diversion Rules**.
- Technical terms stay in Latin script: rule-set, sing-box, DNS, IP, `.srs`, resolve.
- Persian plural on a Latin word: `rule-setها`; use the zero-width non-joiner (U+200C) where Persian spelling needs it (`می‌شود`, `لیست‌ها`).
- Persian digits in prose (`۱۲`, `۱٫۷`); Latin digits inside code, versions in code spans and URLs.
- Every `←` in a step path is followed by U+200F (RIGHT-TO-LEFT MARK): `**Config** ←‏ **Modules**`. Without it, adjacent English labels merge into one left-to-right run and the steps display in reverse order.
