data "aws_kms_alias" "default" {
  name = "${var.default_kms_alias}"
}

data "aws_kms_alias" "core_env_default" {
  name = "${var.default_client_api_env_kms_alias}"
}

data "aws_kms_alias" "internal_sqs" {
  name = "alias/${var.environment}-client-api-internal-sqs"
}

data "aws_kms_alias" "documents_bucket" {
  name = "alias/${var.environment}-client-api-s3-documents"
}
