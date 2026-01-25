# Dead Letter Queue
resource "aws_sqs_queue" "tasks_dlq" {
  name                       = "codingworker-tasks-dlq"
  message_retention_seconds  = var.sqs_message_retention * 7 # 7 days for DLQ
  visibility_timeout_seconds = var.sqs_visibility_timeout

  tags = {
    Name = "codingworker-tasks-dlq"
  }
}

# Main Task Queue
resource "aws_sqs_queue" "tasks" {
  name                       = "codingworker-tasks"
  visibility_timeout_seconds = var.sqs_visibility_timeout
  message_retention_seconds  = var.sqs_message_retention
  receive_wait_time_seconds  = 20 # Long polling

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.tasks_dlq.arn
    maxReceiveCount     = var.dlq_max_receive_count
  })

  tags = {
    Name = "codingworker-tasks"
  }
}

# Queue Policy - Allow GitHub Actions role to send messages
resource "aws_sqs_queue_policy" "tasks" {
  queue_url = aws_sqs_queue.tasks.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowGitHubActionsSend"
        Effect    = "Allow"
        Principal = {
          AWS = aws_iam_role.github_actions.arn
        }
        Action   = "sqs:SendMessage"
        Resource = aws_sqs_queue.tasks.arn
      }
    ]
  })
}
