aws_profile = "prod-us-west-2-saml-roles-deployment"
public_route53_account_profile = "master-us-west-2-saml-roles-deployment"
public_hosted_zone_id = "Z3BUXRPXJI78KB"

aws_region = "us-west-2"

environment = "prd"

environment_name = "prod"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-critical-sns"

# VPC
vpc_id = "vpc-0dbc76a72cb520546"

vpc_cidr_block = "10.17.0.0/16"

csp_vpc_cidr_block = "10.20.0.0/16"

client_api_rds_vpc_cidr_block = "10.17.0.0/16"

public_subnet_ids = [
  "subnet-014188901727d1077",
  "subnet-0ae11e622379832a4",
  "subnet-01f7544c17ed9f3eb",
]

app_subnet_ids = [
  "subnet-09259e808681d9c5a",
  "subnet-06cdb0c22471f3c1d",
  "subnet-0d75875254efbc0b1",
]

default_client_api_kms_alias = "alias/prod-wise-us-vpc"

default_csp_kms_alias = "alias/csp-prod-wise-us-vpc"

# route 53
private_hosted_zone_id = "ZEQGF56E9VNY4"

# kinesis
txn_kinesis_name = "prd-bbva-txn"

txn_kinesis_kms_alias = "alias/prd-bbva-txn-s3"

txn_kinesis_region = "us-west-2"

ntf_kinesis_name = "prd-bbva-ntf"

ntf_kinesis_kms_alias = "alias/prd-bbva-ntf-s3"

ntf_kinesis_region = "us-west-2"

alloy_kinesis_name = "prd-alloy-txn"

alloy_kinesis_kms_alias = "alias/prd-alloy-txn-s3"

alloy_kinesis_region = "us-west-2"

# all services
cw_log_group_retention_in_days = 365

# payments service
payments_old_domain = "money-request.wise.us"

payments_domain = "payments.wise.us"

payments_image_tag = "build7"

payments_desired_container_count = 2

payments_min_container_count = 2

payments_max_container_count = 6

# Stripe Webhook Task
stripe_webhook_domain = "stripe-webhook.wise.us"

stripe_webhook_image_tag = "build7"

stripe_webhook_desired_container_count = 2

stripe_webhook_min_container_count = 2

stripe_webhook_max_container_count = 6

# BBVA Notification Tasks
bbva_env_name = "prod"

bbva_notification_image_tag = "build7"

bbva_notification_desired_container_count = 2

bbva_notification_min_container_count = 2

bbva_notification_max_container_count = 6

# segment analytics
segment_analytics_image_tag = "build7"

segment_analytics_desired_container_count = 2

segment_analytics_min_container_count = 2

segment_analytics_max_container_count = 6

# shopify order
shopify_order_image_tag = "build7"

shopify_order_desired_container_count = 2

shopify_order_min_container_count = 2

shopify_order_max_container_count = 6

# App Notification Task
app_notification_image_tag = "build7"

app_notification_desired_container_count = 2

app_notification_min_container_count = 2

app_notification_max_container_count = 6

# batch account service
batch_account_image_tag = "build7"

batch_account_desired_container_count = 0

batch_account_min_container_count = 0

batch_account_max_container_count = 1

# batch analytics service
batch_analytics_image_tag = "build7"

batch_analytics_desired_container_count = 0

batch_analytics_min_container_count = 0

batch_analytics_max_container_count = 1

# batch monitor service
batch_monitor_image_tag = "build7"

batch_monitor_desired_container_count = 0

batch_monitor_min_container_count = 0

batch_monitor_max_container_count = 1

# batch transaction service
batch_transaction_image_tag = "build7"

batch_transaction_desired_container_count = 0

batch_transaction_min_container_count = 0

batch_transaction_max_container_count = 1

# batch monthly interest service
batch_monthly_interest_image_tag = "build7"

batch_monthly_interest_desired_container_count = 0

batch_monthly_interest_min_container_count = 0

batch_monthly_interest_max_container_count = 1

# merchant logos service
merchant_logo_domain = "merchant-logos-ecs.wise.us"

merchant_logo_image_tag = "build7"

merchant_logo_desired_container_count = 2

merchant_logo_min_container_count = 2

merchant_logo_max_container_count = 6

# hello sign service
hello_sign_domain = "hello-sign-webhook.wise.us"

hello_sign_image_tag = "build7"

hello_sign_desired_container_count = 2

hello_sign_min_container_count = 2

hello_sign_max_container_count = 6

# signature service
signature_image_tag = "build7"

signature_desired_container_count = 2

signature_min_container_count = 2

signature_max_container_count = 6

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"