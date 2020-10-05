aws_profile = "dev-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "ppd"

ssm_environment = "dev1"

environment_name = "preprod"

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

# BBVA
bbva_iam_role_env = "pre" # pre-production or production (values: pre, pro)

# BBVA SNS Connector Task
bbva_sns_connector_image_tag = "build1"

bbva_sns_connector_desired_container_count = 2

bbva_sns_connector_min_container_count = 2

bbva_sns_connector_max_container_count = 3

# this needs at least one to not error the sns policy
sns_allowed_subscribe_accounts = [
  "152334605517",
]
