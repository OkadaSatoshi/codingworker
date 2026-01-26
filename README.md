# CodingWorker

API費用ゼロの自動コーディングシステム

## 概要

MBP 2018 (Intel Mac) を自律型コーディングワーカーとして活用し、ローカルLLM（Ollama + Aider）でGitHub Issuesから自動的にコードを生成してPRを作成するシステム。

## アーキテクチャ

```
GitHub Issues → GitHub Actions → AWS SQS → Go Worker → Aider + Ollama → GitHub PR
```

## セットアップ

### 前提条件

- macOS (Monterey 12.0以上)
- Homebrew
- Git

### 1. 環境構築 (自動)

```bash
# Taskfile を使用した自動セットアップ
task setup
```

これにより以下がインストールされます:
- Ollama + qwen2.5-coder:1.5b
- mise (Goバージョン管理)
- uv + Aider

### 2. 環境構築 (手動)

#### Ollama

```bash
# インストール (https://ollama.ai からダウンロード)
ollama pull qwen2.5-coder:1.5b
ollama serve  # バックグラウンドで起動
```

#### Aider

```bash
# uv 経由でインストール
curl -LsSf https://astral.sh/uv/install.sh | sh
uv tool install --python 3.12 aider-chat

# PATH追加 (~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"
```

#### Go (mise経由)

```bash
brew install mise
mise install go
```

### 3. 動作確認

```bash
# Ollama
ollama run qwen2.5-coder:1.5b "Hello"

# Aider
~/.local/bin/aider --model ollama_chat/qwen2.5-coder:1.5b
```

### 4. Worker ビルド

```bash
cd worker
mise exec -- go build ./cmd/worker
```

## ドキュメント

| ドキュメント | 説明 |
|:---|:---|
| [要求書](docs/requirements.md) | プロジェクト要求事項 |
| [要件定義](docs/specifications.md) | システム仕様 |
| [詳細設計](docs/design.md) | Worker詳細設計 |
| [インフラ設計](docs/infrastructure.md) | AWS/GitHub Actions設計 |
| [タスク一覧](docs/tasks.md) | 進捗管理 |
| [PoC結果](poc/results/performance.md) | パフォーマンス計測結果 |

## 技術スタック

| コンポーネント | 技術 |
|:---|:---|
| コーディングエージェント | Aider |
| ローカルLLM | Ollama (qwen2.5-coder:1.5b) |
| メッセージキュー | AWS SQS |
| ワーカー | Go |
| IaC | Terraform |

## ステータス

**現在のフェーズ**: Phase 2 - Worker開発 (90%完了)

| フェーズ | 状態 |
|:---|:---|
| Phase 0: PoC | ✅ 完了 |
| Phase 1: AWS基盤 | 🟡 コード完了 (apply待ち) |
| Phase 2: Worker | 🟡 90% (統合テスト残) |
| Phase 3: E2Eテスト | ⚪ 未着手 |
| Phase 4: 運用整備 | ⚪ 未着手 |

## ライセンス

MIT
