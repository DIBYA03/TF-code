aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "sbx"

environment_name = "sbx"

# VPC
vpc_id = "vpc-0a9971002df8d25bb"

vpc_cidr_block = "10.4.0.0/16"

csp_rds_vpc_cidr_block = "10.4.0.0/16"

app_subnet_ids = [
  "subnet-0aa0906e487fd35de",
  "subnet-0f00f5cce55e9fdf3",
  "subnet-086bc5764359f56e6",
]

# Route53
route53_domain_name = "sbx-csp.internal.wise.us"

# SNS
sns_non_critical_topic = "arn:aws:sns:us-west-2:345501831767:wise-us-sbx-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:345501831767:wise-us-sbx-critical-sns"

# KMS
default_kms_alias = "alias/csp-wise-us-vpc-csp"

# clientAPI integrations
core_db_cidr_blocks = [
  "10.3.0.0/16",
]

# API Gateway
api_gw_endpoint_configuration = "PRIVATE"

api_gw_server_description = "Sbx CSP API Server"

# cognito
cognito_domain_name = "sbx-csp-auth-wise"
