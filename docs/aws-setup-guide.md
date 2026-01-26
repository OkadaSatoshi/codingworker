# AWS アカウント初期セットアップガイド

AWSルートアカウント作成後に行うべき設定手順。

---

## 1. ルートアカウントのセキュリティ強化（最優先）

ルートアカウントは日常的に使用せず、以下を設定後は IAM ユーザーを使用する。

### 1.1 MFA（多要素認証）の有効化

1. AWS Console にルートアカウントでログイン
2. 右上のアカウント名 → Security credentials
3. MFA デバイスを割り当て（Google Authenticator 等のアプリ推奨）

### 1.2 ルートアカウントのアクセスキー

**ルートアカウントのアクセスキーは絶対に作成しない。**

---

## 2. 管理者用 IAM ユーザーの作成

Terraform 実行や日常の AWS 操作に使用する。

### 2.1 Console で手動作成

1. IAM → Users → Create user
2. ユーザー名: `admin`（任意）
3. 「Provide user access to the AWS Management Console」にチェック
4. Permissions: `AdministratorAccess` ポリシーをアタッチ
5. ユーザー作成完了

### 2.2 MFA の有効化

1. 作成したユーザーの Security credentials タブ
2. MFA device → Assign MFA device

### 2.3 アクセスキーの作成（CLI 用）

1. Security credentials → Access keys → Create access key
2. Use case: Command Line Interface (CLI)
3. アクセスキー ID とシークレットアクセスキーを安全に保存

---

## 3. 請求アラートの設定

予期しない課金を防ぐ。

1. Billing → Budgets → Create budget
2. Budget type: Cost budget
3. 月額予算: $5〜$10（任意）
4. Alert threshold: 80% で通知
5. 通知先メールアドレスを設定

---

## 4. AWS CLI の設定

```bash
# インストール（未インストールの場合）
brew install awscli

# 認証情報の設定
aws configure
# Access Key ID: （管理者 IAM ユーザーのもの）
# Secret Access Key: （管理者 IAM ユーザーのもの）
# Region: ap-northeast-1（東京）
# Output format: json

# 確認
aws sts get-caller-identity
```

---

## 5. IAM ユーザーの設計指針

### 人間用 vs サービス用

| 種類 | 数 | 用途 |
|------|---|------|
| 管理者（人間用） | 1 | Terraform 実行、Console 操作 |
| サービス用 | PJ ごと | アプリケーションが使用 |

### 命名規則

```
{プロジェクト名}-{用途}

例:
- codingworker-github-actions（GitHub Actions 用ロール）
- codingworker-worker（ローカル Worker 用ユーザー）
```

### セキュリティ原則

- 最小権限: 必要な権限のみ付与
- PJ 分離: プロジェクトごとに IAM を分ける（1 つ漏洩しても他に影響なし）

---

## 6. Terraform 実行準備

### 6.1 State 管理

開発初期はローカル State で十分。S3 + DynamoDB は後から移行可能。

### 6.2 初回実行

```bash
cd infra/terraform

# 変数ファイル作成
cp terraform.tfvars.example terraform.tfvars
# terraform.tfvars を編集（github_org 等）

# 初期化
terraform init

# 確認
terraform plan

# 適用
terraform apply
```

---

## 7. 作業チェックリスト

- [ ] ルートアカウント MFA 有効化
- [ ] 管理者 IAM ユーザー作成
- [ ] 管理者 IAM ユーザー MFA 有効化
- [ ] 管理者でログインし直す（ルートは使わない）
- [ ] 請求アラート設定
- [ ] AWS CLI 設定・動作確認
- [ ] Terraform init/plan 実行

---

**作成日**: 2025-01-26
