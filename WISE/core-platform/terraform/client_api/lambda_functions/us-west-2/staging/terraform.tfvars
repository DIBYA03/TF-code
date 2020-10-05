aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "stg"

environment_name = "staging"

vpc_id = "vpc-043a00e2bba3b0dd8"

vpc_cidr_block = "10.3.0.0/16"

app_subnet_ids = [
  "subnet-0e932567d087fdd68",
  "subnet-0b3b84440068c1bf7",
  "subnet-0322b5e9cf73e7eb0",
]

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

# CSP integrations
csp_environment = "stg"

csp_kms_alias = "alias/csp-wise-us-vpc"

# API Gateway
api_gw_5XX_error_alarm_non_critical_threshold = 5

api_gw_5XX_error_alarm_critical_threshold = 20

api_gw_4XX_error_alarm_non_critical_threshold = 20

api_gw_4XX_error_alarm_critical_threshold = 50

api_gw_latency_alarm_threshold = 3000 # 3 seconds

# Lambda
lambda_timeout = 60

# kinesis
txn_kinesis_name = "stg-firehose-txn-kinesis"

txn_kinesis_region = "us-west-2"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"