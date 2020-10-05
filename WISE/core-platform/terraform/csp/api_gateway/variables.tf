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

# VPC
variable "vpc_id" {}

variable "shared_vpc_id" {
  default = "vpc-0426719bb5133d7a0"
}

variable "dev_vpc_id" {
  default = "vpc-0e844a5b87a3cfce5"
}

variable "vpc_cidr_block" {}

variable "csp_rds_vpc_cidr_block" {}

variable "app_subnet_ids" {
  type = "list"
}

# SNS
variable "sns_non_critical_topic" {}

variable "sns_critical_topic" {}

# KMS
variable "default_kms_alias" {}

# clientAPI integrations
variable "core_db_cidr_blocks" {
  type = "list"
}

# grpc
variable "grpc_port" {
  default = 3001
}

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

# Route 53
variable "route53_domain_name" {}

# API Gateway
variable "api_gw_server_description" {}

variable "api_gw_endpoint_configuration" {
  default = "PRIVATE"
}

variable "api_gw_openapi_dir_location" {
  default = "../../../specs/api/csp"
}

variable "api_gw_openapi_file_source" {
  default = "csp-api.yaml"
}

variable "api_gw_stage" {}

# Cognito
variable "cognito_domain_name" {}

variable "cognito_mfa_configuration" {
  default = "ON" # ON, OFF, OPTIONAL
}

variable "cognito_advanced_security_mode" {
  default = "OFF" # OFF, AUDIT or ENFORCED.
}

variable "cognito_username_attributes" {
  type = "list"

  default = [
    "email",
  ]
}

variable "cognito_allow_admin_create_user_only" {
  default = false
}

variable "cognito_unused_account_validity_days" {
  default = 7
}

variable "cognito_password_minimum_length" {
  default = 8
}

variable "cognito_password_require_lowercase" {
  default = true
}

variable "cognito_password_require_numbers" {
  default = true
}

variable "cognito_password_require_symbols" {
  default = true
}

variable "cognito_password_require_uppercase" {
  default = true
}

variable "cognito_challenge_required_on_new_device" {
  default = true
}

variable "cognito_device_only_remembered_on_user_prompt" {
  default = true
}

variable "cognito_auto_verified_attributes" {
  type    = "list"
  default = ["phone_number"]
}

variable "cognito_sms_authentication_message" {
  default = "{####} is your Wise authentication code."
}

variable "cognito_sms_verification_message" {
  default = "{####} is your Wise verification code."
}

# Cognito Clients
variable "cognito_client_web_generate_secret" {
  default = false
}

variable "cognito_client_web_refresh_token_validity" {
  default = 5
}

# Cognito Identity Pool
variable "cognito_client_web_explicit_auth_flows" {
  type    = "list"
  default = []
}

# SQS
variable "document_upload_sqs_kms_data_key_reuse_period_seconds" {
  default = 300 # 5 minutes
}

variable "document_upload_sqs_visibility_timeout_seconds" {
  default = "30"
}

variable "document_upload_sqs_message_retention_seconds" {
  default = "604800" # 7 days
}

variable "document_upload_sqs_dl_message_retention_seconds" {
  default = "1209600" # 14 days
}

variable "document_upload_sqs_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "document_upload_sqs_dl_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "document_upload_sqs_max_message_size" {
  default = 262144 # 256KiB
}

variable "document_upload_sqs_delay_seconds" {
  default = 0
}

variable "document_upload_sqs_receive_wait_time_seconds" {
  default = 0
}

variable "document_upload_sqs_fifo_queue" {
  default = false
}

# review SQS
variable "review_sqs_kms_data_key_reuse_period_seconds" {
  default = 300 # 5 minutes
}

variable "review_sqs_visibility_timeout_seconds" {
  default = "30"
}

variable "review_sqs_message_retention_seconds" {
  default = "604800" # 7 days
}

variable "review_sqs_dl_message_retention_seconds" {
  default = "1209600" # 14 days
}

variable "review_sqs_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "review_sqs_dl_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "review_sqs_max_message_size" {
  default = 262144 # 256KiB
}

variable "review_sqs_delay_seconds" {
  default = 0
}

variable "review_sqs_receive_wait_time_seconds" {
  default = 0
}

variable "review_sqs_fifo_queue" {
  default = false
}
