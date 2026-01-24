#!/bin/bash
set -euo pipefail

# =============================================================================
# Aider Setup Script
# =============================================================================

echo "=== Aider Setup ==="

# Check Python version
echo "[INFO] Checking Python..."
if command -v python3 &> /dev/null; then
    PYTHON_VERSION=$(python3 --version)
    echo "[OK] Python installed: $PYTHON_VERSION"
else
    echo "[ERROR] Python 3 is not installed"
    echo "Please install Python 3.8+ first:"
    echo "  brew install python"
    exit 1
fi

# Check pip
if command -v pip3 &> /dev/null; then
    echo "[OK] pip3 is available"
else
    echo "[ERROR] pip3 is not installed"
    exit 1
fi

# Install/upgrade aider
echo "[INFO] Installing aider-chat..."
if command -v aider &> /dev/null; then
    echo "[INFO] Aider is already installed, upgrading..."
    pip3 install --upgrade aider-chat
else
    pip3 install aider-chat
fi

# Verify installation
echo ""
echo "=== Verification ==="
if command -v aider &> /dev/null; then
    echo "[OK] Aider installed: $(aider --version)"
else
    echo "[WARN] aider command not found in PATH"
    echo "You may need to add ~/.local/bin to your PATH"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "[OK] Aider setup complete!"
echo ""
echo "Usage with Ollama:"
echo "  aider --model ollama_chat/qwen2.5-coder:7b"
