#!/bin/bash
set -euo pipefail

# =============================================================================
# Ollama Setup Script
# =============================================================================

MODEL="qwen2.5-coder:7b"

echo "=== Ollama Setup ==="

# Check if Ollama is installed
if command -v ollama &> /dev/null; then
    echo "[OK] Ollama is already installed: $(ollama --version)"
else
    echo "[INFO] Installing Ollama..."

    # macOS installation
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # Check if Homebrew is available
        if command -v brew &> /dev/null; then
            brew install ollama
        else
            echo "[INFO] Downloading Ollama installer..."
            curl -fsSL https://ollama.ai/install.sh | sh
        fi
    else
        # Linux installation
        curl -fsSL https://ollama.ai/install.sh | sh
    fi

    echo "[OK] Ollama installed: $(ollama --version)"
fi

# Start Ollama service (if not running)
echo "[INFO] Checking Ollama service..."
if ! pgrep -x "ollama" > /dev/null; then
    echo "[INFO] Starting Ollama service..."
    ollama serve &
    sleep 3
fi

# Pull the model
echo "[INFO] Pulling model: $MODEL"
if ollama list | grep -q "$MODEL"; then
    echo "[OK] Model $MODEL is already downloaded"
else
    ollama pull "$MODEL"
    echo "[OK] Model $MODEL downloaded"
fi

# Verify
echo ""
echo "=== Verification ==="
echo "Ollama version: $(ollama --version)"
echo "Available models:"
ollama list

echo ""
echo "[OK] Ollama setup complete!"
