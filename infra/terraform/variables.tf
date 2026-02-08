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

  validation {
    condition     = length(var.github_org) > 0
    error_message = "github_org must not be empty."
  }
}

variable "github_repos" {
  description = "GitHub repository names that can trigger workflows"
  type        = list(string)

  validation {
    condition     = length(var.github_repos) > 0
    error_message = "github_repos must contain at least one repository."
  }
}

variable "sqs_visibility_timeout" {
  description = "SQS visibility timeout in seconds"
  type        = number
  default     = 10800 # 3 hours (Aider processing on low-spec machine)

  validation {
    condition     = var.sqs_visibility_timeout >= 0 && var.sqs_visibility_timeout <= 43200
    error_message = "sqs_visibility_timeout must be between 0 and 43200 (12 hours)."
  }
}

variable "sqs_message_retention" {
  description = "SQS message retention in seconds"
  type        = number
  default     = 604800 # 7 days

  validation {
    condition     = var.sqs_message_retention >= 60 && var.sqs_message_retention <= 1209600
    error_message = "sqs_message_retention must be between 60 (1 minute) and 1209600 (14 days)."
  }
}

variable "sqs_max_receive_count" {
  description = "Max receive count before moving to DLQ"
  type        = number
  default     = 3

  validation {
    condition     = var.sqs_max_receive_count >= 1
    error_message = "sqs_max_receive_count must be at least 1."
  }
}
