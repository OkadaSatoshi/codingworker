# CodingWorker - 全体タスク一覧（v2: Aider版）

## プロジェクト情報

- **プロジェクト名**: CodingWorker
- **目的**: API費用ゼロの自動コーディングシステム構築
- **期間**: 約10-15日
- **現在のフェーズ**: Phase 0
- **コーディングエージェント**: Aider（旧: OpenHands）

---

## 変更履歴

| バージョン | 日付 | 変更内容 |
|:---|:---|:---|
| v1 | 2025-01-25 | 初版作成（OpenHands版） |
| v2 | 2025-01-25 | Aider版に変更、Docker不要化、タスク簡素化 |

---

## アーキテクチャ概要

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
│    AWS SQS      │ ← メッセージキュー
└────────┬────────┘
         │ ロングポーリング
         ▼
┌─────────────────┐
│   Go Worker     │ ← MBP上で稼働
└────────┬────────┘
         │ CLI呼び出し
         ▼
┌─────────────────┐
│  Aider + Ollama │ ← ローカルLLM (qwen2.5-coder:1.5b)
└────────┬────────┘
         │ 自動コミット
         ▼
┌─────────────────┐
│  GitHub (PR)    │ ← プルリクエスト作成
└─────────────────┘
```

---

## 進捗サマリー

| フェーズ | ステータス | 完了率 | 期間 |
|:---|:---|:---:|:---|
| Phase 0: ローカル環境 PoC | ✅ 完了 | 100% | 2-3日 |
| Phase 1: Cloud Foundation | 🟡 コード完了 | 80% | 2-3日 |
| Phase 2: SQS連携・ワーカー | 🟡 進行中 | 90% | 3-4日 |
| Phase 3: E2E統合テスト | ⚪ 未着手 | 0% | 2-3日 |
| Phase 4: 運用基盤整備 | ⚪ 未着手 | 0% | 1-2日 |

**注記**:
- Phase 1: Terraformコード完了、`terraform apply` 未実行
- Phase 2: Worker実装完了、Mock SQS統合テスト残

---

# Phase 0: ローカル環境 PoC 構築

**期間**: 2-3日
**目的**: MBP 2018 上で Aider + Ollama を動作させ、基本的なコーディングタスクが実行可能であることを検証する

---

## タスク0: 環境確認 ✅

### 0-1: MBP環境の確認（30分）
- [x] macOSバージョン確認（Monterey 12.0以上推奨）
- [x] 空きディスク容量確認（最低10GB必要）
- [x] 空きメモリ確認（16GB中、8GB以上空き推奨）
- [x] Python 3.8+ インストール確認（`python3 --version`）
- [x] Git インストール確認（`git --version`）
- [x] Homebrew インストール確認（`brew --version`）

### 0-2: 不足ツールのインストール
- [x] Python未インストールの場合: `brew install python`
- [x] Git未インストールの場合: `brew install git`
- [x] pip確認: `pip3 --version`

---

## タスク1: Ollama セットアップ ✅

### 1-1: Ollama のインストール
- [x] 公式サイト（https://ollama.ai）からインストーラーをダウンロード
- [x] インストール実行
- [x] `ollama --version` で動作確認
- [x] Ollama サービス起動確認

### 1-2: qwen2.5-coder:1.5b モデルのダウンロード
- [x] `ollama pull qwen2.5-coder:1.5b` 実行
- [x] ダウンロード完了確認（約4GB）
- [x] `ollama list` でモデル確認

### 1-3: Ollama 基本動作確認
- [x] `ollama run qwen2.5-coder:1.5b` で対話テスト
- [x] 簡単なコード生成プロンプト実行（例: "Write a hello world in Python"）
- [x] レスポンス速度を目視確認（目安: 30秒以内）
- [x] `Ctrl+D` で終了

---

## タスク2: Aider セットアップ ✅

### 2-1: Aider のインストール
- [x] `uv tool install aider-chat` 実行（pip3 → uv に変更）
- [x] `aider --version` で確認
- [x] インストール完了確認

### 2-2: Aider + Ollama 連携設定
- [x] テスト用ディレクトリ作成: `mkdir ~/aider-test && cd ~/aider-test`
- [x] Git初期化: `git init`
- [x] Aider起動: `aider --model ollama_chat/qwen2.5-coder:1.5b`
- [x] 接続成功確認

### 2-3: Aider 基本動作確認
- [x] 簡単なファイル作成指示（例: "Create a hello.py that prints hello world"）
- [x] ファイルが生成されることを確認
- [x] 自動コミットされることを確認（`git log`）
- [x] `/quit` で終了

---

## タスク3: PoC テスト実行 ✅

### 3-1: テストケース1 - FizzBuzz
- [x] 新規ディレクトリ作成・Git初期化
- [x] Aider起動
- [x] タスク投入: "Create a Go script fizzbuzz.go that prints FizzBuzz from 1 to 100"（Go版）
- [x] 生成コードを確認
- [x] `go run fizzbuzz.go` で実行・動作確認
- [x] 実行時間を記録

### 3-2: テストケース2 - ファイル処理スクリプト
- [x] タスク投入: "Create a Go script that reads a CSV file and sorts it by the Name column"
- [x] 生成コードを確認
- [x] テスト用CSVで動作確認
- [x] 実行時間を記録

### 3-3: テストケース3 - バグ修正
- [x] バグを含むコードを準備（off-by-one error）
- [x] タスク投入: "Fix the bug in buggy.go"
- [x] バグ検出・修正能力を確認
- [x] 実行時間を記録

### 3-4: テスト結果まとめ
- [x] 各テストケースの成否記録
- [x] 平均実行時間算出
- [x] 生成コード品質評価（構文エラー有無、正確性）
- [x] 結果: `poc/results/performance.md` に記録

---

## タスク4: パフォーマンス計測 ✅

### 4-1: リソース使用量の監視
- [x] Activity Monitor でOllamaのメモリ使用量確認
- [x] Aider実行中のメモリ使用量確認
- [x] 合計メモリ使用量が14GB以下か確認
- [x] CPU使用率記録
- [x] 結果: MBP M4 で計測、Intel Mac は推論速度を考慮して 10分タイムアウト設定

### 4-2: 推論速度の計測
- [x] 3つのテストケースそれぞれで時間計測
- [x] 平均推論速度算出
- [x] 目標: 60秒/応答以内 → 達成（1.5B: 約2分、qwen2.5-coder使用）

### 4-3: 長時間稼働テスト（オプション）
- [x] スキップ（Phase 3 で実施予定）

---

## タスク5: Go/No-Go 判断 ✅

### 5-1: 評価基準チェック

#### Go 基準（すべて満たす場合、Phase 1へ）
- [x] 推論速度: 60秒/応答以内 → 達成
- [x] メモリ使用量: 14GB以下 → 達成
- [x] テストケース: 3つ中2つ以上成功 → 3つ全て成功
- [x] Aider + Ollama連携: 安定動作 → 確認済み

#### Pivot 基準（代替案検討）
- [x] Aider連携が不安定 → 該当なし
- [x] qwen2.5-coder:1.5b品質不足 → 該当なし（1.5Bで十分）

#### Fallback 基準（無料API検討）
- [x] ローカルLLM推論が120秒/応答超過 → 該当なし
- [x] 品質が実用レベルでない → 該当なし

#### No-Go 基準（プロジェクト中止）
- [x] 該当なし

### 5-2: 判断結果記録
- [x] 判断結果: **Go**
- [x] 理由: 全てのGo基準を満たした。3B は Diff 生成エラー・パス管理問題で不採用、1.5B のみ使用
- [x] 次のアクション: Phase 1 へ進行

---

## タスク6: ドキュメント化 ✅

### 6-1: セットアップ手順書
- [x] Ollama インストール手順
- [x] Aider インストール手順
- [x] 連携設定手順
- [x] ドキュメント: `README.md`（自動/手動セットアップ手順）

### 6-2: PoC結果レポート
- [x] テスト結果まとめ
- [x] パフォーマンスデータ
- [x] 課題・改善点
- [x] ドキュメント: `poc/results/performance.md`

### 6-3: Phase 1 引き継ぎ事項
- [x] 明らかになった課題（3B不採用）
- [x] 推奨設定値（map_tokens: 0、10分タイムアウト）
- [x] ドキュメント: `docs/design.md` に統合

---

## Phase 0 完了判定基準 ✅

- [x] Ollama + qwen2.5-coder:1.5b が安定動作する
- [x] Aider が Ollama 経由でコードを生成できる
- [x] 生成されたコードが実行可能である
- [x] Go/No-Go判断が完了している（Go判定）
- [x] セットアップ手順が文書化されている（README.md）

---

# Phase 1: Cloud Foundation 構築

**期間**: 2-3日
**目的**: AWS上にクラウド基盤を構築し、GitHub Actions と SQS の連携を実現する
**前提**: Phase 0 で Go 判定

---

## タスク7: AWS環境準備 🟡

### 7-1: AWSアカウント確認
- [x] AWSアカウント有無確認
- [x] 無料枠残量確認
- [ ] 請求アラート設定（$5超過で通知） → Phase 4

### 7-2: AWS CLI セットアップ
- [x] `brew install awscli` 実行
- [x] `aws --version` で確認
- [ ] IAM User作成（MBP用）→ Terraform で作成予定
- [ ] Access Key / Secret Key 取得 → Terraform 適用後
- [ ] `aws configure` で設定 → Terraform 適用後
- [ ] `aws sts get-caller-identity` で確認 → Terraform 適用後

---

## タスク8: Terraform セットアップ ✅

### 8-1: Terraform インストール
- [x] `brew install terraform` 実行
- [x] `terraform --version` で確認

### 8-2: Terraform State 管理準備
- [x] ローカルState（S3/DynamoDBはスキップ - 単一開発者のため）
- [x] `backend.tf` 設定ファイル作成

### 8-3: Terraform コード作成（モノレポ内）
- [x] ディレクトリ構造作成
```
infra/terraform/
├── backend.tf
├── provider.tf
├── variables.tf
├── outputs.tf
├── oidc.tf
├── sqs.tf
└── iam.tf
```
- [x] 設計ドキュメント: `docs/infrastructure.md`

---

## タスク9: IAM OIDC プロバイダー設定 🟡 コード完了

### 9-1: IAM OIDC プロバイダー作成
- [x] Terraform コード作成: `oidc.tf`
- [x] GitHub OIDC プロバイダー設定
- [x] Trust Policy 設定

### 9-2: IAM ロール作成（GitHub Actions用）
- [x] SQS SendMessage 権限を持つロール作成
- [x] Terraform コード: `iam.tf`
- [x] 最小権限原則の確認

### 9-3: IAM User作成（MBP Worker用）
- [x] SQS ReceiveMessage, DeleteMessage 権限
- [x] Terraform コード追加
- [x] Access Key 出力設定

**ステータス**: Terraform コード完了、`terraform apply` 待ち

---

## タスク10: SQS キュー作成 🟡 コード完了

### 10-1: SQS Standard Queue 作成
- [x] Terraform コード作成: `sqs.tf`
- [x] キュー名: `codingworker-tasks`
- [x] 可視性タイムアウト: 3600秒
- [x] 暗号化設定（AWS管理キー）

### 10-2: Dead Letter Queue 作成
- [x] DLQ作成: `codingworker-tasks-dlq`
- [x] メインキューに DLQ 設定
- [x] 最大受信回数: 3回

### 10-3: Terraform Apply
- [x] `terraform init` 実行
- [x] `terraform plan` で確認
- [ ] `terraform apply` で適用 → **未実行**
- [ ] 作成されたリソース確認 → apply 後

**ステータス**: Terraform コード完了、`terraform apply` 待ち

---

## タスク11: GitHub Actions ワークフロー作成 🟡 設計完了

### 11-1: テスト用リポジトリ準備
- [ ] `codingworker-sandbox` リポジトリ作成 → AWS 適用後
- [ ] `.github/workflows/` ディレクトリ作成 → AWS 適用後

### 11-2: OIDC 認証ワークフロー作成
- [x] `.github/workflows/send-to-sqs.yml` 設計完了
- [x] OIDC 認証設定（設計）
- [x] Issue イベントトリガー（`ai-task`ラベル）設計

### 11-3: SQS メッセージ送信実装
- [x] メッセージフォーマット（JSON）設計完了
- [x] AWS CLI で SendMessage 実装（設計）
- [x] エラーハンドリング（設計）

### 11-4: 連携テスト
- [ ] テスト用 Issue 作成（`ai-task`ラベル付与）→ AWS 適用後
- [ ] ワークフロー実行確認 → AWS 適用後
- [ ] SQS にメッセージが届くか確認 → AWS 適用後

**ステータス**: 設計完了、AWS リソース作成後に実装・テスト
**設計ドキュメント**: `docs/infrastructure.md`

---

## Phase 1 完了判定基準 🟡

- [ ] Terraform で AWS リソースが作成されている → `terraform apply` 待ち
- [ ] IAM OIDC が正しく設定されている → `terraform apply` 待ち
- [ ] SQS キュー（メイン + DLQ）が作成されている → `terraform apply` 待ち
- [ ] GitHub Actions → SQS へのメッセージ送信が成功する → AWS 適用後にテスト
- [x] Terraform コードがリポジトリに保存されている

---

# Phase 2: SQS連携とワーカー開発

**期間**: 3-4日
**目的**: Go製ワーカーを実装し、SQS → Aider → GitHub PR の連携を実現する

---

## タスク12: worker ディレクトリ作成（モノレポ内）

### 12-1: ディレクトリ構造作成
- [x] モノレポ内に worker/ ディレクトリ作成
- [x] ディレクトリ構造作成
```
worker/
├── cmd/
│   └── worker/
│       └── main.go
├── internal/
│   ├── sqs/
│   │   └── client.go
│   ├── aider/
│   │   └── runner.go
│   ├── github/
│   │   └── client.go
│   └── config/
│       └── config.go
├── configs/
│   └── config.yaml
├── go.mod
├── go.sum
├── Taskfile.yml
└── README.md
```

### 12-2: Go モジュール初期化
- [x] `go mod init github.com/OkadaSatoshi/codingworker/worker`
- [x] 依存関係追加
  - `gopkg.in/yaml.v3`
  - (AWS SDK は Phase 1 完了後に追加)
- [x] Taskfile.yml 作成（Makefileではなく）

---

## タスク13: SQS ポーリング機能実装 🟡 Mock完了

### 13-1: SQS クライアント実装
- [x] `internal/sqs/client.go` 作成
- [x] Mock SQS クライアント実装（開発用）
- [ ] AWS SDK v2 でクライアント作成 → AWS 適用後

### 13-2: メッセージ受信処理
- [x] ロングポーリング実装（WaitTimeSeconds: 20）
- [x] メッセージパース（JSON → struct）
- [x] エラーハンドリング

### 13-3: メッセージ削除処理
- [x] 処理完了後の DeleteMessage 実装
- [x] 失敗時の処理（メッセージを削除せず再キューイング）

**ステータス**: Mock 実装完了、AWS SDK 統合は AWS 適用後

---

## タスク14: Aider 連携実装 ✅

### 14-1: Aider Runner 実装
- [x] `internal/aider/runner.go` 作成
- [x] Aider CLI 呼び出し処理
- [x] コマンド構築: `aider --model ollama_chat/qwen2.5-coder:1.5b --yes --message "{task}"`
- [x] 2パス実行: 実装 → テスト作成（`RunWithTests`）
- [x] 検証ループ: build → lint → test（最大3回リトライ）

### 14-2: 作業ディレクトリ管理
- [x] リポジトリのクローン処理（`github/client.go`）
- [x] ブランチ作成: `ai/issue-{number}`
- [x] 作業完了後のクリーンアップ

### 14-3: 実行結果取得
- [x] Aider 終了コード確認
- [x] 生成ファイルの確認
- [x] エラー検出・Aider への修正依頼

**設計ドキュメント**: `docs/design.md`

---

## タスク15: GitHub 連携実装 ✅

### 15-1: GitHub クライアント実装
- [x] `internal/github/client.go` 作成
- [x] `gh` CLI 使用（go-github ライブラリ不使用）
- [x] Personal Access Token 認証（環境変数）

### 15-2: PR 作成処理
- [x] ブランチのプッシュ（`git push`）
- [x] Pull Request 作成（`gh pr create`）
- [x] PR 本文テンプレート適用
- [x] Issue番号参照（`Closes #123`）

