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

# Route53
route53_domain_name = "aws-apigw-client-api.internal.wise.us"

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-critical-sns"

# KMS
default_kms_alias = "alias/prod-wise-us-vpc"

# Lambda
lambda_timeout = 60

# API Gateway
api_gw_server_description = "Prod Server"

# BBVA Notifications
bbva_wise_profile = "prod-us-west-2-saml-roles-deployment"

bbva_sqs_environment = "prd"

bbva_notifications_env = "prd"

use_transaction_service = "true"

use_banking_service = "true"

use_invoice_service = "true"