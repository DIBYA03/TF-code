aws_profile = "prod-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "prd"

environment_name = "prod"

# VPC
vpc_id = "vpc-0a3d5ba3cf7256441"

vpc_cidr_block = "10.20.0.0/16"

csp_rds_cidr_block = "10.20.0.0/16"

app_subnet_ids = [
  "subnet-0ecaac978d90a9d35",
  "subnet-0746233f88b758a7d",
  "subnet-0c4a66fe32ca4033f",
]

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-csp-prod-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-csp-prod-critical-sns"

default_kms_alias = "alias/csp-prod-wise-us-vpc"

# clientAPI integrations
default_client_api_env_kms_alias = "alias/prod-wise-us-vpc"

core_db_cidr_blocks = [
  "10.17.0.0/16",
]

# csp frontend service
csp_frontend_domain = "csp.internal.wise.us"

csp_frontend_image_tag = "prod-build7"

csp_frontend_desired_container_count = 2

csp_frontend_min_container_count = 2

ccsp_frontend_max_container_count = 6

# document upload service
csp_business_upload_image_tag = "build7"

csp_business_upload_desired_container_count = 2

csp_business_upload_min_container_count = 2

csp_business_upload_max_container_count = 6

# review service
csp_review_image_tag = "build7"

csp_review_desired_container_count = 2

csp_review_min_container_count = 2

csp_review_max_container_count = 6

# batch business service
batch_business_cw_event_scheduled_expression = "cron(0/10 14-02 ? * 2-6 *)"

batch_business_image_tag = "build7"

batch_business_desired_container_count = 0

batch_business_min_container_count = 0

batch_business_max_container_count = 1

use_transaction_service = "true"

use_banking_service = "true"

# batch account closure service

batch_account_closure_image_tag = "build7"

batch_account_closure_desired_container_count = 0

batch_account_closure_min_container_count = 0

batch_account_closure_max_container_count = 1
