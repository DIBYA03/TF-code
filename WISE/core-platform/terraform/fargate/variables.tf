variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "client"
}

variable "component" {
  default = "ecs"
}

variable "team" {
  default = "cloud-ops"
}

# Route 53
variable "public_route53_account_profile" {
  default = "master-us-west-2-saml-roles-deployment"
}

# SNS
variable "non_critical_sns_topic" {}

variable "critical_sns_topic" {}

# VPC
variable "vpc_id" {}

variable "vpc_cidr_block" {}
variable "csp_vpc_cidr_block" {}

variable "shared_vpc_cidr_block" {
  default = "10.2.0.0/16"
}

variable "client_api_rds_vpc_cidr_block" {}

variable "public_subnet_ids" {
  type = "list"
}

variable "app_subnet_ids" {
  type = "list"
}

variable "public_hosted_zone_id" {
  default = "Z09518463MDYVGIEDF9HV"
}

variable "private_hosted_zone_id" {}

variable "default_client_api_kms_alias" {}

variable "default_csp_kms_alias" {}

# kinesis
variable "txn_kinesis_name" {}

variable "txn_kinesis_kms_alias" {}

variable "txn_kinesis_region" {}

variable "ntf_kinesis_name" {}

variable "ntf_kinesis_kms_alias" {}

variable "ntf_kinesis_region" {}

variable "alloy_kinesis_name" {}

variable "alloy_kinesis_kms_alias" {}

variable "alloy_kinesis_region" {}

# all services
variable "cw_log_group_retention_in_days" {}

variable "grpc_port" {
  default = 3001
}

variable "use_transaction_service" {
  default = false
}

variable "use_banking_service" {
  default = false
}

variable "use_invoice_service" {
  default = false
}

variable "services_container_port" {
  default = 3000
}

variable "card_reader_max_request_amount" {
  default = "1000"
}

variable "card_online_max_request_amount" {
  default = "2500"
}

# ach
variable "s3_ach_pull_list_config_object" {
  default = "config/ach_pull_whitelist.json"
}

# payments service
variable "payments_add_monitoring" {
  default = true
}

variable "payments_maintenance_enabled" {
  default = "false"
}

variable "payments_old_domain" {}
variable "payments_domain" {}
variable "payments_image_tag" {}
variable "payments_desired_container_count" {}
variable "payments_min_container_count" {}
variable "payments_max_container_count" {}

variable "payments_name" {
  default = "core-app-payments"
}

variable "payments_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-app-payments"
}

variable "payments_cpu" {
  default = 256
}

variable "payments_mem" {
  default = 512
}

variable "payments_add_owasp10_waf" {
  default = true
}

# stripe_webhook service
variable "stripe_webhook_add_monitoring" {
  default = true
}

variable "stripe_webhook_domain" {}

variable "stripe_webhook_image_tag" {}
variable "stripe_webhook_desired_container_count" {}
variable "stripe_webhook_min_container_count" {}
variable "stripe_webhook_max_container_count" {}

variable "stripe_webhook_name" {
  default = "core-application-stripe-webhook"
}

variable "stripe_webhook_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-stripe-webhook"
}

variable "stripe_webhook_cpu" {
  default = 256
}

variable "stripe_webhook_mem" {
  default = 512
}

variable "stripe_webhook_ip_list" {
  type = "list"

  default = [
    # https://stripe.com/files/ips/ips_webhooks.txt
    "54.187.174.169/32",

    "54.187.205.235/32",
    "54.187.216.72/32",
    "54.241.31.99/32",
    "54.241.31.102/32",
    "54.241.34.107/32",
  ]
}

# BBVA Notification Task
variable "bbva_notification_add_monitoring" {
  default = true
}

variable "bbva_env_name" {}

variable "bbva_notification_name" {
  default = "core-application-bbva-notification"
}

variable "bbva_notification_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-bbva-notification"
}

variable "bbva_notification_cpu" {
  default = 256
}

variable "bbva_notification_mem" {
  default = 512
}

variable "bbva_notification_image_tag" {}
variable "bbva_notification_desired_container_count" {}
variable "bbva_notification_min_container_count" {}
variable "bbva_notification_max_container_count" {}

# App Notification Task
variable "app_notification_add_monitoring" {
  default = true
}

variable "app_notification_name" {
  default = "core-application-app-notification"
}

variable "app_notification_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-app-notifications"
}

variable "app_notification_cpu" {
  default = 256
}

variable "app_notification_mem" {
  default = 512
}

variable "app_notification_image_tag" {}
variable "app_notification_desired_container_count" {}
variable "app_notification_min_container_count" {}
variable "app_notification_max_container_count" {}

# batch timezone
variable "batch_default_timezone" {
  default = "America/Los_Angeles"
}

# batch account
variable "batch_account_add_monitoring" {
  default = true
}

variable "batch_account_name" {
  default = "core-application-batch-account"
}

variable "batch_account_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-batch-account"
}

variable "batch_account_cpu" {
  default = 256
}

variable "batch_account_mem" {
  default = 512
}

variable "batch_account_image_tag" {}

variable "batch_account_desired_container_count" {}
variable "batch_account_min_container_count" {}
variable "batch_account_max_container_count" {}

# batch monitor
variable "batch_monitor_add_monitoring" {
  default = true
}

variable "batch_monitor_name" {
  default = "core-application-batch-monitor"
}

