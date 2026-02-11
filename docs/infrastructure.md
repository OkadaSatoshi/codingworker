# CodingWorker インフラストラクチャ設計

## 1. 概要

本ドキュメントはCodingWorkerのAWSインフラストラクチャとGitHub Actions連携の設計を記述する。

---

## 2. システム構成図

```
┌─────────────────────────────────────────────────────────────────┐
│                         GitHub                                   │
│  ┌─────────────┐     ┌──────────────────┐                       │
│  │   Issues    │────▶│  GitHub Actions  │                       │
│  │ (ai-task)   │     │  send-to-sqs.yml │                       │
│  └─────────────┘     └────────┬─────────┘                       │
└───────────────────────────────┼─────────────────────────────────┘
                                │ OIDC認証
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                           AWS                                    │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    IAM OIDC Provider                      │   │
│  │              (token.actions.githubusercontent.com)        │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                │                                 │
│                                ▼                                 │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              IAM Role: github-actions                     │   │
│  │              (sqs:SendMessage)                            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                │                                 │
│                                ▼                                 │
│  ┌─────────────────────┐     ┌─────────────────────┐            │
│  │  SQS: tasks         │────▶│  SQS: tasks-dlq     │            │
│  │  (メインキュー)      │     │  (デッドレターキュー) │            │
│  └──────────┬──────────┘     └─────────────────────┘            │
└─────────────┼───────────────────────────────────────────────────┘
              │ ロングポーリング
              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      MBP (ローカル)                              │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              IAM User: worker                             │   │
│  │              (sqs:ReceiveMessage, sqs:DeleteMessage)      │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                │                                 │
│                                ▼                                 │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Go Worker                              │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. AWS リソース

### 3.1 SQS キュー

#### メインキュー: `codingworker-tasks`

| 設定項目 | 値 | 説明 |
|:---|:---|:---|
| visibility_timeout | 10800秒 (3時間) | 低スペックマシンでのAider処理時間を考慮 |
| message_retention | 604800秒 (7日) | 未処理メッセージの保持期間 |
| receive_wait_time | 20秒 | ロングポーリング |
| redrive_policy | maxReceiveCount: 3 | 3回失敗でDLQへ移動 |
| sqs_managed_sse_enabled | true | サーバーサイド暗号化 (SSE-SQS) |
| Queue Policy | DenyNonSSLAccess | HTTPS以外のアクセスを拒否 |

#### デッドレターキュー: `codingworker-tasks-dlq`

| 設定項目 | 値 | 説明 |
|:---|:---|:---|
| message_retention | 1209600秒 (14日) | SQS最大値。調査用に長めに保持 |
| sqs_managed_sse_enabled | true | サーバーサイド暗号化 (SSE-SQS) |
| redrive_allow_policy | byQueue | メインキューのみリドライブ元として許可 |
| Queue Policy | DenyNonSSLAccess | HTTPS以外のアクセスを拒否 |

### 3.2 IAM OIDC Provider

GitHub Actionsからのアクセスを許可するためのOIDCプロバイダー。

| 設定項目 | 値 |
|:---|:---|
| URL | `https://token.actions.githubusercontent.com` |
| Audience | `sts.amazonaws.com` |

### 3.3 IAM Role: `codingworker-github-actions`

GitHub Actionsが使用するロール。

**信頼ポリシー:**
```json
{
  "Effect": "Allow",
  "Principal": {
    "Federated": "arn:aws:iam::{account}:oidc-provider/token.actions.githubusercontent.com"
  },
  "Action": "sts:AssumeRoleWithWebIdentity",
  "Condition": {
    "StringEquals": {
      "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
    },
    "StringEquals": {
      "token.actions.githubusercontent.com:sub": "repo:{org}/{repo}:ref:refs/heads/main"
    }
  }
}
```

**権限:**
- `sqs:SendMessage`
- `sqs:GetQueueUrl`
- `sqs:GetQueueAttributes`

### 3.4 IAM User: `codingworker-worker`

MBP上のWorkerが使用するユーザー。

**権限:**
- `sqs:ReceiveMessage`
- `sqs:DeleteMessage`
- `sqs:GetQueueUrl`
- `sqs:GetQueueAttributes`
- `sqs:ChangeMessageVisibility`

---

## 4. GitHub Actions ワークフロー

### 4.1 トリガー条件

| イベント | 条件 |
|:---|:---|
| `issues.labeled` | ラベル名が `ai-task` |
| 有効化フラグ | `vars.ENABLE_SQS_WORKFLOW == 'true'` |

