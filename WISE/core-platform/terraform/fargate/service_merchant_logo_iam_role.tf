resource "aws_iam_role" "merchant_logo_execution" {
  name = "${module.naming.aws_iam_role}-merchant-logo-exec-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
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
    Name        = "${module.naming.aws_iam_role}-merchant-logo-exec-role"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "merchant_logo_execution" {
  name = "${module.naming.aws_iam_policy}-default-execution"
  role = "${aws_iam_role.merchant_logo_execution.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    },
    {
      "Sid": "AllowAccessToEnvSSMParametersOnly",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:GetParametersByPath"
      ],
      "Resource": [
        "${data.aws_kms_alias.env_default.target_key_arn}",
        "${data.aws_ssm_parameter.clear_bit_api_key.arn}"
      ]
    }
  ]
}
EOF
}
