# IAM role for business submission lambda function
resource "aws_iam_role" "clientapi_business_lambda" {
  name = "${module.naming.aws_iam_role}-bus-lambda-${var.api_gw_stage}"

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
    Name        = "${module.naming.aws_iam_role}-bus-lambda-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "clientapi_business_lambda_vpc" {
  role       = "${aws_iam_role.clientapi_business_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "clientapi_business_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.clientapi_business_lambda.name}"

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
        "${data.aws_kms_alias.csp_default.target_key_arn}",
        "${data.aws_sqs_queue.review.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "clientapi_business_lambda_s3_doc_read" {
  name = "${module.naming.aws_iam_role_policy}-s3-document-read"
  role = "${aws_iam_role.clientapi_business_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
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
        "${data.aws_kms_alias.documents_bucket.target_key_arn}",
        "${data.aws_s3_bucket.documents.arn}",
        "${data.aws_s3_bucket.documents.arn}/${var.s3_ach_pull_list_config_object}"
      ]
    }
  ]
}
EOF
}
