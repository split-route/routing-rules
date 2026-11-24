#!/usr/bin/env bash
# Publishes every generated file as one orphan commit on the `release` branch, only when the
# content changed. The branch is force-pushed on purpose: it always holds a single commit, so
# thousands of files never pile up as history. Layout:
#   /                      the release set ($REL: .list, .srs, .mrs, configs, geosite.dat ...) and SHA256SUMS,
#                          so <raw base>/release/<name>.list works like releases/latest/download/<name>.list
#   /clash/                the Clash text sources
#   /rule-set-geosite/     every upstream geosite category as .srs (SagerNet/sing-geosite)
#   /rule-set-geoip/       every upstream geoip country as .srs (SagerNet/sing-geoip)
#   /v2ray/                the full upstream geosite.dat (v2fly/domain-list-community) and geoip.dat (v2fly/geoip)
# Upstream sets live in their own folders so they never shadow our own geosite-*/geoip-* files.
# Run it from the repo root after the build, compile and assemble steps.
# Needs REL, GH_TOKEN, GITHUB_REPOSITORY, GITHUB_SHA, BOT_NAME and BOT_EMAIL.
# PUBLISH_DRY_RUN=1 builds the tree and prints what it would do, without pushing.
set -euo pipefail

branch=release
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
tree=$work/tree
mkdir "$tree"
cp -R "$REL"/. "$tree/"
rm -f "$tree/SHA256SUMS"   # rebuilt below for the whole tree
cp -R out/clash "$tree/clash"

# fetch_geo <repo> <dir> <min files>: the rule-set branch of a SagerNet geo repo, checked before use.
# A failed or implausible fetch keeps the branch as it is: yesterday's lists beat none.
fetch_geo() {
  local repo=$1 dir=$2 min=$3 src=$work/$2.src
  git clone -q --depth 1 --branch rule-set "https://github.com/SagerNet/$repo.git" "$src" || return 1
  mkdir "$tree/$dir"
  find "$src" -maxdepth 1 -name '*.srs' -exec cp {} "$tree/$dir/" \;
  local n; n=$(find "$tree/$dir" -name '*.srs' | wc -l | tr -d ' ')
  [ "$n" -ge "$min" ] || { echo "::warning::$repo has only $n rule-sets (want $min+)"; return 1; }
  for f in "$tree/$dir"/*.srs; do   # "SRS" magic: a rule-set, not an HTML error page
    [ "$(head -c 3 "$f")" = "SRS" ] || { echo "::warning::$f is not a rule-set"; return 1; }
  done
}
# fetch_dat <repo> <asset> <name> <min bytes>: a v2fly release file, checked against its published sha256.
fetch_dat() {
  local url=https://github.com/v2fly/$1/releases/latest/download/$2 f=$tree/v2ray/$3
  mkdir -p "$tree/v2ray"
  curl -fsSL --retry 3 -o "$f" "$url" && curl -fsSL --retry 3 -o "$work/$3.sum" "$url.sha256sum" || return 1
  [ "$(sha256sum "$f" | cut -d' ' -f1)" = "$(cut -d' ' -f1 "$work/$3.sum")" ] || { echo "::warning::$3 does not match its sha256"; return 1; }
  [ "$(wc -c < "$f")" -ge "$4" ] || { echo "::warning::$3 is too small"; return 1; }
}
if ! { fetch_geo sing-geosite rule-set-geosite 1000 && fetch_geo sing-geoip rule-set-geoip 200 \
  && fetch_dat domain-list-community dlc.dat geosite.dat 1000000 && fetch_dat geoip geoip.dat geoip.dat 5000000; }; then
  echo "::warning::upstream geo sets unavailable; $branch branch left as it is"
  exit 0
fi

# GitHub rejects files over 100 MB; fail before the push with a clear message.
big=$(find "$tree" -type f -size +95M)
if [ -n "$big" ]; then
  echo "::error::files too large for a git push (over 95 MB): $big"
  exit 1
fi

# Unchanged sources give an identical SHA256SUMS.
(cd "$tree" && find . -type f -print0 | LC_ALL=C sort -z | xargs -0 sha256sum | sed 's#  \./#  #' > "$work/SHA256SUMS")
mv "$work/SHA256SUMS" "$tree/SHA256SUMS"
echo "$branch branch: $(wc -l < "$tree/SHA256SUMS") files, $(du -sm "$tree" | cut -f1) MB"

if [ "${PUBLISH_DRY_RUN:-}" = 1 ]; then
  echo "dry run: not pushing"
  exit 0
fi

if gh api "repos/$GITHUB_REPOSITORY/contents/SHA256SUMS?ref=$branch" -H 'Accept: application/vnd.github.raw' > "$work/prev" 2>/dev/null \
  && cmp -s "$work/prev" "$tree/SHA256SUMS"; then
  echo "No changes since the last push to $branch."
  exit 0
fi

cd "$tree"
git init -q -b "$branch"
git config user.name "$BOT_NAME"
git config user.email "$BOT_EMAIL"
git add -A
git commit -q -m "chore: generated files from ${GITHUB_SHA::7}"
git push -q --force "https://x-access-token:${GH_TOKEN}@github.com/${GITHUB_REPOSITORY}.git" "$branch"
echo "Pushed $branch."
