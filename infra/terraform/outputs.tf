# =============================================================================
# SQS Outputs
# =============================================================================

output "sqs_queue_url" {
  description = "URL of the main SQS queue"
  value       = aws_sqs_queue.tasks.url
}

output "sqs_queue_arn" {
  description = "ARN of the main SQS queue"
  value       = aws_sqs_queue.tasks.arn
}

output "sqs_dlq_url" {
  description = "URL of the dead letter queue"
  value       = aws_sqs_queue.tasks_dlq.url
}

# =============================================================================
# IAM Outputs
# =============================================================================

output "github_actions_role_arn" {
  description = "ARN of the IAM role for GitHub Actions"
  value       = aws_iam_role.github_actions.arn
}

output "worker_access_key_id" {
  description = "Access key ID for the worker"
  value       = aws_iam_access_key.worker.id
}

output "worker_secret_access_key" {
  description = "Secret access key for the worker (sensitive)"
  value       = aws_iam_access_key.worker.secret
  sensitive   = true
}

# =============================================================================
# OIDC Outputs
# =============================================================================

output "oidc_provider_arn" {
  description = "ARN of the GitHub OIDC provider"
  value       = aws_iam_openid_connect_provider.github.arn
}
