variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "csp"
}

variable "component" {
  default = "api"
}

variable "team" {
  default = "cloud-ops"
}

variable "domain_name" {}
variable "api_gw_stage" {}

variable "vpc_id" {}

variable "vpc_cidr_block" {}

variable "csp_rds_vpc_cidr_block" {}

variable "app_subnet_ids" {
  type = "list"
}

# clientAPI integrations
variable "core_db_cidr_blocks" {
  type = "list"
}

# grpc
variable "grpc_port" {
  default = 3001
}

variable "use_transaction_service" {}
variable "use_banking_service" {}
variable "use_invoice_service" {}

# SNS
variable "sns_non_critical_topic" {}

variable "sns_critical_topic" {}

# Lambda
variable "lambda_timeout" {
  default = 600
}

variable "default_lambda_non_critical_alarm_error_count" {
  default = 1
}

variable "default_lambda_critical_alarm_error_count" {
  default = 10
}

# KMS
variable "default_kms_alias" {
  default = "alias/csp-wise-us-vpc"
}

# kinesis
variable "txn_kinesis_name" {}

variable "txn_kinesis_region" {}

variable "batch_default_timezone" {
  default = "America/Los_Angeles"
}
