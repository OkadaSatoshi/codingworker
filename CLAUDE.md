# CodingWorker - AI Assistant Instructions

## Project Overview

CodingWorker is an automated coding system using local LLM with zero API cost.
Goal: GitHub Issues → GitHub Actions → AWS SQS → Go Worker → Aider + Ollama → GitHub PR

## Architecture Decisions

### Tool Selection
- **Aider** (not OpenHands): Chosen for lighter weight, no Docker required, CLI-based
- **Ollama**: Local LLM runtime
- **Model**: `qwen2.5-coder:1.5b` (for MBP 2018 compatibility, CPU inference optimized)
- **uv**: Python package manager (not pip) for isolated Aider installation
- **mise**: Tool version manager for Go
- **Taskfile**: Task automation (not Makefile)

### Repository Structure
- **Monorepo**: Single repository for AI-coding friendliness
- `docs/` - Project documentation (English filenames)
- `poc/` - Proof of concept test cases
- `scripts/` - Automation scripts
- `worker/` - Go Worker (SQS polling, Aider execution, GitHub PR creation)

## Worker Component

### Directory Structure
```
worker/
├── cmd/worker/main.go      # Entry point
├── internal/
│   ├── config/config.go    # Configuration loader
│   ├── sqs/client.go       # SQS client (Mock supported)
│   ├── aider/runner.go     # Aider CLI executor
│   └── github/client.go    # GitHub operations
├── configs/config.yaml     # Configuration file
├── go.mod
└── Taskfile.yml
```

### Building Worker
```bash
cd worker
mise exec -- go mod tidy
mise exec -- go build ./cmd/worker
```

### Worker Configuration
`worker/configs/config.yaml`:
- `sqs.use_mock: true` - Development mode (no AWS required)
- `github.token` - Uses `$GITHUB_TOKEN` environment variable
- `aider.bin_path` - Path to Aider binary

### Current Status
- SQS: Mock implementation (AWS integration pending Phase 1)
- Aider: CLI execution implemented
- GitHub: Clone/Branch/PR via git and gh CLI

## Key Configuration Files

### Aider Model Configuration
`.aider.model.metadata.json` - Required to suppress "Unknown context window size" warning:
```json
{
    "ollama_chat/qwen2.5-coder:1.5b": {
        "max_tokens": 32768,
        "max_input_tokens": 32768,
        "max_output_tokens": 8192,
        "input_cost_per_token": 0,
        "output_cost_per_token": 0
    }
}
```

### Python Version
`.python-version` contains `3.12` for:
- Dependabot compatibility
- uv tool install version specification
- No need to install Python via mise (uv handles it)

### Go Version
`mise.toml` specifies Go version for mise to manage.

## Environment Setup

Run `task setup` for full environment setup (idempotent).

Key tasks:
- `task setup:ollama` - Install Ollama and download model
- `task setup:uv` - Install uv package manager
- `task setup:aider` - Install Aider via uv with Python 3.12
- `task setup:go` - Install Go via mise
- `task setup:path` - Add ~/.local/bin to PATH in .zshrc

## PoC Test Cases

All test cases use Go:
1. `poc/test-case-1-fizzbuzz/` - FizzBuzz generation
2. `poc/test-case-2-csv/` - CSV sorting by name
3. `poc/test-case-3-bugfix/` - Bug detection and fix (off-by-one error)

Run with: `task poc:fizzbuzz`, `task poc:csv`, `task poc:bugfix`

## Important Notes

### Aider Behavior
- Aider uses git working directory as base
- Files created by Aider may appear in git root, not subdirectories
- Must run Aider from the specific test case directory

### Running Aider
```bash
~/.local/bin/aider --model ollama_chat/qwen2.5-coder:1.5b
```

### Running Go with mise
```bash
mise exec -- go run main.go
mise exec -- go test -v
```

## Cleanup

Run `task cleanup` or follow `docs/cleanup.md` to remove:
- Ollama service and data
- Aider and cache
- Downloaded models

## Target Environment

- Primary: MBP 2018 (Intel, limited resources)
- PoC validated on: MBP M4
- OS: macOS (Linux support planned but not implemented)
