aws_profile = "prod-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "prd"

environment_name = "prod"

# VPC
vpc_id = "vpc-0a3d5ba3cf7256441"

vpc_cidr_block = "10.20.0.0/16"

csp_rds_vpc_cidr_block = "10.20.0.0/16"

app_subnet_ids = [
  "subnet-0ecaac978d90a9d35",
  "subnet-0746233f88b758a7d",
  "subnet-0c4a66fe32ca4033f",
]

# SNS
sns_non_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-csp-prod-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-csp-prod-critical-sns"

# KMS
default_kms_alias = "alias/csp-prod-wise-us-vpc"

# clientAPI integrations
core_db_cidr_blocks = [
  "10.17.0.0/16",
]

# Route53
domain_name = "csp-api.internal.wise.us"

# kinesis
txn_kinesis_name = "prd-bbva-txn"

txn_kinesis_region = "us-west-2"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"