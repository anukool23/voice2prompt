#!/usr/bin/env bash
# Install Ollama and pull the cleanup model used by the Phase 2 LLM layer.
# Everything runs locally afterward — no network needed at runtime.
set -euo pipefail

MODEL="${1:-qwen2.5:3b}"

if ! command -v ollama >/dev/null 2>&1; then
  echo "==> Installing Ollama via Homebrew"
  brew install ollama
else
  echo "==> Ollama already installed: $(ollama --version 2>/dev/null || echo present)"
fi

# Make sure the server is up (the app can also start it on demand).
if ! curl -fsS http://127.0.0.1:11434/api/tags >/dev/null 2>&1; then
  echo "==> Starting 'ollama serve' in the background"
  ollama serve >/tmp/ollama-serve.log 2>&1 &
  for _ in $(seq 1 30); do
    curl -fsS http://127.0.0.1:11434/api/tags >/dev/null 2>&1 && break
    sleep 1
  done
fi

echo "==> Pulling model: $MODEL  (a few GB, one time)"
ollama pull "$MODEL"

echo "==> Done. Installed models:"
ollama list
