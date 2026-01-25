# 要件定義: CodingWorker - 無料エージェントコーディングシステム

## 1. ドキュメント情報

- **プロジェクト名**: CodingWorker
- **ドキュメント種別**: 要件定義書
- **バージョン**: 1.0
- **作成日**: 2025-01-25
- **承認者**: OkadaSatoshi

---

## 2. システム全体構成

### 2.1 システムアーキテクチャ

```
┌─────────────────┐
│  GitHub Issues  │ ← ユーザーがタスク起票
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ GitHub Actions  │ ← イベントトリガー
└────────┬────────┘
         │ OIDC認証
         ▼
┌─────────────────┐
│    AWS SQS      │ ← メッセージキュー（非同期処理）
└────────┬────────┘
         │ ロングポーリング
         ▼
┌─────────────────┐
│   Go Worker     │ ← MBP上で稼働
└────────┬────────┘
         │ Docker API
         ▼
┌─────────────────┐
│   OpenHands     │ ← コーディングエージェント
└────────┬────────┘
         │ API呼び出し
         ▼
┌─────────────────┐
│     Ollama      │ ← ローカルLLM (qwen2.5-coder:1.5b)
└────────┬────────┘
         │
         ▼
    コード生成
         │
         ▼
┌─────────────────┐
│  GitHub (PR)    │ ← プルリクエスト作成
└─────────────────┘
```

### 2.2 技術スタック

| レイヤー | 技術 | バージョン | 備考 |
|:---|:---|:---|:---|
| **クラウド** | AWS | - | 無料枠活用 |
| **メッセージング** | SQS | Standard Queue | 暗号化有効 |
| **IaC** | Terraform | 1.5+ | State管理: S3 |
| **CI/CD** | GitHub Actions | - | OIDC認証 |
| **ワーカー** | Go | 1.21+ | AWS SDK v2 |
| **コンテナ** | Docker | 24.0+ | Docker Desktop for Mac |
| **エージェント** | OpenHands | Latest | Headless Mode |
| **LLM** | Ollama | 0.1.20+ | qwen2.5-coder:1.5b |
| **ハードウェア** | MacBook Pro 2018 | Intel i7, 16GB RAM | - |

---

## 3. 機能要件

### 3.1 タスク起票機能

#### FR-001: GitHub Issues によるタスク起票

**説明**: ユーザーが GitHub Issues にタスクを起票することで、自動コーディングプロセスを開始する。

**要件**:
- [ ] Issue作成時に GitHub Actions がトリガーされる
- [ ] Issueのタイトルと本文をSQSメッセージに変換する
- [ ] 特定のラベル（例: `auto-code`）が付与されたIssueのみを対象とする

**入力**:
- Issueタイトル
- Issue本文（Markdown形式）
- ラベル

**出力**:
- SQSメッセージ（JSON形式）

**非機能要件**:
- 処理時間: 30秒以内
- 信頼性: メッセージ送信成功率 99%以上

---

### 3.2 メッセージキュー機能

#### FR-002: AWS SQS によるメッセージ管理

**説明**: GitHub Actions から送信されたタスクメッセージをキューイングし、ローカルワーカーに配信する。

**要件**:
- [ ] Standard Queue を使用（FIFO不要）
- [ ] メッセージの可視性タイムアウト: 3600秒（1時間）
- [ ] 最大受信回数: 3回
- [ ] Dead Letter Queue（DLQ）への自動転送

**メッセージフォーマット**:
```json
{
  "issue_number": 123,
  "repository": "owner/repo",
  "title": "Implement user authentication",
  "body": "Add JWT-based authentication...",
  "labels": ["auto-code", "backend"],
  "created_at": "2025-01-25T10:00:00Z"
}
```

**非機能要件**:
- メッセージ保持期間: 4日間
- 暗号化: AWS管理キー（SSE-SQS）
- スループット: 100メッセージ/時

---

### 3.3 ワーカー機能

#### FR-003: Go Worker によるタスク処理

