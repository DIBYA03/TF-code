# IAM roles for the lambda function
resource "aws_iam_role" "clientapi_lambda" {
  name = "${module.naming.aws_iam_role}-lambda"

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
    Name        = "${module.naming.aws_iam_role}-lambda"
    Team        = "${var.team}"
  }
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "clientapi_lambda_cw_logs" {
  role       = "${aws_iam_role.clientapi_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "clientapi_lambda_vpc_access" {
  role       = "${aws_iam_role.clientapi_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.clientapi_lambda.name}"

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

# Cognito IAM role
resource "aws_iam_role" "cognito" {
  name = "${module.naming.aws_iam_role}-cognito"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "cognito-idp.amazonaws.com"
        },
        "Effect": "Allow",
        "Condition": {
          "StringEquals": {
            "sts:ExternalId": "${random_uuid.cognito.result}"
          }
        }
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_iam_role}-cognito"
    Team        = "${var.team}"
  }
}

data "template_file" "cognito_policy" {
  template = "${file("./policies/cognito.json")}"
}

resource "aws_iam_policy" "cognito_policy" {
  name        = "${module.naming.aws_iam_policy}-cognito"
  description = "IAM policy for ${var.environment_name} Cognito"
  policy      = "${data.template_file.cognito_policy.rendered}"
}

resource "aws_iam_role_policy_attachment" "cognito" {
  role       = "${aws_iam_role.cognito.name}"
  policy_arn = "${aws_iam_policy.cognito_policy.arn}"
}
