resource "aws_iam_role" "batch_account_closure_execution_role" {
  name = "${module.naming.aws_iam_role}-batch-account-closure-execution-role"

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
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_iam_role}-batch-account-closure-execution-role"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "batch_account_closure_execution_role" {
  name = "${module.naming.aws_iam_policy}-default-execution"
  role = "${aws_iam_role.batch_account_closure_execution_role.id}"

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
        "${data.aws_kms_alias.core_env_default.target_key_arn}",
        "${data.aws_ssm_parameter.bbva_app_env.arn}",
        "${data.aws_ssm_parameter.bbva_app_id.arn}",
        "${data.aws_ssm_parameter.bbva_app_name.arn}",
        "${data.aws_ssm_parameter.bbva_app_secret.arn}",
        "${data.aws_ssm_parameter.rds_master_endpoint.arn}",
        "${data.aws_ssm_parameter.rds_read_endpoint.arn}",
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
        "${data.aws_ssm_parameter.wise_clearing_business_id.arn}",
        "${data.aws_ssm_parameter.wise_clearing_user_id.arn}",
        "${data.aws_ssm_parameter.wise_refund_clearing_linked_account_id.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "batch_account_closure_s3_read" {
  name = "${module.naming.aws_iam_role_policy}-s3-doc-read"
  role = "${aws_iam_role.batch_account_closure_execution_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "GetObjects",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
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

resource "aws_iam_role_policy" "batch_account_csp_ssm_param_role" {
  name = "${module.naming.aws_iam_policy}-csp-ssm-param-role"
  role = "${aws_iam_role.batch_account_closure_execution_role.id}"

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
        "${data.aws_kms_alias.default.target_key_arn}",
        "${data.aws_ssm_parameter.csp_rds_read_endpoint.arn}",
        "${data.aws_ssm_parameter.csp_rds_port.arn}",
        "${data.aws_ssm_parameter.csp_rds_master_endpoint.arn}",
        "${data.aws_ssm_parameter.csp_rds_db_name.arn}",
        "${data.aws_ssm_parameter.csp_rds_username.arn}",
        "${data.aws_ssm_parameter.csp_rds_password.arn}"
      ]
    }
  ]
}
EOF
}
