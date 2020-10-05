data "aws_kms_alias" "env_default" {
  name = "${var.default_client_api_kms_alias}"
}

data "aws_kms_alias" "documents_bucket" {
  name = "alias/${var.environment}-client-api-s3-documents"
}

data "aws_kms_alias" "internal_sqs" {
  name = "alias/${var.environment}-client-api-internal-sqs"
}

data "aws_kms_alias" "logging_kinesis" {
  name = "${var.ntf_kinesis_kms_alias}"
}

data "aws_kms_alias" "transactions_kinesis" {
  name = "${var.txn_kinesis_kms_alias}"
}

data "aws_kms_alias" "alloy_kinesis" {
  name = "${var.alloy_kinesis_kms_alias}"
}

data "aws_kms_alias" "csp_default" {
  name = "${var.default_csp_kms_alias}"
}
