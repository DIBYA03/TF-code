aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "sbx"

environment_name = "sbx"

vpc_id = "vpc-06a457dbfa4e8ed8d"

vpc_cidr_block = "10.3.0.0/16"

app_subnet_ids = [
  "subnet-0981342e0da7bbe7d",
  "subnet-07947a3404ce42939", 
  "subnet-0792b862247faeaf8",
]

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

# CSP integrations
csp_environment = "sbx"

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
txn_kinesis_name = "sbx-firehose-txn-kinesis"

txn_kinesis_region = "us-west-2"

use_transaction_service = "true"

use_banking_service = "false"

use_invoice_service = "true"