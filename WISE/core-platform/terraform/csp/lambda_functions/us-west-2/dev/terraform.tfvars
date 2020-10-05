aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "dev"

environment_name = "dev"

# VPC
vpc_id = "vpc-0aef975032e418d38"

vpc_cidr_block = "10.4.0.0/16"

csp_rds_vpc_cidr_block = "10.4.0.0/16"

app_subnet_ids = [
  "subnet-0f20d4c37f4242b7b",
  "subnet-051e08c06df7de0f0",
  "subnet-07bd6a1a5586635b0",
]

# SNS
sns_non_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-csp-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-csp-noncritical-sns"

# KMS
default_kms_alias = "alias/csp-wise-us-vpc"

# clientAPI integrations
core_db_cidr_blocks = [
  "10.1.0.0/16",
]

# Route53
domain_name = "dev-csp-api.internal.wise.us"

# kinesis
txn_kinesis_name = "dev-firehose-txn-kinesis"

txn_kinesis_region = "us-west-2"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"