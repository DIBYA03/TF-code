variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "client"
}

variable "component" {
  default = "api"
}

variable "team" {
  default = "cloud-ops"
}

variable "vpc_id" {}
variable "vpc_cidr_block" {}

variable "app_subnet_ids" {
  type = "list"
}

# CSP integrations
variable "csp_kms_alias" {}

variable "csp_environment" {
  description = "since prod csp is used in beta-prod and prd, we need to know this"
}

# ach
variable "s3_ach_pull_list_config_object" {
  default = "config/ach_pull_whitelist.json"
}

# SNS
variable "sns_non_critical_topic" {}

variable "sns_critical_topic" {}

# Lambda
variable "lambda_timeout" {}

variable "enable_user_delete_lambda" {
  default = true
}

variable "default_lambda_non_critical_alarm_error_count" {
  default = 1
}

variable "default_lambda_critical_alarm_error_count" {
  default = 10
}

# API Gateway
variable "api_gw_stage" {}

variable "api_gw_5XX_error_alarm_non_critical_threshold" {}

variable "api_gw_5XX_error_alarm_critical_threshold" {}

variable "api_gw_4XX_error_alarm_non_critical_threshold" {}

variable "api_gw_4XX_error_alarm_critical_threshold" {}

variable "api_gw_latency_alarm_threshold" {}

# kinesis
variable "txn_kinesis_name" {}

variable "txn_kinesis_region" {}

# Wise Clearing
variable "wise_clearing_max_request_amount" {
  default = "10000"
}

variable "card_reader_max_request_amount" {
  default = "1000"
}

variable "card_online_max_request_amount" {
  default = "2500"
}

variable "max_check_amount_allowed" {
  default = "2500"
}

variable "ach_max_amount" {
  default = "5000"
}

# grpc
variable "grpc_port" {
  default = 3001
}

variable "use_transaction_service" {}
variable "use_banking_service" {}
variable "use_invoice_service" {
  default = false
}
