# IAM role for stripe webhook lambda function
resource "aws_iam_role" "clientapi_banking_business_contact_moneyrequest_lambda" {
  name = "${module.naming.aws_iam_role}-ban-bus-con-myr-lambda-${var.api_gw_stage}"

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
    Name        = "${module.naming.aws_iam_role}-ban-bus-con-myr-lambda"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "clientapi_banking_business_contact_moneyrequest_lambda" {
  name = "${module.naming.aws_iam_role_policy}-ban-bus-con-myr"
  role = "${aws_iam_role.clientapi_banking_business_contact_moneyrequest_lambda.name}"

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
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_sqs_queue.stripe_webhook.arn}"
      ]
    },
    {
      "Sid": "GetObjects",
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
        "${data.aws_kms_alias.documents_bucket.target_key_arn}",
        "${data.aws_s3_bucket.documents.arn}",
        "${data.aws_s3_bucket.documents.arn}/*"
      ]
    }
  ]
}
EOF
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "clientapi_banking_business_contact_moneyrequest_lambda_cw" {
  role       = "${aws_iam_role.clientapi_banking_business_contact_moneyrequest_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "clientapi_banking_business_contact_moneyrequest_lambda_vpc" {
  role       = "${aws_iam_role.clientapi_banking_business_contact_moneyrequest_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "clientapi_banking_business_contact_moneyrequest_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.clientapi_banking_business_contact_moneyrequest_lambda.name}"

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
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_sqs_queue.segment_analytics.arn}"
      ]
    }
  ]
}
EOF
}

# inline policy to allow access to banking sqs
resource "aws_iam_role_policy" "clientapi_banking_business_contact_moneyrequest_banking_sqs" {
  name = "${module.naming.aws_iam_role_policy}-banking-sqs"
  role = "${aws_iam_role.clientapi_banking_business_contact_moneyrequest_lambda.name}"

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
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_sqs_queue.internal_banking.arn}"
      ]
    }
  ]
}
EOF
}
