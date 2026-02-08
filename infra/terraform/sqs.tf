# Dead Letter Queue
resource "aws_sqs_queue" "tasks_dlq" {
  name                       = "codingworker-tasks-dlq"
  message_retention_seconds  = 1209600 # 14 days (SQS maximum)
  visibility_timeout_seconds = var.sqs_visibility_timeout
  sqs_managed_sse_enabled    = true

  tags = {
    Name = "codingworker-tasks-dlq"
  }
}

# Allow only the main queue to use this DLQ
resource "aws_sqs_queue_redrive_allow_policy" "tasks_dlq" {
  queue_url = aws_sqs_queue.tasks_dlq.id

  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue"
    sourceQueueArns   = [aws_sqs_queue.tasks.arn]
  })
}

# Main Task Queue
resource "aws_sqs_queue" "tasks" {
  name                       = "codingworker-tasks"
  visibility_timeout_seconds = var.sqs_visibility_timeout
  message_retention_seconds  = var.sqs_message_retention
  receive_wait_time_seconds  = 20 # Long polling
  sqs_managed_sse_enabled    = true

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.tasks_dlq.arn
    maxReceiveCount     = var.sqs_max_receive_count
  })

  tags = {
    Name = "codingworker-tasks"
  }
}

# Main queue policy
resource "aws_sqs_queue_policy" "tasks" {
  queue_url = aws_sqs_queue.tasks.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "DenyNonSSLAccess"
        Effect    = "Deny"
        Principal = "*"
        Action    = "sqs:*"
        Resource  = aws_sqs_queue.tasks.arn
        Condition = {
          Bool = {
            "aws:SecureTransport" = "false"
          }
        }
      }
    ]
  })
}

# DLQ policy
resource "aws_sqs_queue_policy" "tasks_dlq" {
  queue_url = aws_sqs_queue.tasks_dlq.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "DenyNonSSLAccess"
        Effect    = "Deny"
        Principal = "*"
        Action    = "sqs:*"
        Resource  = aws_sqs_queue.tasks_dlq.arn
        Condition = {
          Bool = {
            "aws:SecureTransport" = "false"
          }
        }
      }
    ]
  })
}