**説明**: MBP上で稼働するGoプログラムが、SQSをポーリングし、タスクを取得・処理する。

**要件**:
- [ ] SQS ロングポーリング（WaitTimeSeconds: 20）
- [ ] メッセージ取得後、OpenHandsコンテナを起動
- [ ] 処理完了後、メッセージを削除
- [ ] エラー発生時、メッセージを再キューイング（最大3回）

**処理フロー**:
```
1. SQS からメッセージ取得
2. メッセージをパース
3. Docker API 経由で OpenHands コンテナ起動
4. OpenHands にタスクを渡す
5. 完了まで待機（最大1時間）
6. 結果を取得
7. GitHub へ Push & PR作成
8. SQS メッセージ削除
```

**非機能要件**:
- 同時実行: 1タスクのみ
- メモリ使用量: 2GB以下（Worker自体）
- ログレベル: INFO以上

---

### 3.4 コード生成機能

#### FR-004: OpenHands + Ollama によるコード生成

**説明**: OpenHandsフレームワークを使用し、ローカルOllama（qwen2.5-coder:1.5b）でコードを生成する。

**要件**:
- [ ] Headless Mode で動作
- [ ] Ollama エンドポイント: `http://localhost:11434`
- [ ] モデル: `qwen2.5-coder:1.5b`
- [ ] タイムアウト: 3600秒（1時間）

**入力**:
- タスク説明（Issue本文）

**出力**:
- 生成されたコードファイル
- 変更の説明（コミットメッセージ）

**非機能要件**:
- 推論速度: 10-30秒/応答（Intel Mac想定）
- メモリ使用量: 最大10GB（Ollama + OpenHands）
- 品質: 生成コードが構文エラーなしで実行可能

---

### 3.5 GitHub連携機能

#### FR-005: GitHub へのコミット・PR作成

**説明**: 生成されたコードをGitHubリポジトリにコミットし、プルリクエストを作成する。

**要件**:
- [ ] 新しいブランチを作成（命名規則: `auto-code/issue-{number}`）
- [ ] コードをコミット
- [ ] PRを作成（Issue番号を参照）
- [ ] PR本文に生成プロセスの説明を含める

**認証**:
- Personal Access Token（スコープ: `repo`, `workflow`）

**PRテンプレート**:
```markdown
## 自動生成されたコード

このPRは CodingWorker によって自動生成されました。

**関連Issue**: #123

**生成モデル**: Ollama qwen2.5-coder:1.5b  
**生成日時**: 2025-01-25 10:30:00

### 変更内容
- ファイル1: 説明
- ファイル2: 説明

### 確認事項
- [ ] コードが正しく動作するか
- [ ] テストが通るか
- [ ] コーディング規約に準拠しているか
```

**非機能要件**:
- 処理時間: 60秒以内
- 失敗時のリトライ: 3回

---

## 4. 非機能要件

### 4.1 パフォーマンス要件

| 項目 | 目標値 | 計測方法 |
|:---|:---|:---|
| タスク処理時間 | 5分以内（平均） | ログから計測 |
| スループット | 10-20タスク/日 | 日次レポート |
| 推論速度 | 10-30秒/応答 | Ollama ログ |
| メモリ使用量 | 12GB以下（ピーク） | Activity Monitor |
| CPU使用率 | 80%以下（平均） | Activity Monitor |

### 4.2 信頼性要件

| 項目 | 目標値 | 対策 |
|:---|:---|:---|
| 稼働率 | 80%以上 | 自動再起動設定 |
| メッセージ損失率 | 0% | DLQ + リトライ |
| エラーリカバリ | 自動リトライ3回 | Worker実装 |
| 長時間稼働 | 3時間以上 | メモリリーク対策 |

### 4.3 セキュリティ要件

#### SEC-001: 認証・認可

| 接続 | 認証方式 | 権限 |
|:---|:---|:---|
| GitHub Actions → AWS | IAM OIDC | SQS SendMessage |
| MBP → AWS SQS | IAM User (Access Key) | SQS ReceiveMessage, DeleteMessage |
| MBP → GitHub | Personal Access Token | repo, workflow |