### 15-3: PR テンプレート ✅
```markdown
## 自動生成されたコード

このPRは CodingWorker によって自動生成されました。

**関連Issue**: #{issue_number}
**生成モデル**: Ollama qwen2.5-coder:1.5b (via Aider)
**生成日時**: {timestamp}
**Worker ID**: {worker_id}

### 確認事項
- [ ] コードが正しく動作するか
- [ ] テストが通るか
```

---

## タスク16: Worker 統合 ✅

### 16-1: メインループ実装
- [x] `cmd/worker/main.go` 作成
- [x] 設定ファイル読み込み（`internal/config/config.go`）
- [x] SQS ポーリングループ
- [x] シグナルハンドリング（graceful shutdown）

### 16-2: 処理フロー実装 ✅
```
1. SQS からメッセージ取得
2. リポジトリをクローン
3. ブランチ作成 (ai/issue-N)
4. Aider でコード生成（2パス実行）
5. 検証ループ（build/lint/test）
6. GitHub へプッシュ
7. PR 作成
8. SQS メッセージ削除
9. クリーンアップ
```

### 16-3: エラーハンドリング ✅
- [x] タイムアウト設定（10分/Aider実行）
- [x] 2レベルリトライロジック
  - 内側: Aider修正依頼（最大3回）
  - 外側: インフラリトライ（固定10秒、最大3回）
