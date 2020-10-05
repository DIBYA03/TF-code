resource "aws_iam_role" "clientapi_user_submission_lambda" {
  name = "${module.naming.aws_iam_role}-usr-sub-lambda-${var.api_gw_stage}"

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
    Name        = "${module.naming.aws_iam_role}-usr-sub-lambda-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "clientapi_user_submission_lambda_cw" {
  role       = "${aws_iam_role.clientapi_user_submission_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "clientapi_user_submission_lambda_vpc" {
  role       = "${aws_iam_role.clientapi_user_submission_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "clientapi_business_document_csp_review_sqs" {
  name = "${module.naming.aws_iam_role_policy}-csp-review-sqs"
  role = "${aws_iam_role.clientapi_user_submission_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToCSPReviewSQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey",
        "sqs:GetQueue*",
        "sqs:SendMessage*"
      ],
      "Resource": [
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_kms_alias.csp_default.target_key_arn}",
        "${data.aws_sqs_queue.review.arn}",
        "${data.aws_sqs_queue.segment_analytics.arn}"
      ]
    }
  ]
}
EOF
}
