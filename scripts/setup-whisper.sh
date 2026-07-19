#!/usr/bin/env bash
# Clone and build whisper.cpp (with Metal on Apple Silicon) to get the
# `whisper-server` sidecar binary that the Go app talks to over HTTP.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VENDOR="$ROOT/third_party/whisper.cpp"

if [ ! -d "$VENDOR/.git" ]; then
  echo "==> Cloning whisper.cpp"
  git clone --depth 1 https://github.com/ggerganov/whisper.cpp "$VENDOR"
else
  echo "==> whisper.cpp already cloned"
fi

echo "==> Building whisper.cpp (Metal enabled by default on Apple Silicon)"
cmake -S "$VENDOR" -B "$VENDOR/build" \
  -DCMAKE_BUILD_TYPE=Release \
  -DGGML_METAL=ON \
  -DWHISPER_BUILD_EXAMPLES=ON \
  -DWHISPER_BUILD_TESTS=OFF
cmake --build "$VENDOR/build" -j --config Release --target whisper-server

echo "==> Done. Server binary:"
find "$VENDOR/build" -name "whisper-server" -type f