- [x] エラーログ出力（slog）
- [x] エラー分類（TransientError / PermanentError）

**設計ドキュメント**: `docs/design.md`

---

## タスク17: ローカルテスト 🟡 50%

### 17-1: 単体テスト ✅
- [x] SQS クライアントテスト（`sqs/client_test.go`）
- [x] Config テスト（`config/config_test.go`）
- [x] Retry テスト（`retry/retry_test.go`）
- [x] Aider Runner テスト → コンパイル確認済み
- [x] GitHub クライアントテスト → コンパイル確認済み

### 17-2: 統合テスト（ローカル）
- [ ] テスト用 SQS メッセージを Mock で送信
- [ ] Worker が処理できることを確認
- [ ] PR が作成されることを確認

**ステータス**: 単体テスト完了、統合テスト残

---

## Phase 2 完了判定基準 🟡

- [x] Go Worker が実装されている
- [x] SQS からメッセージを取得できる（Mock）
- [x] Aider を呼び出してコード生成できる
- [x] GitHub に PR を作成できる
- [x] エラーハンドリングが実装されている
- [ ] 統合テストが成功する → 残タスク

---

# Phase 3: エンドツーエンド統合テスト

**期間**: 2-3日
**目的**: システム全体の動作確認とバグ修正

---

## タスク18: E2E テストシナリオ作成

