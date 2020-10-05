# Cognito Web Client
resource "aws_cognito_user_pool_client" "web" {
  name                   = "web"
  user_pool_id           = "${aws_cognito_user_pool.default.id}"
  generate_secret        = "${var.cognito_client_web_generate_secret}"
  refresh_token_validity = "${var.cognito_client_web_refresh_token_validity}"
  explicit_auth_flows    = ["${var.cognito_client_web_explicit_auth_flows}"]  # USER_PASSWORD_AUTH

  callback_urls = [
    "${var.environment == "dev" ? "http://localhost:8080/user/login" : "https://${var.route53_domain_name}/user/login"}", # only want on dev environment
    "https://${var.route53_domain_name}/user/login",
  ]

  logout_urls = [
    "${var.environment == "dev" ? "http://localhost:8080/user/logout" : "https://${var.route53_domain_name}/user/logout"}", # only want on dev environment
    "https://${var.route53_domain_name}/user/logout",
  ]

  allowed_oauth_flows                  = ["code"]
  allowed_oauth_flows_user_pool_client = true

  allowed_oauth_scopes = [
    "email",
    "openid",
    "profile",
    "aws.cognito.signin.user.admin",
  ]

  supported_identity_providers = ["Google"]

  depends_on = [
    "aws_cognito_identity_provider.google",
  ]
}
