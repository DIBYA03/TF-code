resource "aws_iam_role" "csp_consumer_document_submit_lambda" {
  name = "${module.naming.aws_iam_role}-cmr-doc-sbt-${var.api_gw_stage}"

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
    Name        = "${module.naming.aws_iam_role}-cmr-doc-sbt-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "csp_consumer_document_submit_lambda_cw" {
  role       = "${aws_iam_role.csp_consumer_document_submit_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "csp_consumer_document_submit_lambda_vpc" {
  role       = "${aws_iam_role.csp_consumer_document_submit_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to business upload sqs
resource "aws_iam_role_policy" "csp_consumer_document_submit_business_upload_sqs" {
  name = "${module.naming.aws_iam_role_policy}-cmr-doc-sbt-bus-upload"
  role = "${aws_iam_role.csp_consumer_document_submit_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToBusinessDocUploadQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey",
        "sqs:GetQueue*",
        "sqs:SendMessage*"
      ],
      "Resource": [
        "${data.aws_kms_alias.env_default.target_key_arn}",
        "${data.aws_sqs_queue.business_document_upload.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "csp_consumer_document_submit_lambda_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.csp_consumer_document_submit_lambda.name}"

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
