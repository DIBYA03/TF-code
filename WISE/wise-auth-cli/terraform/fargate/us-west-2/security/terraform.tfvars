aws_profile = "security-us-west-2-saml-roles-admin"

aws_region = "us-west-2"

environment = "sec"

environment_name = "security"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:842431896317:wise-us-security-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:842431896317:wise-us-security-critical-sns"

# VPC
vpc_id = "vpc-07fc7f07a5f63dd99"

vpc_cidr_block = "10.23.0.0/16"

app_subnet_ids = [
  "subnet-00469bf0cdd0ddbe7",
  "subnet-06e9158b401c91ed4",
  "subnet-0e45f33374b7f9870",
]

# kms
default_kms_alias = "alias/security-wise-us-vpc"

# cloudwatch
cw_log_group_retention_in_days = 365

# route53
public_route53_hosted_zone = "Z3BUXRPXJI78KB"

private_route53_hosted_zone = "Z7Q2TZ46NCJYD"

# aws_vpn_auth service
aws_vpn_auth_domain = "aws.us-west-2.internal.wise.us"

aws_vpn_auth_allowed_account_ids = [
  "arn:aws:iam::379379777492:root",
]

aws_vpn_auth_image_tag = "latest"

aws_vpn_auth_desired_container_count = 1

aws_vpn_auth_min_container_count = 1

aws_vpn_auth_max_container_count = 2

# wise aws auth cli
aws_vpn_auth_wise_image_tag = "latest"

# endpoint service
endpoint_service_vpc_id = "vpc-0426719bb5133d7a0"

endpoint_service_subnet_ids = [
  "subnet-056282f599d63b8af",
  "subnet-006d54bf1bdc1e424",
]

endpoint_service_allowed_cidr_blocks = [
  "10.2.0.0/16",
]
