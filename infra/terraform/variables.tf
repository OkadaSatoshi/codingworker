variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "github_org" {
  description = "GitHub organization or username"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name for the test project"
  type        = string
  default     = "codingworker-sandbox"
}

variable "sqs_visibility_timeout" {
  description = "SQS visibility timeout in seconds"
  type        = number
  default     = 3600 # 1 hour (Aider processing time)
}

variable "sqs_message_retention" {
  description = "SQS message retention in seconds"
  type        = number
  default     = 86400 # 1 day
}

variable "dlq_max_receive_count" {
  description = "Max receive count before moving to DLQ"
  type        = number
  default     = 3
}
