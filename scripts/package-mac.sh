#!/usr/bin/env bash
# Build Voice2Prompt.app and bundle the whisper-server binary, its dylibs, and the
# whisper model inside it so the app is self-contained and runs from anywhere.
# (Ad-hoc signed — for a distributable build you also need Developer ID + notarization.)
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MODEL="${1:-models/ggml-base.en.bin}"
WHISPER_BIN_DIR="$ROOT/third_party/whisper.cpp/build/bin"
APP="$ROOT/desktop/build/bin/Voice2Prompt.app"
RES="$APP/Contents/Resources"

export PATH="/opt/homebrew/bin:$HOME/go/bin:$PATH"
export GOSUMDB=off GOFLAGS=-mod=mod

echo "==> Building app"
(cd "$ROOT/desktop" && wails build -skipbindings >/dev/null 2>&1)
[ -d "$APP" ] || { echo "build failed: $APP missing"; exit 1; }

echo "==> Bundling whisper-server + dylibs into Resources"
cp "$WHISPER_BIN_DIR/whisper-server" "$RES/"
cp "$WHISPER_BIN_DIR"/*.dylib "$RES/"

echo "==> Rewriting rpath to @loader_path (so colocated dylibs resolve anywhere)"
fixrpath() {
  local f="$1"
  # Drop the absolute build-dir rpath if present (ignore if absent).
  install_name_tool -delete_rpath "$WHISPER_BIN_DIR" "$f" 2>/dev/null || true
  # Add @loader_path so each binary finds its siblings.
  install_name_tool -add_rpath "@loader_path" "$f" 2>/dev/null || true
}
fixrpath "$RES/whisper-server"
for dylib in "$RES"/*.dylib; do fixrpath "$dylib"; done

echo "==> Bundling model: $MODEL"
mkdir -p "$RES/models"
cp "$ROOT/$MODEL" "$RES/models/"

# Prefer a stable local identity so Microphone/Accessibility permissions persist
# across rebuilds. Falls back to ad-hoc (permissions reset each build).
SIGN_ID="${V2P_SIGN_ID:-}"
if [ -z "$SIGN_ID" ] && security find-identity -v -p codesigning 2>/dev/null | grep -q "Voice2Prompt Local"; then
  SIGN_ID="Voice2Prompt Local"
fi
if [ -n "$SIGN_ID" ]; then
  echo "==> Signing with stable identity: $SIGN_ID (permissions persist across rebuilds)"
  codesign --force --deep --sign "$SIGN_ID" "$APP"
else
  echo "==> Re-signing (ad-hoc). Tip: run scripts/make-signing-cert.sh so permissions"
  echo "    don't reset on each rebuild."
  codesign --force --deep --sign - "$APP"
fi

echo "==> Done: $APP"
du -sh "$APP" | awk '{print "    bundle size:", $1}'
