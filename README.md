# CodingWorker

API費用ゼロの自動コーディングシステム

## 概要

MBP 2018 (Intel Mac) を自律型コーディングワーカーとして活用し、ローカルLLM（Ollama + Aider）でGitHub Issuesから自動的にコードを生成してPRを作成するシステム。

## アーキテクチャ

```
GitHub Issues → GitHub Actions → AWS SQS → Go Worker → Aider + Ollama → GitHub PR
```

## ドキュメント

- [要求書](docs/requirements.md)
- [要件定義](docs/specifications.md)
- [タスク一覧](docs/tasks.md)

## 技術スタック

| コンポーネント | 技術 |
|:---|:---|
| コーディングエージェント | Aider |
| ローカルLLM | Ollama (qwen2.5-coder:1.5b) |
| メッセージキュー | AWS SQS |
| ワーカー | Go |
| IaC | Terraform |

## ステータス

**現在のフェーズ**: Phase 0 - ローカル環境 PoC

## ライセンス

MIT
