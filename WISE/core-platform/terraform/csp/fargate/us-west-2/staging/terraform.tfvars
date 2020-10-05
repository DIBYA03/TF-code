aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "stg"

environment_name = "staging"

# VPC
vpc_id = "vpc-0aef975032e418d38"

vpc_cidr_block = "10.4.0.0/16"

csp_rds_cidr_block = "10.4.0.0/16"

app_subnet_ids = [
  "subnet-0f20d4c37f4242b7b",
  "subnet-051e08c06df7de0f0",
  "subnet-07bd6a1a5586635b0",
]

# KMS
default_kms_alias = "alias/csp-wise-us-vpc"

# clientAPI integrations
default_client_api_env_kms_alias = "alias/staging-wise-us-vpc"

core_db_cidr_blocks = [
  "10.3.0.0/16",
]

# csp frontend service
csp_frontend_domain = "staging-csp.internal.wise.us"

csp_frontend_image_tag = "staging-build4"

csp_frontend_desired_container_count = 1

csp_frontend_min_container_count = 1

ccsp_frontend_max_container_count = 6

# document upload service
csp_business_upload_image_tag = "build4"

csp_business_upload_desired_container_count = 1

csp_business_upload_min_container_count = 1

csp_business_upload_max_container_count = 6

# review service
csp_review_image_tag = "build4"

csp_review_desired_container_count = 1

csp_review_min_container_count = 1

csp_review_max_container_count = 6

# batch business service
batch_business_cw_event_scheduled_expression = "cron(0 1 1 1 ? 1970)" # never really run

batch_business_image_tag = "build4"

batch_business_desired_container_count = 0

batch_business_min_container_count = 0

batch_business_max_container_count = 1

use_transaction_service = "true"

use_banking_service = "true"

# batch account closure service

batch_account_closure_image_tag = "build4"

batch_account_closure_desired_container_count = 0

batch_account_closure_min_container_count = 0

batch_account_closure_max_container_count = 1

