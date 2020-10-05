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

#CSP DB
data "aws_ssm_parameter" "csp_rds_read_endpoint" {
  name            = "/${var.environment}/csp/rds/read_endpoint"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_rds_port" {
  name            = "/${var.environment}/csp/rds/db_port"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_rds_master_endpoint" {
  name            = "/${var.environment}/csp/rds/master_endpoint"
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

# Money Request
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

data "aws_ssm_parameter" "wise_clearing_linked_account_id" {
  name            = "/${var.environment}/wise/clearing/linked_account_id"
  with_decryption = true
}

# Stripe
data "aws_ssm_parameter" "stripe_webhook_secret" {
  name            = "/${var.environment}/stripe/webhook_secret"
  with_decryption = true
}

data "aws_ssm_parameter" "stripe_key" {
  name            = "/${var.environment}/stripe/key"
  with_decryption = true
}

data "aws_ssm_parameter" "stripe_publish_key" {
  name            = "/${var.environment}/stripe/publish_key"
  with_decryption = true
}

data "aws_ssm_parameter" "firebase_config" {
  name            = "/${var.environment}/firebase/config"
  with_decryption = true
}

data "aws_ssm_parameter" "sendgrid_api_key" {
  name            = "/${var.environment}/sendgrid/api_key"
  with_decryption = true
}

data "aws_ssm_parameter" "segment_write_key" {
  name            = "/${var.environment}/segment/write_key"
  with_decryption = true
}

data "aws_ssm_parameter" "segment_web_write_key" {
  name            = "/${var.environment}/segment/web_write_key"
  with_decryption = true
}

data "aws_ssm_parameter" "wise_invoice_email_address" {
  name            = "/${var.environment}/wise/invoice_email_address"
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

data "aws_ssm_parameter" "alloy_app_env" {
  name            = "/${var.environment}/alloy/app_env"
  with_decryption = true
}

data "aws_ssm_parameter" "alloy_app_name" {
  name            = "/${var.environment}/alloy/app_name"
  with_decryption = true
}

data "aws_ssm_parameter" "alloy_app_id" {
  name            = "/${var.environment}/alloy/app_id"
  with_decryption = true
}

data "aws_ssm_parameter" "alloy_app_secret" {
  name            = "/${var.environment}/alloy/app_secret"
  with_decryption = true
}

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

data "aws_ssm_parameter" "payments_url" {
  name            = "/${var.environment}/wise/payments/url"
  with_decryption = true
}

data "aws_ssm_parameter" "clear_bit_api_key" {
  name = "/${var.environment}/clear_bit/api_key"
}

# Plaid
data "aws_ssm_parameter" "plaid_env" {
  name            = "/${var.environment}/plaid/env"
  with_decryption = true
}

data "aws_ssm_parameter" "plaid_client_id" {
  name            = "/${var.environment}/plaid/client_id"
  with_decryption = true
}

data "aws_ssm_parameter" "plaid_secret" {
  name            = "/${var.environment}/plaid/secret"
  with_decryption = true
}

data "aws_ssm_parameter" "plaid_public_key" {
  name            = "/${var.environment}/plaid/public_key"
  with_decryption = true
}

# hello sign
# hello sign
data "aws_ssm_parameter" "hellosign_api_key" {
  name            = "/${var.environment}/hellosign/api_key"
  with_decryption = true
}
