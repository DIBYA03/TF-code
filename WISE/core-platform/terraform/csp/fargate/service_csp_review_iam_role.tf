resource "aws_iam_role" "csp_review_execution_role" {
  name = "${module.naming.aws_iam_role}-csp-review-execution-role"

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
    Name        = "${module.naming.aws_iam_role}-csp-review-execution-role"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "csp_review_execution_role" {
  name = "${module.naming.aws_iam_policy}-default-execution"
  role = "${aws_iam_role.csp_review_execution_role.id}"

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
    }
  ]
}
EOF
}

# inline policy to allow access to business upload sqs
resource "aws_iam_role_policy" "csp_review_sqs" {
  name = "${module.naming.aws_iam_role_policy}-business-upload-sqs"
  role = "${aws_iam_role.csp_review_execution_role.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToConsumeReviewSQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "sqs:ChangeMessageVisibility*",
        "sqs:DeleteMessage",
        "sqs:DeleteMessageBatch",
        "kms:GenerateDataKey",
        "sqs:GetQueue*",
        "sqs:ReceiveMessage*",
        "sqs:SendMessage*"
      ],
      "Resource": [
        "${data.aws_kms_alias.default.target_key_arn}",
        "${data.aws_sqs_queue.review.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "csp_review_s3_doc_read_write" {
  name = "${module.naming.aws_iam_role_policy}-s3-read-write"
  role = "${aws_iam_role.csp_review_execution_role.id}"

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
        "s3:ListMultipartUploadParts",
        "s3:PutObject*"
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

resource "aws_iam_role_policy" "csp_review_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.csp_review_execution_role.id}"

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

resource "aws_iam_role_policy" "csp_review_ssm" {
  name = "${module.naming.aws_iam_policy}-ssm-parameters"
  role = "${aws_iam_role.csp_review_execution_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
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
        "${data.aws_kms_alias.default.target_key_arn}",
        "${data.aws_ssm_parameter.bbva_app_env.arn}",
        "${data.aws_ssm_parameter.bbva_app_id.arn}",
        "${data.aws_ssm_parameter.bbva_app_name.arn}",
        "${data.aws_ssm_parameter.bbva_app_secret.arn}",
        "${data.aws_ssm_parameter.bank_rds_db_name.arn}",
        "${data.aws_ssm_parameter.bank_rds_password.arn}",
        "${data.aws_ssm_parameter.bank_rds_user_name.arn}",
        "${data.aws_ssm_parameter.core_rds_db_name.arn}",
        "${data.aws_ssm_parameter.core_rds_password.arn}",
        "${data.aws_ssm_parameter.core_rds_user_name.arn}",
        "${data.aws_ssm_parameter.csp_rds_db_name.arn}",
        "${data.aws_ssm_parameter.csp_rds_master_endpoint.arn}",
        "${data.aws_ssm_parameter.csp_rds_password.arn}",
        "${data.aws_ssm_parameter.csp_rds_port.arn}",
        "${data.aws_ssm_parameter.csp_rds_read_endpoint.arn}",
        "${data.aws_ssm_parameter.csp_rds_username.arn}",
        "${data.aws_ssm_parameter.identity_rds_db_name.arn}",
        "${data.aws_ssm_parameter.identity_rds_user_name.arn}",
        "${data.aws_ssm_parameter.identity_rds_password.arn}",
        "${data.aws_ssm_parameter.rds_master_endpoint.arn}",
        "${data.aws_ssm_parameter.rds_read_endpoint.arn}",
        "${data.aws_ssm_parameter.rds_port.arn}",
        "${data.aws_ssm_parameter.sendgrid_api_key.arn}",
        "${data.aws_ssm_parameter.csp_notification_slack_channel.arn}",
        "${data.aws_ssm_parameter.csp_notification_slack_url.arn}"
      ]
    }
  ]
}
EOF
}
