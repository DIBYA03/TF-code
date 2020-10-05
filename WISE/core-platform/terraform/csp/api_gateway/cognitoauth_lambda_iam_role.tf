resource "aws_iam_role" "cognitoauth_lambda" {
  name = "${module.naming.aws_iam_role}-cognitoauth-lambda"

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

resource "aws_iam_role_policy_attachment" "cognitoauth_lambda_cw_logs" {
  role       = "${aws_iam_role.cognitoauth_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "cognitoauth_lambda_vpc_access" {
  role       = "${aws_iam_role.cognitoauth_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}
