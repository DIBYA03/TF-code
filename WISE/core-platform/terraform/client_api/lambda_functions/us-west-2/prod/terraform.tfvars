aws_profile = "prod-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "prd"

environment_name = "prod"

# VPC
vpc_id = "vpc-0dbc76a72cb520546"

vpc_cidr_block = "10.17.0.0/16"

app_subnet_ids = [
  "subnet-09259e808681d9c5a",
  "subnet-06cdb0c22471f3c1d",
  "subnet-0d75875254efbc0b1",
]

# CSP integrations
csp_environment = "prd"

csp_kms_alias = "alias/csp-prod-wise-us-vpc"

## BBVA integrations
# prd-bbva-ntf
# prd-bbva-txn

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-critical-sns"

# API Gateway
api_gw_5XX_error_alarm_non_critical_threshold = 5

api_gw_5XX_error_alarm_critical_threshold = 20

api_gw_4XX_error_alarm_non_critical_threshold = 20

api_gw_4XX_error_alarm_critical_threshold = 50

api_gw_latency_alarm_threshold = 3000 # 3 seconds

# Lambda
lambda_timeout = 60

# kinesis
txn_kinesis_name = "prd-bbva-txn"

txn_kinesis_region = "us-west-2"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"