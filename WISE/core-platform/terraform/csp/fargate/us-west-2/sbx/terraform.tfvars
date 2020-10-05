aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "sbx"

environment_name = "sbx"

# VPC
vpc_id = "vpc-0a9971002df8d25bb"

vpc_cidr_block = "10.4.0.0/16"

csp_rds_cidr_block = "10.4.0.0/16"

app_subnet_ids = [
  "subnet-0aa0906e487fd35de",
  "subnet-0f00f5cce55e9fdf3",
  "subnet-086bc5764359f56e6",
]

# KMS
default_kms_alias = "alias/csp-wise-us-vpc-csp"

# clientAPI integrations
default_client_api_env_kms_alias = "alias/sbx-wise-us-vpc"

core_db_cidr_blocks = [
  "10.3.0.0/16",
]

# csp frontend service
csp_frontend_domain = "sbx-csp.internal.wise.us"

csp_frontend_image_tag = "sbx-build4"

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

use_banking_service = "false"


# batch account closure service

batch_account_closure_image_tag = "build9"

batch_account_closure_desired_container_count = 0

batch_account_closure_min_container_count = 0

batch_account_closure_max_container_count = 1

