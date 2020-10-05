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

data "aws_ssm_parameter" "redis_endpoint" {
  name            = "/${var.environment}/redis/endpoint"
  with_decryption = true
}

data "aws_ssm_parameter" "redis_port" {
  name            = "/${var.environment}/redis/port"
  with_decryption = true
}

data "aws_ssm_parameter" "redis_password" {
  name            = "/${var.environment}/redis/password"
  with_decryption = true
}

# Stripe
data "aws_ssm_parameter" "stripe_key" {
  name            = "/${var.environment}/stripe/key"
  with_decryption = true
}

data "aws_ssm_parameter" "stripe_webhook_secret" {
  name            = "/${var.environment}/stripe/webhook_secret"
  with_decryption = true
}

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

data "aws_ssm_parameter" "wise_support_phone" {
  name            = "/${var.environment}/wise/support_phone"
  with_decryption = true
}

data "aws_ssm_parameter" "sendgrid_api_key" {
  name            = "/${var.environment}/sendgrid/api_key"
  with_decryption = true
}

data "aws_ssm_parameter" "payments_url" {
  name            = "/${var.environment}/wise/payments/url"
  with_decryption = true
}

resource "aws_ssm_parameter" "cognito_pool_id" {
  name      = "/${var.environment}/cognito/user_pool/id"
  type      = "SecureString"
  key_id    = "${var.default_kms_alias}"
  value     = "${aws_cognito_user_pool.default.id}"
  overwrite = true
}

resource "aws_ssm_parameter" "cognito_client_web_id" {
  count     = "${var.cognito_client_web_enable ? 1 : 0}"
  name      = "/${var.environment}/cognito/user_pool/client/web/id"
  type      = "SecureString"
  key_id    = "${var.default_kms_alias}"
  value     = "${aws_cognito_user_pool_client.web.id}"
  overwrite = true
}

resource "aws_ssm_parameter" "cognito_client_mobile_id" {
  count     = "${var.cognito_client_mobile_enable ? 1 : 0}"
  name      = "/${var.environment}/cognito/user_pool/client/mobile/id"
  type      = "SecureString"
  key_id    = "${var.default_kms_alias}"
  value     = "${aws_cognito_user_pool_client.mobile.id}"
  overwrite = true
}

resource "aws_ssm_parameter" "cognito_client_admin_id" {
  count     = "${var.cognito_client_admin_enable ? 1 : 0}"
  name      = "/${var.environment}/cognito/user_pool/client/admin/id"
  type      = "SecureString"
  key_id    = "${var.default_kms_alias}"
  value     = "${aws_cognito_user_pool_client.admin.id}"
  overwrite = true
}

resource "aws_ssm_parameter" "banking_notifications_sqs" {
  name      = "/${var.environment}/sqs/banking_notifications/url"
  type      = "SecureString"
  key_id    = "${var.default_kms_alias}"
  value     = "${aws_sqs_queue.banking_notifications.id}"
  overwrite = true
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
