resource "aws_iam_role" "cognito_lambda_default" {
  name = "${module.naming.aws_iam_role}-cognito-lambda-default"

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
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "cognito_lambda_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.cognito_lambda_default.name}"

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

resource "aws_iam_role_policy_attachment" "cognito_lambda_default_cw_logs" {
  role       = "${aws_iam_role.cognito_lambda_default.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "cognito_lambda_default_vpc_access" {
  role       = "${aws_iam_role.cognito_lambda_default.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}