#### SEC-002: データ保護

- [ ] SQS メッセージの暗号化（at-rest）: AWS管理キー
- [ ] SQS メッセージの暗号化（in-transit）: TLS 1.2+
- [ ] GitHub PAT の安全な保管: macOS Keychain または環境変数

#### SEC-003: アクセス制御

- [ ] IAM ロールの最小権限原則
- [ ] GitHub PAT のスコープ最小化
- [ ] ログに機密情報を記録しない

### 4.4 保守性要件

#### MAINT-001: ログ管理

**要件**:
- [ ] すべての処理をログに記録
- [ ] ログレベル: DEBUG, INFO, WARN, ERROR
- [ ] ログローテーション: 7日間保持、自動削除
- [ ] ログフォーマット: JSON形式

**ログ出力先**:
- Go Worker: `/var/log/bridge-coder/worker.log`
- OpenHands: Docker ログ
- Ollama: `~/.ollama/logs/`

#### MAINT-002: モニタリング

**要件**:
- [ ] CloudWatch Metrics: SQS キュー深度、メッセージ処理時間
- [ ] ローカルメトリクス: CPU, メモリ, ディスクI/O
- [ ] アラート設定: DLQメッセージ数 > 0

#### MAINT-003: ドキュメント

**必須ドキュメント**:
- [ ] システムアーキテクチャ図
- [ ] セットアップ手順書
- [ ] トラブルシューティングガイド
- [ ] API仕様書（SQSメッセージフォーマット）
- [ ] 運用マニュアル

### 4.5 スケーラビリティ要件

**現状**:
- 単一ワーカー、同時1タスク

**将来の拡張**:
- [ ] 複数ワーカーの並行実行
- [ ] 複数プロジェクトの同時処理
- [ ] クラウドLLMへのフォールバック

---

## 5. インターフェース定義

### 5.1 GitHub Actions → SQS

**エンドポイント**: 
```
https://sqs.{region}.amazonaws.com/{account-id}/{queue-name}
```

**認証**: IAM OIDC

**リクエスト形式**:
```json
{
  "MessageBody": "{...JSON...}",
  "MessageAttributes": {
    "IssueNumber": {
      "StringValue": "123",
      "DataType": "Number"
    },
    "Repository": {
      "StringValue": "owner/repo",
      "DataType": "String"
    }
  }
}
```

### 5.2 Go Worker → OpenHands

**インターフェース**: Docker API

**起動コマンド例**:
```bash
docker run -it \
  -e OLLAMA_BASE_URL=http://host.docker.internal:11434 \
  -e LLM_MODEL=qwen2.5-coder:1.5b \
  -v $(pwd)/workspace:/workspace \
  openhands/openhands:latest \
  --task "Implement user authentication"
```

### 5.3 OpenHands → Ollama

**エンドポイント**: `http://localhost:11434/api/generate`

**リクエスト形式**:
```json
{
  "model": "qwen2.5-coder:1.5b",
  "prompt": "Create a Python function...",
  "stream": false
}
```

**レスポンス形式**:
```json
{
  "response": "def function_name():\n    ...",
  "done": true
}
```

---

## 6. データ構造

### 6.1 SQSメッセージ構造

```json
{
  "issue_number": 123,
  "repository": "owner/repo",
  "title": "Implement feature X",
  "body": "Description of the task...",
  "labels": ["auto-code", "enhancement"],
  "assignees": ["username"],
  "created_at": "2025-01-25T10:00:00Z",
  "updated_at": "2025-01-25T10:00:00Z"
}
```

### 6.2 ワーカーログ構造

```json
{
  "timestamp": "2025-01-25T10:30:00Z",
  "level": "INFO",
  "message": "Task processing started",
  "issue_number": 123,
  "repository": "owner/repo",
  "worker_id": "mbp-2018-001"
}
```

### 6.3 コミットメッセージ構造

