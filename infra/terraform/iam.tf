# =============================================================================
# GitHub Actions Role (OIDC)
# =============================================================================

# Trust policy for GitHub Actions
data "aws_iam_policy_document" "github_actions_assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:sub"
      values   = [for repo in var.github_repos : "repo:${var.github_org}/${repo}:*"]
    }
  }
}

# Role for GitHub Actions
resource "aws_iam_role" "github_actions" {
  name               = "codingworker-github-actions"
  assume_role_policy = data.aws_iam_policy_document.github_actions_assume_role.json

  tags = {
    Name = "codingworker-github-actions"
  }
}

# Policy: Allow sending messages to SQS
data "aws_iam_policy_document" "github_actions_sqs" {
  statement {
    effect = "Allow"
    actions = [
      "sqs:SendMessage",
      "sqs:GetQueueUrl",
      "sqs:GetQueueAttributes"
    ]
    resources = [aws_sqs_queue.tasks.arn]
  }
}

resource "aws_iam_role_policy" "github_actions_sqs" {
  name   = "sqs-send-message"
  role   = aws_iam_role.github_actions.id
  policy = data.aws_iam_policy_document.github_actions_sqs.json
}

# =============================================================================
# Worker User (MBP)
# =============================================================================

# IAM User for the worker running on MBP
resource "aws_iam_user" "worker" {
  name = "codingworker-worker"

  tags = {
    Name = "codingworker-worker"
  }
}

# Policy: Allow receiving and deleting messages from SQS
data "aws_iam_policy_document" "worker_sqs" {
  statement {
    effect = "Allow"
    actions = [
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueUrl",
      "sqs:GetQueueAttributes",
      "sqs:ChangeMessageVisibility"
    ]
    resources = [aws_sqs_queue.tasks.arn]
  }
}

resource "aws_iam_user_policy" "worker_sqs" {
  name   = "sqs-receive-message"
  user   = aws_iam_user.worker.name
  policy = data.aws_iam_policy_document.worker_sqs.json
}

