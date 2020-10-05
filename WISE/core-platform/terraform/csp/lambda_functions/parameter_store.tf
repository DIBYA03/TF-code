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

data "aws_ssm_parameter" "sendgrid_api_key" {
  name = "/${var.environment}/sendgrid/api_key"
}

# twilio
data "aws_ssm_parameter" "twilio_account_sid" {
  name            = "/${var.environment}/twilio/account_sid"
  with_decryption = true
}

data "aws_ssm_parameter" "twilio_api_sid" {
  name            = "/${var.environment}/twilio/api_sid"
  with_decryption = true
}

data "aws_ssm_parameter" "twilio_api_secret" {
  name            = "/${var.environment}/twilio/api_secret"
  with_decryption = true
}

data "aws_ssm_parameter" "twilio_sender_phone" {
  name            = "/${var.environment}/twilio/phone_number"
  with_decryption = true
}

# Wise Related
data "aws_ssm_parameter" "wise_clearing_account_id" {
  name            = "/${var.environment}/wise/clearing/account_id"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_clearing_user_id" {
  name            = "/${var.environment}/wise/clearing/user_id"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_clearing_business_id" {
  name            = "/${var.environment}/wise/clearing/business_id"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_promo_clearing_account_id" {
  name            = "/${var.environment}/wise/clearing/promo_account_id"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_promo_linked_clearing_account_id" {
  name            = "/${var.environment}/wise/clearing/promo_linked_account_id"
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

data "aws_ssm_parameter" "intercom_access_token" {
  name            = "/${var.environment}/intercom/access_token"
  with_decryption = true
}

data "aws_ssm_parameter" "cognito_user_pool_id" {
  name = "/${var.environment}/cognito/user_pool/id"
  with_decryption = true
}

data "aws_ssm_parameter" "vgs_cert" {
  name            = "/${var.environment}/vgs/vgs_cert"
  with_decryption = true
}

data "aws_ssm_parameter" "vgs_https_proxy_url" {
  name            = "/${var.environment}/vgs/vgs_https_proxy_url"
  with_decryption = true
}