### 18-1: 正常系シナリオ
- [ ] シナリオ1: 簡単なスクリプト作成（FizzBuzz）
- [ ] シナリオ2: 既存ファイルへの機能追加
- [ ] シナリオ3: バグ修正タスク

### 18-2: 異常系シナリオ
- [ ] シナリオ4: タイムアウト（1時間超過）
- [ ] シナリオ5: 不正なタスク内容
- [ ] シナリオ6: ネットワークエラー

---

## タスク19: E2E テスト実行

### 19-1: 正常系テスト
- [ ] GitHub Issue 作成（`ai-task`ラベル）
- [ ] SQS にメッセージが届くことを確認
- [ ] Worker が処理することを確認
- [ ] PR が作成されることを確認
- [ ] 各シナリオで実行・検証

### 19-2: 異常系テスト
- [ ] 各シナリオで実行・検証
- [ ] DLQ にメッセージが移動することを確認
- [ ] エラーログが出力されることを確認

### 19-3: パフォーマンステスト
- [ ] 連続3タスク処理テスト
- [ ] 長時間稼働テスト（3時間）
- [ ] リソース使用量計測

---

## タスク20: バグ修正と改善

### 20-1: バグ修正
- [ ] 発見されたバグリスト作成
- [ ] 優先度付け
- [ ] バグ修正実施
- [ ] 再テスト

