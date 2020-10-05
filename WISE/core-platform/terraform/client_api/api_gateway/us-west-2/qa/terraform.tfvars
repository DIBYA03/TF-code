aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "qa"

environment_name = "qa"

# VPC
vpc_id = "vpc-0e844a5b87a3cfce5"

vpc_cidr_block = "10.1.0.0/16"

app_subnet_ids = [
  "subnet-0e1481139ef24777a",
  "subnet-045b3c542f4822bb0",
  "subnet-0ebe63bb6037423d0",
]

# Route53
route53_domain_name = "qa-aws-apigw-client-api.internal.wise.us"

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-dev-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-dev-noncritical-sns"

# KMS
default_kms_alias = "alias/dev-wise-us-vpc"

# Lambda
lambda_timeout = 60

# API Gateway
api_gw_server_description = "QA Server"

# BBVA
bbva_wise_profile = "master-us-west-2-saml-roles-deployment"

bbva_notifications_env = "ppd"

bbva_sqs_environment = "ppd"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"