### 4.2 必要なシークレット/変数

| 名前 | 種類 | 説明 |
|:---|:---|:---|
| `AWS_ROLE_ARN` | Secret | GitHub Actions用IAMロールのARN |
| `ENABLE_SQS_WORKFLOW` | Variable | ワークフロー有効化フラグ |

### 4.3 処理フロー

```
1. Issue に "ai-task" ラベル付与
2. GitHub Actions 起動
3. OIDC で AWS 認証
4. SQS にメッセージ送信
5. Issue にコメント通知
```

---

## 5. SQS メッセージフォーマット

### 5.1 メッセージ構造

```json
{
  "issue_number": 123,
  "repository": "OkadaSatoshi/codingworker",
  "title": "Add FizzBuzz function",
  "body": "Create a function that prints FizzBuzz from 1 to 100",
  "labels": ["ai-task"],
  "created_at": "2025-01-26T00:00:00Z"
}
```

### 5.2 フィールド説明

| フィールド | 型 | 必須 | 説明 |
|:---|:---|:---|:---|
| issue_number | number | ✅ | GitHub Issue番号 |
| repository | string | ✅ | リポジトリ (owner/repo形式) |
| title | string | ✅ | Issueタイトル |
| body | string | ✅ | Issue本文 (タスク内容) |
| labels | string[] | ✅ | Issueのラベル一覧 |
| created_at | string | ✅ | メッセージ作成日時 (ISO 8601) |

---

## 6. Terraform 構成

### 6.1 ファイル構成

```
infra/terraform/
├── backend.tf      # State管理 (S3)
├── provider.tf     # AWSプロバイダー
├── variables.tf    # 入力変数
├── outputs.tf      # 出力値
├── oidc.tf         # GitHub OIDC Provider
├── sqs.tf          # SQSキュー
└── iam.tf          # IAM Role/User
```

### 6.2 入力変数

| 変数名 | デフォルト値 | 説明 |
|:---|:---|:---|
| aws_region | ap-northeast-1 | AWSリージョン |
| environment | dev | 環境名 |
| github_org | - | GitHubオーナー名 (必須) |
| github_repos | - | 対象リポジトリ名リスト (必須) |
| sqs_visibility_timeout | 10800 | 可視性タイムアウト (秒) |
| sqs_message_retention | 604800 | メッセージ保持期間 (秒) |
| sqs_max_receive_count | 3 | DLQ移動までの最大受信回数 |

### 6.3 出力値

| 出力名 | 説明 |
|:---|:---|
| sqs_queue_url | SQSキューURL |
| sqs_queue_arn | SQSキューARN |
| sqs_dlq_url | DLQキューURL |
| github_actions_role_arn | GitHub Actions用ロールARN |
| oidc_provider_arn | GitHub OIDC ProviderのARN |

### 6.4 適用手順

```bash
cd infra/terraform

# 変数ファイル作成
cp terraform.tfvars.example terraform.tfvars
# terraform.tfvars を編集

# 初期化
terraform init

# 確認
terraform plan

# 適用
terraform apply
```

---

## 7. セキュリティ

### 7.1 最小権限の原則

| リソース | 許可されたアクション |
|:---|:---|
| GitHub Actions Role | SQS送信のみ |
| Worker User | SQS受信・削除のみ |

### 7.2 認証情報の管理

| 認証情報 | 保存場所 | 備考 |
|:---|:---|:---|
| AWS_ROLE_ARN | GitHub Secrets | GitHub Actions OIDC認証用 |
| Worker Access Key | ~/.aws/credentials | CLI手動作成（Terraform管理外） |
| GITHUB_TOKEN | ローカル環境変数 | Worker の GitHub 操作用 |

---

### 7.3 Terraform Backend

| 項目 | 値 |
|:---|:---|
| Backend | S3 |
| Bucket | `codingworker-dev-tfstate` |
| Key | `codingworker/terraform.tfstate` |
| DynamoDB Lock Table | `codingworker-dev-tfstate-lock` |
| 暗号化 | AES256 (SSE-S3) |
| バケットポリシー | DenyNonSSLAccess |

S3バケット・DynamoDBテーブルはTerraform管理外（CLI手動作成）。

---

## 8. 更新履歴

| 日付 | 内容 |
|:---|:---|
| 2025-01-26 | 初版作成 |
| 2025-02-08 | SQS設定値更新、暗号化・Queue Policy追加、OIDC条件修正、S3 backend移行、Access Key管理方針変更 |
