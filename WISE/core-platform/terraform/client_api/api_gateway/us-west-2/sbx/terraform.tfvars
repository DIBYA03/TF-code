aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "sbx"

environment_name = "sbx"

# VPC
vpc_id = "vpc-06a457dbfa4e8ed8d"

vpc_cidr_block = "10.3.0.0/16"

app_subnet_ids = [
  "subnet-0981342e0da7bbe7d",
  "subnet-07947a3404ce42939",
  "subnet-0792b862247faeaf8",
]

# Route53
route53_domain_name = "aws-apigw-client-api.internal.sbx.wise.us"

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:345501831767:wise-us-sbx-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:345501831767:wise-us-sbx-critical-sns"

# KMS
default_kms_alias = "alias/sbx-wise-us-vpc"

# Lambda
lambda_timeout = 60

# API Gateway
api_gw_server_description = "Staging Server"

# BBVA
bbva_wise_profile = "wiseus"

bbva_notifications_env = "ppd"

bbva_sqs_environment = "ppd"

use_transaction_service = "true"

use_banking_service = "false"

use_invoice_service = "true"