data "aws_kms_alias" "documents_bucket" {
  name = "alias/${var.environment}-client-api-s3-documents"
}

resource "aws_kms_key" "lambda_default" {
  description         = "${var.environment_name} KMS Key for CSP Lambda functions"
  enable_key_rotation = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${var.environment}-csp-api-lambda-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

data "aws_kms_alias" "internal_sqs" {
  name = "alias/${var.environment}-client-api-internal-sqs"
}

data "aws_kms_alias" "env_default" {
  name = "${var.default_kms_alias}"
}
