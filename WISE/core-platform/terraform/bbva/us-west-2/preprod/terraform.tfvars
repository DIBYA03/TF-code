aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "ppd"

ssm_environment = "stg"

environment_name = "preprod"

vpc_id = "vpc-043a00e2bba3b0dd8"

vpc_cidr_block = "10.3.0.0/16"

app_subnet_ids = [
  "subnet-0e932567d087fdd68",
  "subnet-0b3b84440068c1bf7",
  "subnet-0322b5e9cf73e7eb0",
]

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

# KMS
default_kms_alias = "alias/staging-wise-us-vpc"

# BBVA
bbva_iam_role_env = "pre" # pre-production or production (values: pre, pro)

# BBVA SNS Connector Task
bbva_sns_connector_image_tag = "build1"

bbva_sns_connector_desired_container_count = 2

bbva_sns_connector_min_container_count = 2

bbva_sns_connector_max_container_count = 5
