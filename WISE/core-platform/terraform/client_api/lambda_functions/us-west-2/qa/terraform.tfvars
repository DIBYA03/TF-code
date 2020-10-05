aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "qa"

environment_name = "qa"

vpc_id = "vpc-0e844a5b87a3cfce5"

vpc_cidr_block = "10.1.0.0/16"

app_subnet_ids = [
  "subnet-0e1481139ef24777a",
  "subnet-045b3c542f4822bb0",
  "subnet-0ebe63bb6037423d0",
]

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-dev-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-dev-noncritical-sns"

# CSP integrations
csp_environment = "qa"

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
txn_kinesis_name = "qa-firehose-txn-kinesis"

txn_kinesis_region = "us-west-2"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"