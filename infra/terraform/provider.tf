terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region  = var.aws_region
  profile = "codingworker-dev"

  default_tags {
    tags = {
      Project     = "codingworker"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}