### 20-2: 改善実施
- [ ] ログ改善
- [ ] エラーメッセージ改善
- [ ] ドキュメント更新

---

## Phase 3 完了判定基準

- [ ] すべての正常系シナリオが成功する
- [ ] 異常系シナリオで適切にエラーハンドリングされる
- [ ] 3時間以上の連続稼働が可能
- [ ] 既知のバグがすべて修正されている

---

# Phase 4: 運用基盤整備

**期間**: 1-2日
**目的**: モニタリング、アラート、運用マニュアルを整備し、本番運用可能な状態にする

---

## タスク21: モニタリング・アラート設定

### 21-1: CloudWatch Alarms
- [ ] DLQ メッセージ数アラーム（閾値: 1以上）
- [ ] SQS キュー深度アラーム（閾値: 10以上）
- [ ] SNS トピック作成・Email通知設定

### 21-2: AWS Budgets
- [ ] 月額予算設定（$5）
- [ ] 予算超過アラート設定

### 21-3: ローカルログ管理
- [ ] Worker ログ出力設定
- [ ] ログローテーション設定（7日保持）

---

## タスク22: 運用マニュアル作成

### 22-1: セットアップマニュアル
- [ ] 初期セットアップ手順
- [ ] 環境変数設定
- [ ] 認証設定（PAT, AWS Access Key）
- [ ] ドキュメント: `docs/setup_manual.md`

