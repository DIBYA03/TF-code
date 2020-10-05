resource "aws_iam_role" "payments_execution" {
  name = "${module.naming.aws_iam_role}-payments-exec-role"

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
    Name        = "${module.naming.aws_iam_role}-payments-exec-role"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "payments_execution" {
  name = "${module.naming.aws_iam_policy}-default-execution"
  role = "${aws_iam_role.payments_execution.id}"

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
        "${data.aws_ssm_parameter.txn_rds_db_name.arn}",
        "${data.aws_ssm_parameter.txn_rds_user_name.arn}",
        "${data.aws_ssm_parameter.txn_rds_password.arn}",
        "${data.aws_ssm_parameter.stripe_key.arn}",
        "${data.aws_ssm_parameter.stripe_publish_key.arn}",
        "${data.aws_ssm_parameter.stripe_webhook_secret.arn}",
        "${data.aws_ssm_parameter.wise_clearing_account_id.arn}",
        "${data.aws_ssm_parameter.wise_clearing_business_id.arn}",
        "${data.aws_ssm_parameter.wise_clearing_user_id.arn}",
        "${data.aws_ssm_parameter.wise_invoice_email_address.arn}",
        "${data.aws_ssm_parameter.plaid_env.arn}",
        "${data.aws_ssm_parameter.plaid_public_key.arn}",
        "${data.aws_ssm_parameter.plaid_client_id.arn}",
        "${data.aws_ssm_parameter.plaid_secret.arn}",
        "${data.aws_ssm_parameter.segment_web_write_key.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "payments_s3_read" {
  name = "${module.naming.aws_iam_role_policy}-s3-docs-read"
  role = "${aws_iam_role.payments_execution.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "GetObjects",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:Encrypt",
        "kms:GenerateDataKey",
        "s3:AbortMultipartUpload",
        "s3:GetObject",
        "s3:GetObjectTagging",
        "s3:GetObjectVersion",
        "s3:GetObjectVersionTagging",
        "s3:ListBucket",
        "s3:ListBucketMultipartUploads",
        "s3:ListMultipartUploadParts"
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

resource "aws_iam_role_policy" "payments_banking_sqs" {
  name = "${module.naming.aws_iam_role_policy}-banking-sqs"
  role = "${aws_iam_role.payments_execution.name}"

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
