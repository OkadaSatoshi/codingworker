# CodingWorker - Go Worker

SQS からタスクを取得し、Aider + Ollama でコードを生成、GitHub PR を作成する Go Worker。

## アーキテクチャ

```
┌─────────────────┐
│    AWS SQS      │ ← メッセージキュー（または Mock）
└────────┬────────┘
         │ ロングポーリング
         ▼
┌─────────────────┐
│   Go Worker     │ ← このコンポーネント
└────────┬────────┘
         │ CLI呼び出し
         ▼
┌─────────────────┐
│  Aider + Ollama │ ← コード生成
└────────┬────────┘
         │ 自動コミット
         ▼
┌─────────────────┐
│  GitHub (PR)    │ ← プルリクエスト作成
└─────────────────┘
```

## ディレクトリ構成

```
worker/
├── cmd/
│   └── worker/
│       └── main.go      # エントリーポイント
├── internal/
│   ├── config/
│   │   └── config.go    # 設定読み込み
│   ├── sqs/
│   │   └── client.go    # SQS クライアント（Mock対応）
│   ├── aider/
│   │   └── runner.go    # Aider 実行
│   └── github/
│       └── client.go    # GitHub 操作
├── configs/
│   └── config.yaml      # 設定ファイル
├── go.mod
├── Taskfile.yml
└── README.md
```

## セットアップ

### 前提条件

- Go 1.23+
- Aider (`~/.local/bin/aider`)
- Ollama + qwen2.5-coder:1.5b
- gh CLI (GitHub CLI)
- Git

### 依存関係のインストール

```bash
task deps
```

### ビルド

```bash
task build
```

## 設定

`configs/config.yaml` を編集:

```yaml
sqs:
  use_mock: true  # 開発時は true、本番では false

aider:
  model: "ollama_chat/qwen2.5-coder:1.5b"
  bin_path: "${HOME}/.local/bin/aider"

github:
  token: "${GITHUB_TOKEN}"  # 環境変数から読み込み
```

### 環境変数

```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
```

## 実行

### 開発モード

```bash
task run:dev
```

### 本番モード

```bash
task run
```

## 開発状況

### 実装済み

- [x] 基本構造
- [x] 設定ファイル読み込み
- [x] SQS Mock クライアント
- [x] Aider Runner
- [x] GitHub Client (Clone/Branch/PR)

### TODO

- [ ] AWS SQS 実接続
- [ ] エラーハンドリング強化
- [ ] リトライロジック
- [ ] ログローテーション
- [ ] メトリクス収集
- [ ] 単体テスト
- [ ] 統合テスト

## ライセンス

Private
