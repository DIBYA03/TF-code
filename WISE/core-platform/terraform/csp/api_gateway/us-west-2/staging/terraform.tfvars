aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "stg"

environment_name = "staging"

# VPC
vpc_id = "vpc-0aef975032e418d38"

vpc_cidr_block = "10.4.0.0/16"

csp_rds_vpc_cidr_block = "10.4.0.0/16"

app_subnet_ids = [
  "subnet-0f20d4c37f4242b7b",
  "subnet-051e08c06df7de0f0",
  "subnet-07bd6a1a5586635b0",
]

# Route53
route53_domain_name = "staging-csp.internal.wise.us"

# SNS
sns_non_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-csp-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-csp-noncritical-sns"

# KMS
default_kms_alias = "alias/csp-wise-us-vpc"

# clientAPI integrations
core_db_cidr_blocks = [
  "10.3.0.0/16",
]

# API Gateway
api_gw_endpoint_configuration = "PRIVATE"

api_gw_server_description = "Staging CSP API Server"

# cognito
cognito_domain_name = "staging-csp-auth-wise"
