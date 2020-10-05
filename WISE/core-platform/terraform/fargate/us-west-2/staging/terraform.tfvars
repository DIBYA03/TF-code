aws_profile = "master-us-west-2-saml-roles-deployment"
public_route53_account_profile = "master-us-west-2-saml-roles-deployment"
public_hosted_zone_id = "Z3BUXRPXJI78KB"


aws_region = "us-west-2"

environment = "stg"

environment_name = "staging"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

# VPC
vpc_id = "vpc-043a00e2bba3b0dd8"

vpc_cidr_block = "10.3.0.0/16"

csp_vpc_cidr_block = "10.4.0.0/16"

client_api_rds_vpc_cidr_block = "10.3.0.0/16"

public_subnet_ids = [
  "subnet-05eab8142e3542619",
  "subnet-03bb66a24fdc36358",
  "subnet-0bf16b4e235d2cdd1",
]

app_subnet_ids = [
  "subnet-0e932567d087fdd68",
  "subnet-0b3b84440068c1bf7",
  "subnet-0322b5e9cf73e7eb0",
]

# route 53
private_hosted_zone_id = "Z1TXQSQLQVDHXJ"

default_client_api_kms_alias = "alias/staging-wise-us-vpc"

default_csp_kms_alias = "alias/csp-wise-us-vpc"

# kinesis
txn_kinesis_name = "stg-firehose-txn-kinesis"

txn_kinesis_kms_alias = "alias/stg-firehose-txn-s3"

txn_kinesis_region = "us-west-2"

ntf_kinesis_name = "stg-logging-kinesis-firehose"

ntf_kinesis_kms_alias = "alias/stg-kinesis-wise-us"

ntf_kinesis_region = "us-west-2"

alloy_kinesis_name = "stg-alloy-txn"

alloy_kinesis_kms_alias = "alias/stg-alloy-txn-s3"

alloy_kinesis_region = "us-west-2"

# all services
cw_log_group_retention_in_days = 30

# payments service
payments_old_domain = "staging-money-request.wise.us"

payments_domain = "staging-payments.wise.us"

payments_image_tag = "build4"

payments_desired_container_count = 1

payments_min_container_count = 1

payments_max_container_count = 2

# Stripe Webhook Task
stripe_webhook_domain = "staging-stripe-webhook.wise.us"

stripe_webhook_image_tag = "build4"

stripe_webhook_desired_container_count = 1

stripe_webhook_min_container_count = 1

stripe_webhook_max_container_count = 2

# BBVA Notification Tasks
bbva_env_name = "preprod"

bbva_notification_image_tag = "build4"

bbva_notification_desired_container_count = 1

bbva_notification_min_container_count = 1

bbva_notification_max_container_count = 2

# segment analytics
segment_analytics_image_tag = "build4"

segment_analytics_desired_container_count = 1

segment_analytics_min_container_count = 1

segment_analytics_max_container_count = 2

# shopify order
shopify_order_image_tag = "build4"

shopify_order_desired_container_count = 1

shopify_order_min_container_count = 1

shopify_order_max_container_count = 2

# App Notification Task
app_notification_image_tag = "build4"

app_notification_desired_container_count = 1

app_notification_min_container_count = 1

app_notification_max_container_count = 2

# batch account service
batch_account_image_tag = "build4"

batch_account_desired_container_count = 0

batch_account_min_container_count = 0

batch_account_max_container_count = 1

# batch analytics service
batch_analytics_image_tag = "build4"

batch_analytics_desired_container_count = 0

batch_analytics_min_container_count = 0

batch_analytics_max_container_count = 1

# batch monitor service
batch_monitor_image_tag = "build4"

batch_monitor_desired_container_count = 0

batch_monitor_min_container_count = 0

batch_monitor_max_container_count = 1

# batch transaction service
batch_transaction_image_tag = "build4"

batch_transaction_desired_container_count = 0

batch_transaction_min_container_count = 0

batch_transaction_max_container_count = 1

# batch monthly interest service
batch_monthly_interest_image_tag = "build4"

batch_monthly_interest_desired_container_count = 0

batch_monthly_interest_min_container_count = 0

batch_monthly_interest_max_container_count = 1

# merchant logos service
merchant_logo_domain = "staging-merchant-logos-ecs.wise.us"

merchant_logo_image_tag = "build4"

merchant_logo_desired_container_count = 1

merchant_logo_min_container_count = 1

merchant_logo_max_container_count = 2

# hello sign service
hello_sign_domain = "staging-hello-sign-webhook.wise.us"

hello_sign_image_tag = "build4"

hello_sign_desired_container_count = 1

hello_sign_min_container_count = 1

hello_sign_max_container_count = 2

# signature service
signature_image_tag = "build4"

signature_desired_container_count = 1

signature_min_container_count = 1

signature_max_container_count = 2

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"
