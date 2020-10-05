aws_profile = "dev-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "dev1"

environment_name = "dev1"

# VPC
vpc_id = "vpc-02d22afaa5a6a4d8a"

vpc_cidr_block = "10.24.0.0/16"

app_subnet_ids = [
  "subnet-0e1481139ef24777a",
  "subnet-045b3c542f4822bb0",
  "subnet-0ebe63bb6037423d0",
]

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:152334605517:wise-us-dev1-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:152334605517:wise-us-dev1-noncritical-sns"

# KMS
default_kms_alias = "alias/dev1-wise-us-vpc"

# Lambda
lambda_timeout = 60

# API Gateway
api_gw_server_description = "Development Server"

api_gw_domain_name = "dev1-aws-apigw-client-api.wise.us"

# BBVA
bbva_wise_profile = "master-us-west-2-saml-roles-deployment"

bbva_notifications_env = "ppd"

bbva_sqs_environment = "ppd"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"