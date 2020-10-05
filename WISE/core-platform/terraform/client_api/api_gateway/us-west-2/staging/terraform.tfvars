aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "stg"

environment_name = "staging"

# VPC
vpc_id = "vpc-043a00e2bba3b0dd8"

vpc_cidr_block = "10.3.0.0/16"

app_subnet_ids = [
  "subnet-0e932567d087fdd68",
  "subnet-0b3b84440068c1bf7",
  "subnet-0322b5e9cf73e7eb0",
]

# Route53
route53_domain_name = "stg-aws-apigw-client-api.internal.wise.us"

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

# KMS
default_kms_alias = "alias/staging-wise-us-vpc"

# Lambda
lambda_timeout = 60

# API Gateway
api_gw_server_description = "Staging Server"

# BBVA
bbva_wise_profile = "master-us-west-2-saml-roles-deployment"

bbva_notifications_env = "ppd"

bbva_sqs_environment = "ppd"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"