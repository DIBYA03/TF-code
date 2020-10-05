# UUID for external ID
resource "random_uuid" "cognito" {}

# Cognito Pool Domain (needed for federated login)
resource "aws_cognito_user_pool_domain" "default" {
  domain       = "${var.cognito_domain_name}"
  user_pool_id = "${aws_cognito_user_pool.default.id}"
}

# Cognito User Pool
resource "aws_cognito_user_pool" "default" {
  name = "${module.naming.aws_cognito_user_pool}"

  mfa_configuration = "${var.cognito_mfa_configuration}"

  alias_attributes = [
    "${var.cognito_username_attributes}",
  ]

  user_pool_add_ons {
    advanced_security_mode = "${var.cognito_advanced_security_mode}"
  }

  admin_create_user_config {
    allow_admin_create_user_only = "${var.cognito_allow_admin_create_user_only}"
    unused_account_validity_days = "${var.cognito_unused_account_validity_days}"
  }

  # SMS Configuration. ExternalID is required, so it's created with UUID
  sms_configuration {
    external_id    = "${random_uuid.cognito.result}"
    sns_caller_arn = "${aws_iam_role.cognito.arn}"
  }

  # Lambda Triggers
  lambda_config {
    post_confirmation    = "${aws_lambda_function.cognitoauth_postconfirm_lambda.arn}"
    pre_sign_up          = "${aws_lambda_function.cognitoauth_presignup_lambda.arn}"
    pre_authentication   = "${aws_lambda_function.cognitoauth_preauthentication_lambda.arn}"
    pre_token_generation = "${aws_lambda_function.cognitoauth_pretoken_lambda.arn}"
  }

  # Password Policy
  password_policy {
    minimum_length    = "${var.cognito_password_minimum_length}"
    require_lowercase = "${var.cognito_password_require_lowercase}"
    require_numbers   = "${var.cognito_password_require_numbers}"
    require_symbols   = "${var.cognito_password_require_symbols}"
    require_uppercase = "${var.cognito_password_require_uppercase}"
  }

  # User Device Configuration
  device_configuration {
    challenge_required_on_new_device      = "${var.cognito_challenge_required_on_new_device}"
    device_only_remembered_on_user_prompt = "${var.cognito_device_only_remembered_on_user_prompt}"
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cognito_user_pool}"
    Team        = "${var.team}"
  }
}
