aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "sbx"

environment_name = "sbx"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:345501831767:wise-us-sbx-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:345501831767:wise-us-sbx-critical-sns"

# VPC
vpc_id = "vpc-06a457dbfa4e8ed8d"

vpc_cidr_block = "10.3.0.0/16"

csp_vpc_cidr_block = "10.4.0.0/16"

client_api_rds_vpc_cidr_block = "10.3.0.0/16"

public_subnet_ids = [
  "subnet-0e9800a72db48a6e9",
  "subnet-0e00be360bd3fd889",
  "subnet-008324c8e6e1e531b",
]

app_subnet_ids = [
  "subnet-0981342e0da7bbe7d",
  "subnet-07947a3404ce42939",
  "subnet-0792b862247faeaf8",
]

# route 53
private_hosted_zone_id = "Z09518463MDYVGIEDF9HV"

default_client_api_kms_alias = "alias/sbx-wise-us-vpc"

default_csp_kms_alias = "alias/csp-wise-us-vpc-csp"

# kinesis
txn_kinesis_name = "sbx-bbva-txn"

txn_kinesis_kms_alias = "alias/sbx-bbva-txn-s3"

txn_kinesis_region = "us-west-2"

ntf_kinesis_name = "sbx-bbva-ntf"

ntf_kinesis_kms_alias = "alias/sbx-bbva-ntf-s3"

ntf_kinesis_region = "us-west-2"

alloy_kinesis_name = "sbx-alloy-txn"

alloy_kinesis_kms_alias = "alias/sbx-alloy-txn-s3"

alloy_kinesis_region = "us-west-2"

# all services
cw_log_group_retention_in_days = 30

# payments service
payments_old_domain = "sbx-money-request.wise.us"

payments_domain = "sbx-payments.wise.us"

payments_image_tag = "build9"

payments_desired_container_count = 1

payments_min_container_count = 1

payments_max_container_count = 2

# Stripe Webhook Task
stripe_webhook_domain = "sbx-stripe-webhook.wise.us"

stripe_webhook_image_tag = "build9"

stripe_webhook_desired_container_count = 1

stripe_webhook_min_container_count = 1

stripe_webhook_max_container_count = 2

# BBVA Notification Tasks
bbva_env_name = "preprod"

bbva_notification_image_tag = "build9"

bbva_notification_desired_container_count = 1

bbva_notification_min_container_count = 1

bbva_notification_max_container_count = 2

# segment analytics
segment_analytics_image_tag = "build9"

segment_analytics_desired_container_count = 1

segment_analytics_min_container_count = 1

segment_analytics_max_container_count = 2

# shopify order
shopify_order_image_tag = "build9"

shopify_order_desired_container_count = 1

shopify_order_min_container_count = 1

shopify_order_max_container_count = 2

# App Notification Task
app_notification_image_tag = "build9"

app_notification_desired_container_count = 1

app_notification_min_container_count = 1

app_notification_max_container_count = 2

# batch account service
batch_account_image_tag = "build9"

batch_account_desired_container_count = 0

batch_account_min_container_count = 0

batch_account_max_container_count = 1

# batch analytics service
batch_analytics_image_tag = "build9"

batch_analytics_desired_container_count = 0

batch_analytics_min_container_count = 0

batch_analytics_max_container_count = 1

# batch monitor service
batch_monitor_image_tag = "build9"

batch_monitor_desired_container_count = 0

batch_monitor_min_container_count = 0

batch_monitor_max_container_count = 1

# batch transaction service
batch_transaction_image_tag = "build9"

batch_transaction_desired_container_count = 0

batch_transaction_min_container_count = 0

batch_transaction_max_container_count = 1

# batch monthly interest service
batch_monthly_interest_image_tag = "build9"

batch_monthly_interest_desired_container_count = 0

batch_monthly_interest_min_container_count = 0

batch_monthly_interest_max_container_count = 1

# merchant logos service
merchant_logo_domain = "sbx-merchant-logos-ecs.wise.us"

merchant_logo_image_tag = "build9"

merchant_logo_desired_container_count = 1

merchant_logo_min_container_count = 1

merchant_logo_max_container_count = 2

# hello sign service
hello_sign_domain = "sbx-hello-sign-webhook.wise.us"

hello_sign_image_tag = "build9"

hello_sign_desired_container_count = 1

hello_sign_min_container_count = 1

hello_sign_max_container_count = 2

# signature service
signature_image_tag = "build9"

signature_desired_container_count = 1

signature_min_container_count = 1

signature_max_container_count = 2

use_transaction_service = "true"

use_banking_service = "false"

use_invoice_service = "true"