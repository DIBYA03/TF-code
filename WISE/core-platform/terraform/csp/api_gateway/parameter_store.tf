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

data "aws_ssm_parameter" "google_idp_client_id" {
  name            = "/${var.environment}/csp/cognito/google/client_id"
  with_decryption = true
}

data "aws_ssm_parameter" "google_idp_client_secret" {
  name            = "/${var.environment}/csp/cognito/google/client_secret"
  with_decryption = true
}

resource "aws_ssm_parameter" "cognito_pool_id" {
  name      = "/${var.environment}/csp/${terraform.workspace}/cognito/user_pool/id"
  type      = "SecureString"
  key_id    = "${var.default_kms_alias}"
  value     = "${aws_cognito_user_pool.default.id}"
  overwrite = true
}

resource "aws_ssm_parameter" "cognito_client_web_id" {
  name      = "/${var.environment}/csp/${terraform.workspace}/cognito/user_pool/client/web/id"
  type      = "SecureString"
  key_id    = "${var.default_kms_alias}"
  value     = "${aws_cognito_user_pool_client.web.id}"
  overwrite = true
}

data "aws_ssm_parameter" "csp_notification_slack_channel" {
  name            = "/${var.environment}/csp/notifications/slack/channel"
  with_decryption = true
}

data "aws_ssm_parameter" "csp_notification_slack_url" {
  name            = "/${var.environment}/csp/notifications/slack/url"
  with_decryption = true
}
