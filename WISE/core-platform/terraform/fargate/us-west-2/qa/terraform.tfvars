aws_profile = "master-us-west-2-saml-roles-deployment"
public_route53_account_profile = "master-us-west-2-saml-roles-deployment"
public_hosted_zone_id = "Z3BUXRPXJI78KB"

aws_region = "us-west-2"

environment = "qa"

environment_name = "qa"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-dev-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-dev-noncritical-sns"

# VPC
vpc_id = "vpc-0e844a5b87a3cfce5"

vpc_cidr_block = "10.1.0.0/16"

csp_vpc_cidr_block = "10.4.0.0/16"

client_api_rds_vpc_cidr_block = "10.1.0.0/16"

public_subnet_ids = [
  "subnet-0333ce9291d82fcc6",
  "subnet-09aaf723dbdc2edf5",
  "subnet-0fc14e8e6df7b0bcf",
]

app_subnet_ids = [
  "subnet-045b3c542f4822bb0",
  "subnet-0e1481139ef24777a",
  "subnet-0ebe63bb6037423d0",
]

# route 53
private_hosted_zone_id = "Z35OJJKVJ3I2IE"

default_client_api_kms_alias = "alias/dev-wise-us-vpc"

default_csp_kms_alias = "alias/csp-wise-us-vpc"

# kinesis
txn_kinesis_name = "qa-firehose-txn-kinesis"

txn_kinesis_kms_alias = "alias/qa-firehose-txn-s3"

txn_kinesis_region = "us-west-2"

ntf_kinesis_name = "qa-logging-kinesis-firehose"

ntf_kinesis_kms_alias = "alias/qa-kinesis-wise-us"

ntf_kinesis_region = "us-west-2"

alloy_kinesis_name = "qa-alloy-txn"

alloy_kinesis_kms_alias = "alias/qa-alloy-txn-s3"

alloy_kinesis_region = "us-west-2"

# all services
cw_log_group_retention_in_days = 30

# payments service
payments_old_domain = "qa-money-request.wise.us"

payments_domain = "qa-payments.wise.us"

payments_image_tag = "build2"

payments_desired_container_count = 1

payments_min_container_count = 1

payments_max_container_count = 2

# Stripe Webhook Task
stripe_webhook_domain = "qa-stripe-webhook.wise.us"

stripe_webhook_image_tag = "build2"

stripe_webhook_desired_container_count = 1

stripe_webhook_min_container_count = 1

stripe_webhook_max_container_count = 2

# BBVA Notification Tasks
bbva_env_name = "preprod"

bbva_notification_image_tag = "build2"

bbva_notification_desired_container_count = 1

bbva_notification_min_container_count = 1

bbva_notification_max_container_count = 3

# segment analytics
segment_analytics_image_tag = "build2"

segment_analytics_desired_container_count = 1

segment_analytics_min_container_count = 1

segment_analytics_max_container_count = 2

# shopify order
shopify_order_image_tag = "build2"

shopify_order_desired_container_count = 1

shopify_order_min_container_count = 1

shopify_order_max_container_count = 2

# App Notification Task
app_notification_image_tag = "build2"

app_notification_desired_container_count = 1

app_notification_min_container_count = 1

app_notification_max_container_count = 2

# batch account service
batch_account_image_tag = "build2"

batch_account_desired_container_count = 0

batch_account_min_container_count = 0

batch_account_max_container_count = 1

# batch analytics service
batch_analytics_image_tag = "build2"

batch_analytics_desired_container_count = 0

batch_analytics_min_container_count = 0

batch_analytics_max_container_count = 1

# batch monitor service
batch_monitor_image_tag = "build2"

batch_monitor_desired_container_count = 0

batch_monitor_min_container_count = 0

batch_monitor_max_container_count = 1

# batch transaction service
batch_transaction_image_tag = "build2"

batch_transaction_desired_container_count = 0

batch_transaction_min_container_count = 0

batch_transaction_max_container_count = 1

# batch monthly interest service
batch_monthly_interest_image_tag = "build2"

batch_monthly_interest_desired_container_count = 0

batch_monthly_interest_min_container_count = 0

batch_monthly_interest_max_container_count = 1

# merchant logos service
merchant_logo_domain = "qa-merchant-logos-ecs.wise.us"

merchant_logo_image_tag = "build2"

merchant_logo_desired_container_count = 1

merchant_logo_min_container_count = 1

merchant_logo_max_container_count = 2

# hello sign service
hello_sign_domain = "qa-hello-sign-webhook.wise.us"

hello_sign_image_tag = "build2"

hello_sign_desired_container_count = 1

hello_sign_min_container_count = 1

hello_sign_max_container_count = 2

# signature service
signature_image_tag = "build2"

signature_desired_container_count = 1

signature_min_container_count = 1

signature_max_container_count = 2

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"