```
[auto-code] Implement user authentication

- Added JWT token generation
- Implemented middleware for auth
- Added tests for authentication flow

Resolves #123
Generated by CodingWorker (Ollama qwen2.5-coder:1.5b)
```

---

## 7. エラーハンドリング

### 7.1 エラー分類

| エラー種別 | 例 | 対応 |
|:---|:---|:---|
| **一時的エラー** | ネットワークタイムアウト | 自動リトライ（最大3回） |
| **恒久的エラー** | 認証失敗、不正なメッセージ | DLQへ移動、アラート |
| **タイムアウト** | 1時間以上かかるタスク | 処理中断、DLQへ移動 |
| **リソース不足** | メモリ不足、ディスク容量不足 | 処理中断、アラート |

### 7.2 リトライポリシー

```
試行1: 即座にリトライ
試行2: 60秒待機後にリトライ
試行3: 300秒待機後にリトライ
失敗 : DLQへ移動
```

### 7.3 Dead Letter Queue（DLQ）

**用途**: 処理に失敗したメッセージを保管

**保持期間**: 14日間

**アラート**: DLQにメッセージが1件以上ある場合、SNS経由でメール通知

---

## 8. 制約事項

### 8.1 技術的制約

- Intel Mac の推論速度は遅い（10-30秒/応答）
- 同時実行は1タスクのみ
- モデルのコンテキスト長制限（qwen2.5-coder:1.5b は32K tokens）

### 8.2 運用上の制約

- MBP の電源が入っている必要がある
- ネットワーク接続が必須
- Docker Desktop が起動している必要がある

### 8.3 コスト制約

- AWS 無料枠内での運用（SQS: 100万リクエスト/月）
- 電気代のみ（月額2,000円程度）

---

## 9. テスト要件

### 9.1 単体テスト

- [ ] Go Worker の各関数をテスト
- [ ] SQSメッセージのパース処理
- [ ] GitHub API 呼び出し処理

### 9.2 統合テスト

- [ ] GitHub Actions → SQS → Worker の連携
- [ ] Worker → OpenHands → Ollama の連携
- [ ] コード生成 → GitHub PR作成

### 9.3 E2Eテスト

- [ ] Issue起票からPR作成までの完全フロー
- [ ] エラーケース（タイムアウト、認証失敗等）
- [ ] 長時間稼働テスト（3時間以上）

### 9.4 パフォーマンステスト

- [ ] 推論速度の計測
- [ ] メモリ使用量の計測
- [ ] CPU使用率の計測

---

## 10. 移行・展開計画

### 10.1 Phase 0: PoC（3-5日）

- [ ] Ollama セットアップ
- [ ] OpenHands セットアップ
- [ ] ローカル動作確認

### 10.2 Phase 1: Cloud Foundation（2-3日）

- [ ] Terraform による IaC化
- [ ] IAM OIDC 設定
- [ ] SQS キュー作成

### 10.3 Phase 2: ワーカー開発（3-5日）

- [ ] Go Worker 実装
- [ ] GitHub Actions ワークフロー作成
- [ ] Docker API 連携

### 10.4 Phase 3: 統合テスト（2-3日）

- [ ] E2Eテスト
- [ ] バグ修正
- [ ] ドキュメント整備

### 10.5 Phase 4: 運用開始（2-3日）

- [ ] モニタリング設定
- [ ] アラート設定
- [ ] 運用マニュアル作成

---

## 11. 用語集

| 用語 | 説明 |
|:---|:---|
| **OpenHands** | オープンソースのコーディングエージェントフレームワーク |
| **Ollama** | ローカルLLMを実行するためのツール |
| **qwen2.5-coder:1.5b** | Alibabaが開発したコーディング特化型LLM（70億パラメータ） |
| **OIDC** | OpenID Connect。GitHub ActionsがAWSに認証するための仕組み |
| **DLQ** | Dead Letter Queue。処理失敗メッセージの保管場所 |
| **IaC** | Infrastructure as Code。インフラをコードで管理する手法 |

---

**承認日**: 2025-01-25  
**バージョン**: 1.0  
**次回レビュー**: Phase 0 完了後
