# IAM role for stripe webhook lambda function
resource "aws_iam_role" "stripe_webhook_lambda" {
  name = "${module.naming.aws_iam_role}-stp-whk-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_iam_role}-stp-whk-lambda"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "stripe_webhook_lambda" {
  name = "${module.naming.aws_iam_role_policy}-stripe-webhook"
  role = "${aws_iam_role.stripe_webhook_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "kms:CreateGrant",
        "kms:Decrypt",
        "kms:DescribeKey",
        "kms:ReEncryptFrom",
        "kms:ReEncryptTo",
        "sqs:ChangeMessageVisibility",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
        "sqs:ReceiveMessage"
      ],
      "Resource": [
        "${aws_kms_alias.internal_sqs.target_key_arn}",
        "${aws_sqs_queue.sqs_stripe_request_payment.arn}"
      ]
    },
    {
      "Sid": "WriteToS3Docs",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:Encrypt",
        "kms:GenerateDataKey",
        "s3:AbortMultipartUpload",
        "s3:ListBucketMultipartUploads",
        "s3:ListMultipartUploadParts",
        "s3:PutObject*"
      ],
      "Resource": [
        "${aws_kms_alias.documents_bucket.target_key_arn}",
        "${aws_s3_bucket.documents.arn}",
        "${aws_s3_bucket.documents.arn}/*"
      ]
    },
    {
      "Sid": "GetACHPullConfigObject",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "s3:GetObject",
        "s3:GetObjectTagging",
        "s3:GetObjectVersion",
        "s3:GetObjectVersionTagging",
        "s3:ListBucket"
      ],
      "Resource": [
        "${aws_kms_alias.documents_bucket.target_key_arn}",
        "${aws_s3_bucket.documents.arn}",
        "${aws_s3_bucket.documents.arn}/${var.s3_ach_pull_list_config_object}"
      ]
    }
  ]
}
EOF
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "stripe_webhook_lambda_cw" {
  role       = "${aws_iam_role.stripe_webhook_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "stripe_webhook_lambda_vpc" {
  role       = "${aws_iam_role.stripe_webhook_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "stripe_webhook_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.stripe_webhook_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToInternalBankingSQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey",
        "sqs:GetQueue*",
        "sqs:SendMessage*"
      ],
      "Resource": [
        "${aws_kms_alias.internal_sqs.target_key_arn}",
        "${aws_sqs_queue.segment_analytics.arn}"
      ]
    }
  ]
}
EOF
}
