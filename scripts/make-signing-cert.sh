#!/usr/bin/env bash
# Create a local self-signed code-signing certificate so that macOS TCC permissions
# (Microphone, Accessibility) PERSIST across rebuilds. Ad-hoc signing changes the
# app's identity on every build, forcing you to re-grant permissions each time.
#
# Run once. Then package-mac.sh auto-detects and uses it.
set -euo pipefail

NAME="Voice2Prompt Local"
KEYCHAIN="$HOME/Library/Keychains/login.keychain-db"

if security find-identity -v -p codesigning 2>/dev/null | grep -q "$NAME"; then
  echo "==> Signing identity '$NAME' already exists. Nothing to do."
  exit 0
fi

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "==> Generating self-signed code-signing certificate: $NAME"
openssl req -x509 -newkey rsa:2048 -nodes -days 3650 \
  -keyout "$TMP/key.pem" -out "$TMP/cert.pem" \
  -subj "/CN=$NAME" \
  -addext "keyUsage=critical,digitalSignature" \
  -addext "extendedKeyUsage=critical,codeSigning" >/dev/null 2>&1

openssl pkcs12 -export -out "$TMP/id.p12" \
  -inkey "$TMP/key.pem" -in "$TMP/cert.pem" -passout pass: >/dev/null 2>&1

echo "==> Importing into login keychain (allow codesign to use it)"
security import "$TMP/id.p12" -k "$KEYCHAIN" -P "" -T /usr/bin/codesign >/dev/null

# Let codesign use the key without an interactive prompt each time.
security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" "$KEYCHAIN" >/dev/null 2>&1 || true

echo "==> Done. Verify:"
security find-identity -v -p codesigning | grep "$NAME" || {
  echo "!! Certificate not found for code signing — you may need to trust it in Keychain Access."
  exit 1
}
echo
echo "Now run: bash scripts/package-mac.sh   (it will sign with '$NAME')"
