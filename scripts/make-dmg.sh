#!/usr/bin/env bash
# Build a shareable Voice2Prompt.dmg (drag-to-Applications installer).
#
# Tiers:
#   * No env vars           → ad-hoc signed DMG. Runs on this Mac; on OTHER Macs the
#                             recipient must bypass Gatekeeper (see the printed note).
#   * DEVELOPER_ID set      → signs the app with a Developer ID (hardened runtime).
#   * DEVELOPER_ID + notary → also notarizes & staples for clean install anywhere.
#
# Env vars for a distributable build:
#   DEVELOPER_ID="Developer ID Application: Your Name (TEAMID)"
#   NOTARY_APPLE_ID="you@apple.id"  NOTARY_TEAM_ID="TEAMID"  NOTARY_PASSWORD="app-specific-pw"
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP="$ROOT/desktop/build/bin/Voice2Prompt.app"
DIST="$ROOT/dist"
DMG="$DIST/Voice2Prompt.dmg"
VOL="Voice2Prompt"

export PATH="/opt/homebrew/bin:$HOME/go/bin:$PATH"

echo "==> Building & bundling the app"
bash "$ROOT/scripts/package-mac.sh" >/dev/null
[ -d "$APP" ] || { echo "app build failed"; exit 1; }

# Optional: re-sign with a real Developer ID + hardened runtime for distribution.
if [ -n "${DEVELOPER_ID:-}" ]; then
  echo "==> Signing with Developer ID: $DEVELOPER_ID"
  codesign --force --deep --options runtime --timestamp \
    --sign "$DEVELOPER_ID" "$APP"
fi

echo "==> Staging DMG contents"
mkdir -p "$DIST"
STAGE="$(mktemp -d)"
trap 'rm -rf "$STAGE"' EXIT
cp -R "$APP" "$STAGE/"
ln -s /Applications "$STAGE/Applications"

echo "==> Creating $DMG"
rm -f "$DMG"
hdiutil create -volname "$VOL" -srcfolder "$STAGE" -ov -format UDZO "$DMG" >/dev/null

# Optional: notarize + staple so Gatekeeper accepts it on any Mac.
if [ -n "${DEVELOPER_ID:-}" ] && [ -n "${NOTARY_APPLE_ID:-}" ]; then
  echo "==> Notarizing (this can take a few minutes)"
  xcrun notarytool submit "$DMG" \
    --apple-id "$NOTARY_APPLE_ID" --team-id "$NOTARY_TEAM_ID" \
    --password "$NOTARY_PASSWORD" --wait
  echo "==> Stapling"
  xcrun stapler staple "$DMG"
  echo "✅ Notarized DMG — installs cleanly on any Apple Silicon Mac."
else
  echo
  echo "⚠️  Ad-hoc DMG (not notarized). On another Mac, the recipient must bypass Gatekeeper:"
  echo "    1) Try to open it, then System Settings → Privacy & Security → 'Open Anyway', OR"
  echo "    2) Terminal:  xattr -dr com.apple.quarantine /Applications/Voice2Prompt.app"
  echo
  echo "    For a clean install-anywhere DMG, re-run with DEVELOPER_ID + NOTARY_* env vars."
fi

echo "==> Done: $DMG"
du -sh "$DMG" | awk '{print "    size:", $1}'
