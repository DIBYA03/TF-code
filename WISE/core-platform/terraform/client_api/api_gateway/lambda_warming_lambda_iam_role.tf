# IAM roles for the lambda function
resource "aws_iam_role" "lambda_warmer_lambda" {
  name = "${module.naming.aws_iam_role}-lambda-warmer"

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
    Name        = "${module.naming.aws_iam_role}-lambda-warmer"
    Team        = "${var.team}"
  }
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "lambda_warmer_lambda_cw_logs" {
  role       = "${aws_iam_role.lambda_warmer_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "lambda_warmer_lambda_vpc_access" {
  role       = "${aws_iam_role.lambda_warmer_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to invoke lambdas for warming
resource "aws_iam_role_policy" "lambda_warmer_lambda_invoking_lambdas" {
  name = "${module.naming.aws_iam_role_policy}-invoke-lambdas"
  role = "${aws_iam_role.lambda_warmer_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToInvokeLambdas",
      "Effect": "Allow",
      "Action": [
        "lambda:InvokeFunction"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

# ssm parameters to get list of lambdas
resource "aws_iam_role_policy" "lambda_warmer_lambda_ssm" {
  name = "${module.naming.aws_iam_role_policy}-ssm-params"
  role = "${aws_iam_role.lambda_warmer_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToLambdaWarmerSSMParameters",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:GetParametersByPath"
      ],
      "Resource": [
        "arn:aws:ssm:us-west-2:*:parameter/${var.environment}/${aws_lambda_function.lambda_warmer_lambda.function_name}/*"
      ]
    }
  ]
}
EOF
}
