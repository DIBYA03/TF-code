# UUID for external ID
resource "random_uuid" "cognito" {}

# Cognito User Pool
resource "aws_cognito_user_pool" "default" {
  name = "${module.naming.aws_cognito_user_pool}"

  # Allow verification automatically with MFA, instead of manual verification
  auto_verified_attributes = ["${var.cognito_auto_verified_attributes}"]

  # User Device Configuration
  device_configuration {
    challenge_required_on_new_device      = "${var.cognito_challenge_required_on_new_device}"
    device_only_remembered_on_user_prompt = "${var.cognito_device_only_remembered_on_user_prompt}"
  }

  # Lambda Triggers
  lambda_config {
    custom_message       = "${aws_lambda_function.cognitoauth_custommessage_lambda.arn}"
    post_confirmation    = "${aws_lambda_function.cognitoauth_postconfirm_lambda.arn}"
    pre_sign_up          = "${aws_lambda_function.cognitoauth_presignup_lambda.arn}"
    pre_token_generation = "${aws_lambda_function.cognitoauth_pretoken_lambda.arn}"
  }

  mfa_configuration = "${var.cognito_mfa_configuration}"

  # Password Policy
  password_policy {
    minimum_length    = "${var.cognito_password_minimum_length}"
    require_lowercase = "${var.cognito_password_require_lowercase}"
    require_numbers   = "${var.cognito_password_require_numbers}"
    require_symbols   = "${var.cognito_password_require_symbols}"
    require_uppercase = "${var.cognito_password_require_uppercase}"
  }

  # SMS Configuration. ExternalID is required, so it's created with UUID
  sms_configuration {
    external_id    = "${random_uuid.cognito.result}"
    sns_caller_arn = "${aws_iam_role.cognito.arn}"
  }

  # Email addresses or phone numbers can be specified as usernames
  # Using instead of alias so only phone number can be used as a username
  username_attributes = [
    "${var.cognito_username_attributes}",
  ]

  user_pool_add_ons {
    advanced_security_mode = "${var.cognito_advanced_security_mode}"
  }

  # How quickly should temporary passwords set by administrators expire if not used?
  admin_create_user_config {
    unused_account_validity_days = "${var.cognito_unused_account_validity_days}"
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cognito_user_pool}"
    Team        = "${var.team}"
  }
}