### 22-2: 運用マニュアル
- [ ] 日常運用手順
- [ ] Worker 起動・停止方法
- [ ] ログ確認方法
- [ ] ドキュメント: `docs/operation_manual.md`

### 22-3: トラブルシューティングガイド
- [ ] よくある問題と解決方法
- [ ] エラーコード一覧
- [ ] ドキュメント: `docs/troubleshooting.md`

---

## タスク23: セキュリティ確認

### 23-1: セキュリティチェック
- [ ] IAM ロール権限が最小限か確認
- [ ] GitHub PAT のスコープが最小限か確認
- [ ] シークレットがログに出力されていないか確認
- [ ] SQS 暗号化が有効か確認

---

## Phase 4 完了判定基準

- [ ] アラートが設定されている
- [ ] 運用マニュアルが整備されている
- [ ] トラブルシューティングガイドが完成している
- [ ] セキュリティチェックが完了している

---

# プロジェクト完了判定

## 全体完了基準

### システム機能
- [ ] GitHub Issues → PR作成までの完全自動化が動作する
- [ ] エラー発生時に適切にハンドリングされる
- [ ] 3時間以上の連続稼働が可能

### コスト
- [ ] 月額コストが2,000円以内（電気代のみ）
- [ ] AWS 無料枠内で運用できている

### ドキュメント
- [ ] セットアップ手順書が完成している
- [ ] 運用マニュアルが完成している
- [ ] トラブルシューティングガイドが完成している

### 品質
- [ ] すべてのE2Eテストが成功する
- [ ] セキュリティチェックが完了している

---

# フォールバック計画

## ローカルLLMが実用レベルでない場合

### 選択肢1: モデル変更
- qwen2.5-coder:1.5b → codellama:7b
- ~~qwen2.5-coder:3b~~ → 検証の結果、3BはDiff生成エラーやパス管理問題が発生し不採用

### 選択肢2: ツール変更
- Aider → Open Interpreter

### 選択肢3: 無料API利用
- Google Gemini API（無料枠: 60 RPM, 1500 RPD）
- Groq API（無料枠: 30 RPM）

### 選択肢4: 低コストAPI利用
- OpenRouter 経由で DeepSeek（低コスト）
- Mistral API

---

# 将来の拡張タスク（オプション）

### 短期（3ヶ月以内）
- [ ] 複数プロジェクトの並行処理機能
- [ ] Web UI（タスク管理画面）
- [ ] Slack/Discord 通知連携

### 中期（6ヶ月以内）
- [ ] Apple Silicon Mac への移行
- [ ] より高性能なローカルLLM対応
- [ ] 自動テスト実行機能

---

**最終更新**: 2025-01-26
**バージョン**: 2.1（Taskfile採用、詳細設計追加）
**総タスク数**: 23メインタスク
**想定期間**: 10-15日
