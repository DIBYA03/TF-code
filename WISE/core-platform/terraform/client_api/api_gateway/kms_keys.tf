resource "aws_kms_key" "documents_bucket" {
  description         = "${var.environment_name} KMS Key for S3 documents bucket"
  enable_key_rotation = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${var.environment}-s3-documents-bucket"
    Team        = "${var.team}"
  }
}

resource "aws_kms_alias" "documents_bucket" {
  name          = "${module.naming.aws_kms_alias}-s3-documents"
  target_key_id = "${aws_kms_key.documents_bucket.key_id}"
}

resource "aws_kms_key" "internal_sqs" {
  description         = "${var.environment_name} KMS Key for Internal SQS"
  enable_key_rotation = true

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "kms:*",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::${data.aws_caller_identity.account.account_id}:root"
      },
      "Resource": "*",
      "Sid": "All All Account Users"
    },
    {
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey*"
      ],
      "Effect": "Allow",
      "Principal": {
        "Service": "sns.amazonaws.com"
      },
      "Resource": "*",
      "Sid": "Allow SNS Topics"
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${var.environment}-internal-sqs"
    Team        = "${var.team}"
  }
}

resource "aws_kms_alias" "internal_sqs" {
  name          = "${module.naming.aws_kms_alias}-internal-sqs"
  target_key_id = "${aws_kms_key.internal_sqs.key_id}"
}

resource "aws_kms_key" "cognito_lambda" {
  description         = "KMS Key for ${var.environment_name} cognito lambdas"
  enable_key_rotation = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${var.environment}-cognito-lambdas"
    Team        = "${var.team}"
  }
}

resource "aws_kms_alias" "cognito_lambda" {
  name          = "${module.naming.aws_kms_alias}-cognito-lambdas"
  target_key_id = "${aws_kms_key.cognito_lambda.key_id}"
}
