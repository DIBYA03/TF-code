resource "aws_kms_key" "cognito_lambda" {
  description         = "kms Key for csp ${var.environment_name} cognito lambdas"
  enable_key_rotation = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${var.environment}-cognito-lambdas"
    Team        = "${var.team}"
  }
}

resource "aws_kms_alias" "cognito_lambda" {
  name          = "${module.naming.aws_kms_alias}-cognito-lambdas"
  target_key_id = "${aws_kms_key.cognito_lambda.key_id}"
}

data "aws_kms_alias" "internal_sqs" {
  name = "alias/${var.environment}-client-api-internal-sqs"
}

data "aws_kms_alias" "env_default" {
  name = "${var.default_kms_alias}"
}
