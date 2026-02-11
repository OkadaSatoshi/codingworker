# AWS アカウント初期セットアップガイド

AWS Organizations + IAM Identity Center (SSO) によるマルチアカウント構成のセットアップ手順。

> **詳細な手順・Q&A**: [Issue #1: AWS環境準備](https://github.com/OkadaSatoshi/codingworker/issues/1) を参照

---

## 1. 構成概要

```
Management Account (個人アカウント)
└── codingworker-dev  ← 開発環境
```

| 項目 | 方式 |
|:---|:---|
| アカウント管理 | AWS Organizations (All features) |
| 認証方式 | IAM Identity Center (SSO) |
| CLI 認証 | `aws configure sso` → SSO プロファイル |
| 日常操作 | SSO ユーザー (AdministratorAccess Permission Set) |

---

## 2. セットアップ手順概要

### 2.1 Organizations の有効化 (Management Account)

1. AWS Console → AWS Organizations → Create an organization
2. 機能: All features

### 2.2 メンバーアカウントの作成

1. Organizations → Add an AWS account → Create
2. Account name: `codingworker-dev`
3. IAM role name: `OrganizationAccountAccessRole`（デフォルト）
4. Tags: `Environment=dev`, `Project=codingworker`

### 2.3 IAM Identity Center の設定 (Management Account)

1. IAM Identity Center を有効化（リージョン: ap-northeast-1）
2. ユーザー作成 + MFA 設定
3. Permission Set 作成: `AdministratorAccess`
4. `codingworker-dev` アカウントにユーザーを割り当て

### 2.4 メンバーアカウントのセキュリティ

1. ルートユーザーのパスワードリセット + MFA 有効化（認証アプリ推奨）
2. admin IAM ユーザー作成 + MFA 有効化（緊急用）

### 2.5 AWS CLI の SSO 設定

```zsh
aws configure sso
# SSO session name: codingworker
# SSO start URL: https://d-xxxxxxxxxx.awsapps.com/start
# SSO region: ap-northeast-1
# → プロファイル名: codingworker-dev
```

### 2.6 動作確認

```zsh
aws sso login --profile codingworker-dev
aws sts get-caller-identity --profile codingworker-dev
```

### 2.7 請求アラート設定 (Management Account)

- Billing → Budgets で無料枠監視アラートを設定

---

## 3. Terraform 実行準備

### 3.1 Terraform インストール

```zsh
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

※ Homebrew 標準の `terraform` は 1.5.7 で更新停止。hashicorp/tap から最新版をインストール。

### 3.2 State 管理

S3 backend + DynamoDB lock を使用。

| リソース | 名前 | 作成方法 |
|:---|:---|:---|
| S3 バケット | `codingworker-dev-tfstate` | CLI 手動作成 |
| DynamoDB テーブル | `codingworker-dev-tfstate-lock` | CLI 手動作成 |

作成手順は [Issue #23](https://github.com/OkadaSatoshi/codingworker/issues/23) を参照。

### 3.3 初回実行

```zsh
cd infra/terraform

# 変数ファイル作成
cp terraform.tfvars.example terraform.tfvars
# terraform.tfvars を編集（github_org, github_repos）

# 初期化
terraform init

# 確認・適用
terraform plan
terraform apply
```

---

## 4. Worker 用 Access Key

Worker (MBP) は 24h 稼働のため、SSO セッション（有効期限あり）ではなく IAM User の Access Key を使用する。

```zsh
aws iam create-access-key --user-name codingworker-worker --profile codingworker-dev
```

Access Key は Terraform 管理外（tfstate にシークレットを残さないため）。

---

## 5. 作業チェックリスト

- [x] Organizations 有効化
- [x] codingworker-dev アカウント作成
- [x] IAM Identity Center 設定
- [x] SSO ユーザー作成 + MFA
- [x] ルートユーザー MFA 有効化
- [x] AWS CLI SSO 設定・動作確認
- [x] 請求アラート設定
- [x] Terraform init/plan/apply 実行
- [ ] Worker 用 Access Key 作成

---

**作成日**: 2025-01-26
**更新日**: 2025-02-08 - Organizations/SSO方式に全面改訂
