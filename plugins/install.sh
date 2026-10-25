#!/usr/bin/env bash
# Install a shokku companion plugin on a dokku host, independently of shokku:
#   curl -fsSL https://raw.githubusercontent.com/randspace0/shokkupan/main/plugins/install.sh | sudo bash -s discord
# Override SHOKKU_REPO / SHOKKU_REF to install from a fork, branch or tag.
set -euo pipefail

NAME="${1:-}"
REPO="${SHOKKU_REPO:-randspace0/shokkupan}"
REF="${SHOKKU_REF:-main}"
PLUGINS_DIR="/var/lib/dokku/plugins/available"

if [[ "$(id -u)" != "0" ]]; then
  echo "This script must be run as root" 1>&2
  exit 1
fi
if ! command -v dokku &>/dev/null; then
  echo "Please install dokku first using the instructions at https://dokku.com" 1>&2
  exit 1
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "=> downloading $REPO@$REF"
curl -fsSL "https://github.com/$REPO/archive/$REF.tar.gz" | tar -xz -C "$tmp" --strip-components=1

if [[ ! "$NAME" =~ ^[a-z][a-z0-9-]*$ || ! -f "$tmp/plugins/$NAME/plugin.toml" ]]; then
  echo "usage: install.sh <plugin>   available: $(cd "$tmp/plugins" && ls -d */ | tr -d / | tr '\n' ' ')" 1>&2
  exit 1
fi

echo "=> installing plugin $NAME"
rm -rf "${PLUGINS_DIR:?}/$NAME"
cp -r "$tmp/plugins/$NAME" "$PLUGINS_DIR/$NAME"
dokku plugin:enable "$NAME" || true
dokku plugin:install
echo "=> done"
