#!/bin/bash
set -euo pipefail

# =============================================================================
# Full Setup Script - Ollama + Aider
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=============================================="
echo "CodingWorker - Full Setup"
echo "=============================================="
echo ""

# Run Ollama setup
"$SCRIPT_DIR/setup-ollama.sh"
echo ""

# Run Aider setup
"$SCRIPT_DIR/setup-aider.sh"
echo ""

echo "=============================================="
echo "Setup Complete!"
echo "=============================================="
echo ""
echo "Next steps:"
echo "  1. cd poc/test-case-1-fizzbuzz"
echo "  2. git init"
echo "  3. aider --model ollama_chat/qwen2.5-coder:7b"
echo ""
