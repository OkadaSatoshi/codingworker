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
│  Aider + Ollama │ ← ローカルLLM (qwen2.5-coder:7b)
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
| Phase 0: ローカル環境 PoC | 🟡 進行中 | 0% | 2-3日 |
| Phase 1: Cloud Foundation | ⚪ 未着手 | 0% | 2-3日 |
| Phase 2: SQS連携・ワーカー | ⚪ 未着手 | 0% | 3-4日 |
| Phase 3: E2E統合テスト | ⚪ 未着手 | 0% | 2-3日 |
| Phase 4: 運用基盤整備 | ⚪ 未着手 | 0% | 1-2日 |

---

# Phase 0: ローカル環境 PoC 構築

**期間**: 2-3日
**目的**: MBP 2018 上で Aider + Ollama を動作させ、基本的なコーディングタスクが実行可能であることを検証する

---

## タスク0: 環境確認

### 0-1: MBP環境の確認（30分）
- [ ] macOSバージョン確認（Monterey 12.0以上推奨）
- [ ] 空きディスク容量確認（最低10GB必要）
- [ ] 空きメモリ確認（16GB中、8GB以上空き推奨）
- [ ] Python 3.8+ インストール確認（`python3 --version`）
- [ ] Git インストール確認（`git --version`）
- [ ] Homebrew インストール確認（`brew --version`）

### 0-2: 不足ツールのインストール
- [ ] Python未インストールの場合: `brew install python`
- [ ] Git未インストールの場合: `brew install git`
- [ ] pip確認: `pip3 --version`

---

## タスク1: Ollama セットアップ

### 1-1: Ollama のインストール
- [ ] 公式サイト（https://ollama.ai）からインストーラーをダウンロード
- [ ] インストール実行
- [ ] `ollama --version` で動作確認
- [ ] Ollama サービス起動確認

### 1-2: qwen2.5-coder:7b モデルのダウンロード
- [ ] `ollama pull qwen2.5-coder:7b` 実行
- [ ] ダウンロード完了確認（約4GB）
- [ ] `ollama list` でモデル確認

### 1-3: Ollama 基本動作確認
- [ ] `ollama run qwen2.5-coder:7b` で対話テスト
- [ ] 簡単なコード生成プロンプト実行（例: "Write a hello world in Python"）
- [ ] レスポンス速度を目視確認（目安: 30秒以内）
- [ ] `Ctrl+D` で終了

---

## タスク2: Aider セットアップ

### 2-1: Aider のインストール
- [ ] `pip3 install aider-chat` 実行
- [ ] `aider --version` で確認
- [ ] インストール完了確認

### 2-2: Aider + Ollama 連携設定
- [ ] テスト用ディレクトリ作成: `mkdir ~/aider-test && cd ~/aider-test`
- [ ] Git初期化: `git init`
- [ ] Aider起動: `aider --model ollama_chat/qwen2.5-coder:7b`
- [ ] 接続成功確認

### 2-3: Aider 基本動作確認
- [ ] 簡単なファイル作成指示（例: "Create a hello.py that prints hello world"）
- [ ] ファイルが生成されることを確認
- [ ] 自動コミットされることを確認（`git log`）
- [ ] `/quit` で終了

---

## タスク3: PoC テスト実行

### 3-1: テストケース1 - FizzBuzz
- [ ] 新規ディレクトリ作成・Git初期化
- [ ] Aider起動
- [ ] タスク投入: "Create a Python script fizzbuzz.py that prints FizzBuzz from 1 to 100"
- [ ] 生成コードを確認
- [ ] `python3 fizzbuzz.py` で実行・動作確認
- [ ] 実行時間を記録

### 3-2: テストケース2 - ファイル処理スクリプト
- [ ] タスク投入: "Create a Python script that reads a CSV file and sorts it by the first column"
- [ ] 生成コードを確認
- [ ] テスト用CSVで動作確認
- [ ] 実行時間を記録

### 3-3: テストケース3 - バグ修正
- [ ] バグを含むコードを準備（例: off-by-one error）
- [ ] タスク投入: "Fix the bug in buggy.py"
- [ ] バグ検出・修正能力を確認
- [ ] 実行時間を記録

