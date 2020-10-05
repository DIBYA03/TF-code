data "aws_kms_alias" "documents_bucket" {
  name = "${module.naming.aws_kms_alias}-s3-documents"
}

data "aws_kms_alias" "internal_sqs" {
  name = "${module.naming.aws_kms_alias}-internal-sqs"
}

data "aws_kms_alias" "csp_default" {
  name = "${var.csp_kms_alias}"
}

resource "aws_kms_key" "api_lambda" {
  description         = "KMS Key for ${var.environment_name} lambda functions"
  enable_key_rotation = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${var.environment}-api-gw-lambda-kms-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}
