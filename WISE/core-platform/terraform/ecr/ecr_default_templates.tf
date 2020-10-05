data "template_file" "default_lifecycle_policy" {
  template = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Keep last ${var.tagged_image_count_limit} tagged images",
      "selection": {
        "tagStatus": "tagged",
        "tagPrefixList": ["build"],
        "countType": "imageCountMoreThan",
        "countNumber": ${var.tagged_image_count_limit}
      },
      "action": {
        "type": "expire"
      }
    },
    {
      "rulePriority": 2,
      "description": "Expire untagged images when more than ${var.untagged_image_count_limit}",
      "selection": {
        "tagStatus": "untagged",
        "countType": "imageCountMoreThan",
        "countNumber": ${var.untagged_image_count_limit}
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF
}

data "template_file" "default_repo_policy" {
  template = <<EOF
  {
    "Version": "2008-10-17",
    "Statement": [
      {
        "Sid": "allow-wise-cross-account",
        "Effect": "Allow",
        "Principal": {
          "AWS": ${jsonencode(var.allowed_account_principals)}
        },
        "Action": [
          "ecr:BatchCheckLayerAvailability",
          "ecr:BatchGetImage",
          "ecr:DescribeImages",
          "ecr:DescribeRepositories",
          "ecr:GetDownloadUrlForLayer",
          "ecr:GetAuthorizationToken",
          "ecr:ListImages"
        ]
      }
    ]
  }
EOF
}
