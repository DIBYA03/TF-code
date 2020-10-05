# RDS
data "aws_ssm_parameter" "rds_master_endpoint" {
  name            = "/${var.environment}/rds/master_endpoint"
  with_decryption = true
}

data "aws_ssm_parameter" "rds_read_endpoint" {
  name            = "/${var.environment}/rds/read_endpoint"
  with_decryption = true
}

data "aws_ssm_parameter" "rds_port" {
  name            = "/${var.environment}/rds/db_port"
  with_decryption = true
}

data "aws_ssm_parameter" "core_rds_db_name" {
  name            = "/${var.environment}/rds/core_db_name"
  with_decryption = true
}

data "aws_ssm_parameter" "core_rds_user_name" {
  name            = "/${var.environment}/rds/core_username"
  with_decryption = true
}

data "aws_ssm_parameter" "core_rds_password" {
  name            = "/${var.environment}/rds/core_password"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_rds_master_endpoint" {
  name            = "/${var.environment}/csp/rds/master_endpoint"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_rds_read_endpoint" {
  name            = "/${var.environment}/csp/rds/read_endpoint"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_rds_port" {
  name            = "/${var.environment}/csp/rds/db_port"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_rds_db_name" {
  name            = "/${var.environment}/csp/rds/db_name"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_rds_username" {
  name            = "/${var.environment}/csp/rds/username"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_rds_password" {
  name            = "/${var.environment}/csp/rds/password"
  with_decryption = true
}

data "aws_ssm_parameter" "bank_rds_db_name" {
  name            = "/${var.environment}/rds/bank_db_name"
  with_decryption = true
}

data "aws_ssm_parameter" "bank_rds_user_name" {
  name            = "/${var.environment}/rds/bank_username"
  with_decryption = true
}

data "aws_ssm_parameter" "bank_rds_password" {
  name            = "/${var.environment}/rds/bank_password"
  with_decryption = true
}

data "aws_ssm_parameter" "identity_rds_db_name" {
  name            = "/${var.environment}/rds/identity_db_name"
  with_decryption = true
}

data "aws_ssm_parameter" "identity_rds_user_name" {
  name            = "/${var.environment}/rds/identity_username"
  with_decryption = true
}

data "aws_ssm_parameter" "identity_rds_password" {
  name            = "/${var.environment}/rds/identity_password"
  with_decryption = true
}

data "aws_ssm_parameter" "txn_rds_db_name" {
  name            = "/${var.environment}/rds/txn_db_name"
  with_decryption = true
}

data "aws_ssm_parameter" "txn_rds_user_name" {
  name            = "/${var.environment}/rds/txn_username"
  with_decryption = true
}

data "aws_ssm_parameter" "txn_rds_password" {
  name            = "/${var.environment}/rds/txn_password"
  with_decryption = true
}

data "aws_ssm_parameter" "sendgrid_api_key" {
  name            = "/${var.environment}/sendgrid/api_key"
  with_decryption = true
}

# BBVA
data "aws_ssm_parameter" "bbva_app_env" {
  name            = "/${var.environment}/bbva/app_env"
  with_decryption = true
}

data "aws_ssm_parameter" "bbva_app_id" {
  name            = "/${var.environment}/bbva/app_id"
  with_decryption = true
}

data "aws_ssm_parameter" "bbva_app_name" {
  name            = "/${var.environment}/bbva/app_name"
  with_decryption = true
}

data "aws_ssm_parameter" "bbva_app_secret" {
  name            = "/${var.environment}/bbva/app_secret"
  with_decryption = true
}

data "aws_ssm_parameter" "bbva_requeue_s3_object" {
  name            = "/${var.environment}/dev/bbva/app_env"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_support_email_address" {
  name            = "/${var.environment}/wise/support_email_address"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_support_email_name" {
  name            = "/${var.environment}/wise/support_email_name"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_notification_slack_channel" {
  name            = "/${var.environment}/csp/notifications/slack/channel"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_notification_slack_url" {
  name            = "/${var.environment}/csp/notifications/slack/url"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_clearing_business_id" {
  name            = "/${var.environment}/wise/clearing/business_id"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_clearing_linked_account_id" {
  name            = "/${var.environment}/wise/clearing/linked_account_id"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_refund_clearing_linked_account_id" {
  name            = "/${var.environment}/wise/clearing/refund/linked_account_id"
  with_decryption = true
}


data "aws_ssm_parameter" "wise_clearing_user_id" {
  name            = "/${var.environment}/wise/clearing/user_id"
  with_decryption = true
}