### 3-4: テスト結果まとめ
- [ ] 各テストケースの成否記録
- [ ] 平均実行時間算出
- [ ] 生成コード品質評価（構文エラー有無、正確性）

---

## タスク4: パフォーマンス計測

### 4-1: リソース使用量の監視
- [ ] Activity Monitor でOllamaのメモリ使用量確認
- [ ] Aider実行中のメモリ使用量確認
- [ ] 合計メモリ使用量が14GB以下か確認
- [ ] CPU使用率記録

### 4-2: 推論速度の計測
- [ ] 3つのテストケースそれぞれで時間計測
- [ ] 平均推論速度算出
- [ ] 目標: 60秒/応答以内

### 4-3: 長時間稼働テスト（オプション）
- [ ] 1時間連続で複数タスク実行
- [ ] メモリリークの有無確認
- [ ] 安定性確認

---

## タスク5: Go/No-Go 判断

### 5-1: 評価基準チェック

#### Go 基準（すべて満たす場合、Phase 1へ）
- [ ] 推論速度: 60秒/応答以内
- [ ] メモリ使用量: 14GB以下
- [ ] テストケース: 3つ中2つ以上成功
- [ ] Aider + Ollama連携: 安定動作

#### Pivot 基準（代替案検討）
- [ ] Aider連携が不安定 → Open Interpreter検証
- [ ] qwen2.5-coder:7b品質不足 → codellama:7b検証

#### Fallback 基準（無料API検討）
- [ ] ローカルLLM推論が120秒/応答超過
- [ ] 品質が実用レベルでない
- [ ] → Google Gemini API（無料枠）を検討

#### No-Go 基準（プロジェクト中止）
- [ ] すべての選択肢で実用レベルに達しない

### 5-2: 判断結果記録
- [ ] 判断結果: Go / Pivot / Fallback / No-Go
- [ ] 理由記録
- [ ] 次のアクション決定

---

## タスク6: ドキュメント化

### 6-1: セットアップ手順書
- [ ] Ollama インストール手順
- [ ] Aider インストール手順
- [ ] 連携設定手順
- [ ] ドキュメント: `docs/setup_guide.md`

### 6-2: PoC結果レポート
- [ ] テスト結果まとめ
- [ ] パフォーマンスデータ
- [ ] 課題・改善点
- [ ] ドキュメント: `docs/poc_report.md`

### 6-3: Phase 1 引き継ぎ事項
- [ ] 明らかになった課題
- [ ] 推奨設定値
- [ ] ドキュメント: `docs/phase1_handover.md`

---

## Phase 0 完了判定基準

- [ ] Ollama + qwen2.5-coder:7b が安定動作する
- [ ] Aider が Ollama 経由でコードを生成できる
- [ ] 生成されたコードが実行可能である
- [ ] Go/No-Go判断が完了している
- [ ] セットアップ手順が文書化されている

---

# Phase 1: Cloud Foundation 構築

**期間**: 2-3日
**目的**: AWS上にクラウド基盤を構築し、GitHub Actions と SQS の連携を実現する
**前提**: Phase 0 で Go 判定

---

## タスク7: AWS環境準備

### 7-1: AWSアカウント確認
- [ ] AWSアカウント有無確認
- [ ] 無料枠残量確認
- [ ] 請求アラート設定（$5超過で通知）

### 7-2: AWS CLI セットアップ
- [ ] `brew install awscli` 実行
- [ ] `aws --version` で確認
- [ ] IAM User作成（MBP用、プログラムアクセス）
- [ ] Access Key / Secret Key 取得
- [ ] `aws configure` で設定
- [ ] `aws sts get-caller-identity` で確認

---

## タスク8: Terraform セットアップ

### 8-1: Terraform インストール
- [ ] `brew install terraform` 実行
- [ ] `terraform --version` で確認

