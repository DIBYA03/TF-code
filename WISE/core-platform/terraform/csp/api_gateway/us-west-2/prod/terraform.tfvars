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

# KMS
default_kms_alias = "alias/csp-prod-wise-us-vpc"

# SNS
sns_non_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-csp-prod-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-csp-prod-critical-sns"

# Route53
route53_domain_name = "csp.internal.wise.us"

# clientAPI integrations
core_db_cidr_blocks = [
  "10.17.0.0/16",
]

# API Gateway
api_gw_endpoint_configuration = "PRIVATE"

api_gw_server_description = "Prod CSP API Server"

# cognito
cognito_domain_name = "csp-auth-wise"
