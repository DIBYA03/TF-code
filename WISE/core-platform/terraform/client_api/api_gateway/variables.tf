variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}
variable "bbva_sqs_environment" {}

variable "application" {
  default = "client"
}

variable "component" {
  default = "api"
}

variable "team" {
  default = "cloud-ops"
}

# VPC
variable "vpc_id" {}

variable "vpc_cidr_block" {}

variable "core_rds_cidr_block" {
  default = ""
}

variable "app_subnet_ids" {
  type = "list"
}

# SNS
variable "sns_non_critical_topic" {}

variable "sns_critical_topic" {}

# KMS
variable "default_kms_alias" {}

# ach
variable "s3_ach_pull_list_config_object" {
  default = "config/ach_pull_whitelist.json"
}

# bbva requeue
variable "s3_bbva_requeue_object" {
  default = "config/bbva_requeue.pdf"
}

# Lambda
variable "lambda_timeout" {}

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
  default = "../../../specs/api/clientapi"
}

variable "api_gw_openapi_file_source" {
  default = "client-api.yaml"
}

variable "api_gw_stage" {}

# Cognito
variable "cognito_sms_prepend_message" {
  default = "is your Wise verification code. It will expire in 15 minutes."
}

variable "cognito_mfa_configuration" {
  default = "ON" # ON, OFF, OPTIONAL
}

variable "cognito_advanced_security_mode" {
  default = "AUDIT" # OFF, AUDIT or ENFORCED
}

variable "cognito_username_attributes" {
  default = ["phone_number"]
  type    = "list"
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
  default = ["phone_number"]
  type    = "list"
}

variable "cognito_unused_account_validity_days" {
  default = 2
}

# Cognito Clients
variable "cognito_client_web_enable" {
  default = true
}

variable "cognito_client_web_generate_secret" {
  default = false
}

variable "cognito_client_web_refresh_token_validity" {
  default = 1
}

variable "cognito_client_web_explicit_auth_flows" {
  default = []
  type    = "list"
}

variable "cognito_client_mobile_enable" {
  default = true
}

variable "cognito_client_mobile_generate_secret" {
  default = true
}

variable "cognito_client_mobile_refresh_token_validity" {
  default = 30
}

variable "cognito_client_mobile_explicit_auth_flows" {
  default = []
  type    = "list"
}

variable "cognito_client_admin_enable" {
  default = false
}

variable "cognito_client_admin_generate_secret" {
  default = false
}

variable "cognito_client_admin_refresh_token_validity" {
  default = 5
}

variable "cognito_client_admin_explicit_auth_flows" {
  default = ["ADMIN_NO_SRP_AUTH"]
  type    = "list"
}

# SQS
variable "internal_sqs_visibility_timeout_seconds" {
  default = "30"
}

variable "internal_sqs_message_retention_seconds" {
  default = "604800" # 7 days
}

variable "internal_sqs_dl_message_retention_seconds" {
  default = "1209600" # 14 days
}

variable "internal_sqs_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "internal_sqs_dl_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "internal_sqs_max_message_size" {
  default = 262144 # 256KiB
}

variable "internal_sqs_delay_seconds" {
  default = 0
}

variable "internal_sqs_receive_wait_time_seconds" {
  default = 0
}

variable "internal_sqs_fifo_queue" {
  default = false
}

# SQS
variable "internal_sqs_kms_data_key_reuse_period_seconds" {
  default = 300 # 5 minutes
}

# Stripe request payment SQS
variable "stripe_request_payment_sqs_visibility_timeout_seconds" {
  default = 60 # This needs to be at least what the stripe lambda is
}

# Segment Analytics SQS
variable "segment_analytics_sqs_visibility_timeout_seconds" {
  default = "30"
}

variable "segment_analytics_sqs_message_retention_seconds" {
  default = "604800" # 7 days
}

variable "segment_analytics_sqs_dl_message_retention_seconds" {
  default = "1209600" # 14 days
}

variable "segment_analytics_sqs_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "segment_analytics_sqs_dl_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "segment_analytics_sqs_max_message_size" {
  default = 262144 # 256KiB
}

variable "segment_analytics_sqs_delay_seconds" {
  default = 0
}

variable "segment_analytics_sqs_receive_wait_time_seconds" {
  default = 0
}

variable "segment_analytics_sqs_fifo_queue" {
  default = false
}

variable "segment_analytics_sqs_kms_data_key_reuse_period_seconds" {
  default = 300 # 5 minutes
}

# Shopify Order SQS
variable "shopify_order_sqs_visibility_timeout_seconds" {
  default = "30"
}

variable "shopify_order_sqs_message_retention_seconds" {
  default = "604800" # 7 days
}

variable "shopify_order_sqs_dl_message_retention_seconds" {
  default = "1209600" # 14 days
}

variable "shopify_order_sqs_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "shopify_order_sqs_dl_content_based_deduplication" {
  default = false # use a SHA-256 hash to generate the message deduplication ID using the body of the message
}

variable "shopify_order_sqs_max_message_size" {
  default = 262144 # 256KiB
}

variable "shopify_order_sqs_delay_seconds" {
  default = 0
}

variable "shopify_order_sqs_receive_wait_time_seconds" {
  default = 0
}

variable "shopify_order_sqs_fifo_queue" {
  default = false
}

variable "shopify_order_sqs_kms_data_key_reuse_period_seconds" {
  default = 300 # 5 minutes
}

# BBVA Notifications SQS
variable "bbva_wise_profile" {}

variable "bbva_sqs_visibility_timeout_seconds" {
  default = 30
}

variable "bbva_sqs_message_retention_seconds" {
  default = 604800
}

variable "bbva_sqs_dl_message_retention_seconds" {
  default = 1209600
}

variable "bbva_sqs_content_based_deduplication" {
  default = false
}

variable "bbva_sqs_dl_content_based_deduplication" {
  default = false
}

variable "bbva_sqs_max_message_size" {
  default = 262144
}

variable "bbva_sqs_delay_seconds" {
  default = 0
}

variable "bbva_sqs_receive_wait_time_seconds" {
  default = 0
}

variable "bbva_sqs_fifo_queue" {
  default = false
}

variable "bbva_sqs_kms_data_key_reuse_period_seconds" {
  default = 300
}

# Banking Notifications SQS
variable "banking_notifications_sqs_delay_seconds" {
  default = 3
}

variable "banking_notifications_dead_letter_sqs_delay_seconds" {
  default = 0
}

# BBVA Notifications
variable "bbva_notifications_env" {}

# Lambda Warmer
variable "cognito_lambda_warmer_payload" {
  default = "{ \"request\": { \"userAttributes\": { \"lambda_warmer\": \"true\" } } }"
}

# grpc
variable "grpc_port" {
  default = 3001
}

variable "use_transaction_service" {}
variable "use_banking_service" {}

variable "batch_default_timezone" {
  default = "America/Los_Angeles"
}

variable "use_invoice_service" {
  default = false
}