aws_profile = "prod-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "prd"

ssm_environment = "prd"

environment_name = "prod"

vpc_id = "vpc-0dbc76a72cb520546"

vpc_cidr_block = "10.17.0.0/16"

app_subnet_ids = [
  "subnet-09259e808681d9c5a",
  "subnet-06cdb0c22471f3c1d",
  "subnet-0d75875254efbc0b1",
]

# SNS Topics
# Give permissinos to other accounts for getting SNS and subscribing
# Do not put the account number where the SNS topic lives
sns_allowed_subscribe_accounts = [
  "178124264531", # QA Prod
]

sns_non_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-critical-sns"

# KMS
default_kms_alias = "alias/prod-wise-us-vpc"

# BBVA
bbva_iam_role_env = "pro" # pre-production or production (values: pre, pro)

# BBVA SNS Connector Task
bbva_sns_connector_image_tag = "build5"

bbva_sns_connector_desired_container_count = 2

bbva_sns_connector_min_container_count = 2

bbva_sns_connector_max_container_count = 5
