resource "aws_iam_role" "stripe_webhook_execution_role" {
  name = "${module.naming.aws_iam_role}stripe-wh-execution-role"

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
    Name        = "${module.naming.aws_iam_role}stripe-wh-execution-role"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "stripe_webhook_execution_role" {
  name = "${module.naming.aws_iam_policy}-default-execution"
  role = "${aws_iam_role.stripe_webhook_execution_role.id}"

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
        "${data.aws_ssm_parameter.bbva_app_env.arn}",
        "${data.aws_ssm_parameter.bbva_app_id.arn}",
        "${data.aws_ssm_parameter.bbva_app_name.arn}",
        "${data.aws_ssm_parameter.bbva_app_secret.arn}",
        "${data.aws_ssm_parameter.rds_master_endpoint.arn}",
        "${data.aws_ssm_parameter.rds_read_endpoint.arn}",
        "${data.aws_ssm_parameter.rds_port.arn}",
        "${data.aws_ssm_parameter.rds_port.arn}",
        "${data.aws_ssm_parameter.core_rds_db_name.arn}",
        "${data.aws_ssm_parameter.core_rds_user_name.arn}",
        "${data.aws_ssm_parameter.core_rds_password.arn}",
        "${data.aws_ssm_parameter.bank_rds_db_name.arn}",
        "${data.aws_ssm_parameter.bank_rds_user_name.arn}",
        "${data.aws_ssm_parameter.bank_rds_password.arn}",
        "${data.aws_ssm_parameter.identity_rds_db_name.arn}",
        "${data.aws_ssm_parameter.identity_rds_user_name.arn}",
        "${data.aws_ssm_parameter.identity_rds_password.arn}",
        "${data.aws_ssm_parameter.stripe_key.arn}",
        "${data.aws_ssm_parameter.stripe_webhook_secret.arn}",
        "${data.aws_ssm_parameter.wise_clearing_business_id.arn}",
        "${data.aws_ssm_parameter.wise_clearing_user_id.arn}",
        "${data.aws_ssm_parameter.wise_support_email_address.arn}",
        "${data.aws_ssm_parameter.wise_support_email_name.arn}"
      ]
    },
    {
      "Sid": "AllowAccessToSQS",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey",
        "sqs:SendMessage",
        "sqs:SendMessageBatch"
      ],
      "Resource": [
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_sqs_queue.stripe_webhook.arn}"
      ]
    }
  ]
}
EOF
}