### 8-2: Terraform State 管理準備
- [ ] S3バケット作成（State保存用）: `codingworker-tfstate-{account-id}`
- [ ] DynamoDB テーブル作成（State Lock用）: `codingworker-tfstate-lock`
- [ ] `backend.tf` 設定ファイル作成

### 8-3: codingworker-infra リポジトリ作成
- [ ] GitHub に新規リポジトリ作成
- [ ] ディレクトリ構造作成
```
codingworker-infra/
├── terraform/
│   ├── backend.tf
│   ├── provider.tf
│   ├── variables.tf
│   ├── outputs.tf
│   ├── oidc.tf
│   ├── sqs.tf
│   └── iam.tf
└── README.md
```
- [ ] README.md 作成

---

## タスク9: IAM OIDC プロバイダー設定

### 9-1: IAM OIDC プロバイダー作成
- [ ] Terraform コード作成: `oidc.tf`
- [ ] GitHub OIDC プロバイダー設定
- [ ] Trust Policy 設定

### 9-2: IAM ロール作成（GitHub Actions用）
- [ ] SQS SendMessage 権限を持つロール作成
- [ ] Terraform コード: `iam.tf`
- [ ] 最小権限原則の確認

### 9-3: IAM User作成（MBP Worker用）
- [ ] SQS ReceiveMessage, DeleteMessage 権限
- [ ] Terraform コード追加
- [ ] Access Key 出力設定

---

## タスク10: SQS キュー作成

### 10-1: SQS Standard Queue 作成
- [ ] Terraform コード作成: `sqs.tf`
- [ ] キュー名: `codingworker-tasks`
- [ ] 可視性タイムアウト: 3600秒
- [ ] 暗号化設定（AWS管理キー）

### 10-2: Dead Letter Queue 作成
- [ ] DLQ作成: `codingworker-tasks-dlq`
- [ ] メインキューに DLQ 設定
- [ ] 最大受信回数: 3回

### 10-3: Terraform Apply
- [ ] `terraform init` 実行
- [ ] `terraform plan` で確認
- [ ] `terraform apply` で適用
- [ ] 作成されたリソース確認

---

## タスク11: GitHub Actions ワークフロー作成

### 11-1: テスト用リポジトリ準備
- [ ] `codingworker-test-project` リポジトリ作成
- [ ] `.github/workflows/` ディレクトリ作成

### 11-2: OIDC 認証ワークフロー作成
- [ ] `.github/workflows/send-to-sqs.yml` 作成
- [ ] OIDC 認証設定
- [ ] Issue イベントトリガー（`auto-code`ラベル）

### 11-3: SQS メッセージ送信実装
- [ ] メッセージフォーマット（JSON）作成
- [ ] AWS CLI で SendMessage 実装
- [ ] エラーハンドリング

### 11-4: 連携テスト
- [ ] テスト用 Issue 作成（`auto-code`ラベル付与）
- [ ] ワークフロー実行確認
- [ ] SQS にメッセージが届くか確認（AWS Console）

---

## Phase 1 完了判定基準

- [ ] Terraform で AWS リソースが作成されている
- [ ] IAM OIDC が正しく設定されている
- [ ] SQS キュー（メイン + DLQ）が作成されている
- [ ] GitHub Actions → SQS へのメッセージ送信が成功する
- [ ] Terraform コードがリポジトリに保存されている

---

# Phase 2: SQS連携とワーカー開発

**期間**: 3-4日
**目的**: Go製ワーカーを実装し、SQS → Aider → GitHub PR の連携を実現する

---

## タスク12: codingworker-host リポジトリ作成