variable "batch_monitor_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-batch-monitor"
}

variable "batch_monitor_cpu" {
  default = 256
}

variable "batch_monitor_mem" {
  default = 512
}

variable "batch_monitor_image_tag" {}

variable "batch_monitor_desired_container_count" {}
variable "batch_monitor_min_container_count" {}
variable "batch_monitor_max_container_count" {}

# batch transaction
variable "batch_transaction_add_monitoring" {
  default = true
}

variable "batch_transaction_name" {
  default = "core-application-batch-transaction"
}

variable "batch_transaction_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-batch-transaction"
}

variable "batch_transaction_cpu" {
  default = 256
}

variable "batch_transaction_mem" {
  default = 512
}

variable "batch_transaction_image_tag" {}

variable "batch_transaction_desired_container_count" {}
variable "batch_transaction_min_container_count" {}
variable "batch_transaction_max_container_count" {}

# Segment Integration Service
variable "segment_analytics_add_monitoring" {
  default = true
}

variable "segment_analytics_name" {
  default = "client-ecs-segment-analytics"
}

variable "segment_analytics_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-segment-analytics"
}

variable "segment_analytics_cpu" {
  default = 256
}

variable "segment_analytics_mem" {
  default = 512
}

variable "segment_analytics_image_tag" {}
variable "segment_analytics_desired_container_count" {}
variable "segment_analytics_min_container_count" {}
variable "segment_analytics_max_container_count" {}

# Shopify Order Service
variable "shopify_order_add_monitoring" {
  default = true
}

variable "shopify_order_name" {
  default = "client-ecs-shopify_order"
}

variable "shopify_order_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-shopify-order"
}

variable "shopify_order_cpu" {
  default = 256
}

variable "shopify_order_mem" {
  default = 512
}

variable "shopify_order_image_tag" {}
variable "shopify_order_desired_container_count" {}
variable "shopify_order_min_container_count" {}
variable "shopify_order_max_container_count" {}

# batch monthly interest
variable "batch_monthly_interest_add_monitoring" {
  default = true
}

variable "batch_monthly_interest_name" {
  default = "core-application-batch-account"
}

variable "batch_monthly_interest_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-batch-monthly-interest"
}

variable "batch_monthly_interest_cpu" {
  default = 256
}

variable "batch_monthly_interest_mem" {
  default = 512
}

variable "batch_monthly_interest_image_tag" {}

variable "batch_monthly_interest_cw_event_scheduled_expression" {
  default = "cron(0 08 1 * ? *)"
}

variable "batch_monthly_interest_desired_container_count" {}
variable "batch_monthly_interest_min_container_count" {}
variable "batch_monthly_interest_max_container_count" {}

# daily batch step function
variable "batch_daily_cw_event_scheduled_expression" {
  default = "cron(0 08 2-31 * ? *)"
}

# merchant logo service
variable "merchant_logo_add_monitoring" {
  default = true
}

variable "merchant_logo_domain" {}
variable "merchant_logo_image_tag" {}
variable "merchant_logo_desired_container_count" {}
variable "merchant_logo_min_container_count" {}
variable "merchant_logo_max_container_count" {}

variable "merchant_logo_name" {
  default = "core-app-payments"
}

variable "merchant_logo_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-merchant-logo"
}

variable "merchant_logo_cpu" {
  default = 256
}

variable "merchant_logo_mem" {
  default = 512
}

variable "merchant_logo_add_owasp10_waf" {
  default = true
}

# hello sign service
variable "hello_sign_add_monitoring" {
  default = true
}

variable "hello_sign_domain" {}
variable "hello_sign_image_tag" {}
variable "hello_sign_desired_container_count" {}
variable "hello_sign_min_container_count" {}
variable "hello_sign_max_container_count" {}

variable "hello_sign_name" {
  default = "core-application-hello-sign"
}

variable "hello_sign_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-hello-sign"
}

variable "hello_sign_cpu" {
  default = 256
}

variable "hello_sign_mem" {
  default = 512
}

variable "hello_sign_webhook_ip_list" {
  type = "list"

  # https://faq.hellosign.com/hc/en-us/articles/115012650807-How-can-I-secure-my-callback-url-
  default = [
    "52.200.252.64/32",
    "34.198.117.22/32",
    "34.198.205.50/32",
  ]
}

# signature service
variable "signature_add_monitoring" {
  default = true
}

variable "signature_image_tag" {}
variable "signature_desired_container_count" {}
variable "signature_min_container_count" {}
variable "signature_max_container_count" {}

variable "signature_name" {
  default = "core-application-signature"
}

variable "signature_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-signature"
}

variable "signature_cpu" {
  default = 256
}

variable "signature_mem" {
  default = 512
}

# batch analytics
variable "batch_analytics_add_monitoring" {
  default = true
}

variable "batch_analytics_name" {
  default = "core-application-batch-analytics"
}

variable "batch_analytics_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-batch-analytics"
}

variable "batch_analytics_cpu" {
  default = 256
}

variable "batch_analytics_mem" {
  default = 512
}

variable "batch_analytics_image_tag" {}

variable "batch_analytics_desired_container_count" {}
variable "batch_analytics_min_container_count" {}
variable "batch_analytics_max_container_count" {}

variable "batch_analytics_cw_event_scheduled_expression" {
  default = "cron(45 11 * * ? *)"
}
