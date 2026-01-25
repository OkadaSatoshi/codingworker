# Terraform Backend Configuration
# 開発環境: ローカルバックエンド
# 本番環境: S3バックエンドに切り替え

terraform {
  backend "local" {
    path = "terraform.tfstate"
  }

  # TODO: AWS環境準備後、以下に切り替え
  # backend "s3" {
  #   bucket         = "codingworker-tfstate-${ACCOUNT_ID}"
  #   key            = "codingworker/terraform.tfstate"
  #   region         = "ap-northeast-1"
  #   encrypt        = true
  #   dynamodb_table = "codingworker-tfstate-lock"
  # }
}