### 12-1: リポジトリ作成
- [ ] GitHub に新規リポジトリ作成
- [ ] ディレクトリ構造作成
```
codingworker-host/
├── cmd/
│   └── worker/
│       └── main.go
├── internal/
│   ├── sqs/
│   │   └── client.go
│   ├── aider/
│   │   └── runner.go
│   └── github/
│       └── pr.go
├── configs/
│   └── config.yaml
├── scripts/
│   └── start.sh
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 12-2: Go モジュール初期化
- [ ] `go mod init github.com/{user}/codingworker-host`
- [ ] 依存関係追加
  - `github.com/aws/aws-sdk-go-v2`
  - `github.com/google/go-github/v50`
  - `gopkg.in/yaml.v3`
- [ ] Makefile 作成

---

## タスク13: SQS ポーリング機能実装

### 13-1: SQS クライアント実装
- [ ] `internal/sqs/client.go` 作成
- [ ] AWS SDK v2 でクライアント作成
- [ ] 認証設定（環境変数 or 設定ファイル）

### 13-2: メッセージ受信処理
- [ ] ロングポーリング実装（WaitTimeSeconds: 20）
- [ ] メッセージパース（JSON → struct）
- [ ] エラーハンドリング

### 13-3: メッセージ削除処理
- [ ] 処理完了後の DeleteMessage 実装
- [ ] 失敗時の処理（メッセージを削除せず再キューイング）

---

## タスク14: Aider 連携実装

### 14-1: Aider Runner 実装
- [ ] `internal/aider/runner.go` 作成
- [ ] Aider CLI 呼び出し処理
- [ ] コマンド構築: `aider --model ollama_chat/qwen2.5-coder:7b --yes --message "{task}"`

### 14-2: 作業ディレクトリ管理
- [ ] リポジトリのクローン処理
- [ ] ブランチ作成: `auto-code/issue-{number}`
- [ ] 作業完了後のクリーンアップ

### 14-3: 実行結果取得
- [ ] Aider 終了コード確認
- [ ] 生成ファイルの確認
- [ ] エラー検出

---

## タスク15: GitHub 連携実装

### 15-1: GitHub クライアント実装
- [ ] `internal/github/pr.go` 作成
- [ ] go-github ライブラリ使用
- [ ] Personal Access Token 認証

### 15-2: PR 作成処理
- [ ] ブランチのプッシュ
- [ ] Pull Request 作成
- [ ] PR 本文テンプレート適用
- [ ] Issue番号参照（`Closes #123`）

### 15-3: PR テンプレート
```markdown
## 自動生成されたコード

このPRは CodingWorker によって自動生成されました。

**関連Issue**: #{issue_number}
**生成モデル**: Ollama qwen2.5-coder:7b (via Aider)
**生成日時**: {timestamp}

### 確認事項
- [ ] コードが正しく動作するか
- [ ] テストが通るか
```

---

## タスク16: Worker 統合

### 16-1: メインループ実装
- [ ] `cmd/worker/main.go` 作成
- [ ] 設定ファイル読み込み
- [ ] SQS ポーリングループ
- [ ] シグナルハンドリング（graceful shutdown）

### 16-2: 処理フロー実装
```
1. SQS からメッセージ取得
2. リポジトリをクローン
3. ブランチ作成
4. Aider でコード生成
5. 変更をコミット（Aiderが自動実行）
6. GitHub へプッシュ
7. PR 作成
8. SQS メッセージ削除
9. クリーンアップ
```

### 16-3: エラーハンドリング
- [ ] タイムアウト設定（1時間）
- [ ] リトライロジック
- [ ] エラーログ出力

---

## タスク17: ローカルテスト

### 17-1: 単体テスト
- [ ] SQS クライアントテスト
- [ ] Aider Runner テスト
- [ ] GitHub クライアントテスト

### 17-2: 統合テスト（ローカル）
- [ ] テスト用 SQS メッセージを手動送信
- [ ] Worker が処理できることを確認
- [ ] PR が作成されることを確認

---

## Phase 2 完了判定基準

- [ ] Go Worker が実装されている
- [ ] SQS からメッセージを取得できる
- [ ] Aider を呼び出してコード生成できる
- [ ] GitHub に PR を作成できる
- [ ] エラーハンドリングが実装されている

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
- [ ] GitHub Issue 作成（`auto-code`ラベル）
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
- qwen2.5-coder:7b → codellama:7b
- qwen2.5-coder:7b → qwen2.5-coder:3b（軽量版）

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

**最終更新**: 2025-01-25
**バージョン**: 2.0（Aider版）
**総タスク数**: 23メインタスク
**想定期間**: 10-15日
