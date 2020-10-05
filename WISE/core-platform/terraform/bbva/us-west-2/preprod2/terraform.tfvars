aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "ppd"

ssm_environment = "sbx"

environment_name = "preprod"

vpc_id = "vpc-06a457dbfa4e8ed8d"

vpc_cidr_block = "10.3.0.0/16"

app_subnet_ids = [
  "subnet-0981342e0da7bbe7d",
  "subnet-07947a3404ce42939",
  "subnet-0792b862247faeaf8",
]

# SNS Topics
sns_non_critical_topic = "arn:aws:sns:us-west-2:345501831767:wise-us-sbx-noncritical-sns"

sns_critical_topic = "arn:aws:sns:us-west-2:345501831767:wise-us-sbx-critical-sns"

# KMS
default_kms_alias = "alias/sbx-wise-us-vpc"

# BBVA
bbva_iam_role_env = "pre" # pre-production or production (values: pre, pro)

# BBVA SNS Connector Task
bbva_sns_connector_image_tag = "build1"

bbva_sns_connector_desired_container_count = 2

bbva_sns_connector_min_container_count = 2

bbva_sns_connector_max_container_count = 5
