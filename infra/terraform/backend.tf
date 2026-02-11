terraform {
  backend "s3" {
    bucket         = "codingworker-dev-tfstate"
    key            = "codingworker/terraform.tfstate"
    region         = "ap-northeast-1"
    encrypt        = true
    dynamodb_table = "codingworker-dev-tfstate-lock"
    profile        = "codingworker-dev"
  }
}
