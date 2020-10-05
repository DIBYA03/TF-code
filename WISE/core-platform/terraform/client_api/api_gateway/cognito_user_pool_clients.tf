# Cognito Web Client
resource "aws_cognito_user_pool_client" "web" {
  count                  = "${var.cognito_client_web_enable ? 1 : 0}"
  name                   = "web"
  user_pool_id           = "${aws_cognito_user_pool.default.id}"
  generate_secret        = "${var.cognito_client_web_generate_secret}"
  refresh_token_validity = "${var.cognito_client_web_refresh_token_validity}"
  explicit_auth_flows    = ["${var.cognito_client_web_explicit_auth_flows}"]
}

# Cognito Mobile Client
resource "aws_cognito_user_pool_client" "mobile" {
  count                  = "${var.cognito_client_mobile_enable ? 1 : 0}"
  name                   = "mobile"
  user_pool_id           = "${aws_cognito_user_pool.default.id}"
  generate_secret        = "${var.cognito_client_mobile_generate_secret}"
  refresh_token_validity = "${var.cognito_client_mobile_refresh_token_validity}"
  explicit_auth_flows    = ["${var.cognito_client_mobile_explicit_auth_flows}"]
}

# Cognito Admin Client
resource "aws_cognito_user_pool_client" "admin" {
  count                  = "${var.cognito_client_admin_enable ? 1 : 0}"
  name                   = "admin"
  user_pool_id           = "${aws_cognito_user_pool.default.id}"
  generate_secret        = "${var.cognito_client_admin_generate_secret}"
  refresh_token_validity = "${var.cognito_client_admin_refresh_token_validity}"
  explicit_auth_flows    = ["${var.cognito_client_admin_explicit_auth_flows}"]
}
