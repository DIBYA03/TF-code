variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "csp"
}

variable "component" {
  default = "ecs"
}

variable "team" {
  default = "cloud-ops"
}

# SNS
variable "non_critical_sns_topic" {
  default = "arn:aws:sns:us-west-2:379379777492:wise-us-csp-noncritical-sns"
}

variable "critical_sns_topic" {
  default = "arn:aws:sns:us-west-2:379379777492:wise-us-csp-critical-sns"
}

# VPC
variable "vpc_id" {}

variable "vpc_cidr_block" {}

variable "shared_cidr_block" {
  default = "10.2.0.0/16"
}

variable "csp_rds_cidr_block" {}

variable "app_subnet_ids" {
  type = "list"
}

# Route 53
variable "route53_public_hosted_zone_id" {
  default = "Z3BUXRPXJI78KB"
}

variable "route53_private_hosted_zone_id" {
  default = "Z7Q2TZ46NCJYD"
}

variable "public_route53_account_profile" {
  default = "master-us-west-2-saml-roles-deployment"
}

# ach
variable "s3_ach_pull_list_config_object" {
  default = "config/ach_pull_whitelist.json"
}

# KMS
variable "default_kms_alias" {}

# clientAPI integrations
variable "default_client_api_env_kms_alias" {}

variable "core_db_cidr_blocks" {
  type = "list"
}

# ECS
variable "ecs_autoscaling_role" {
  default = "arn:aws:iam::379379777492:role/aws-service-role/ecs.application-autoscaling.amazonaws.com/AWSServiceRoleForApplicationAutoScaling_ECSService"
}

# CSP Frontend
variable "csp_frontend_add_monitoring" {
  default = true
}

variable "csp_frontend_domain" {}

variable "csp_frontend_image_tag" {}

variable "csp_frontend_name" {
  default = "csp-app-frontend"
}

variable "csp_frontend_container_port" {
  default = 80
}

variable "csp_frontend_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/csp-app-frontend"
}

variable "csp_frontend_cpu" {
  default = 256
}

variable "csp_frontend_mem" {
  default = 512
}

variable "csp_frontend_add_owasp10_waf" {
  default = true
}

variable "csp_frontend_desired_container_count" {
  default = 1
}

variable "csp_frontend_min_container_count" {
  default = 1
}

variable "csp_frontend_max_container_count" {
  default = 3
}

# document upload service
variable "csp_business_upload_add_monitoring" {
  default = true
}

variable "csp_business_upload_image_tag" {}

variable "csp_business_upload_desired_container_count" {}
variable "csp_business_upload_min_container_count" {}
variable "csp_business_upload_max_container_count" {}

variable "csp_business_upload_container_port" {
  default = 3000
}

variable "csp_business_upload_name" {
  default = "csp-api-document-upload"
}

variable "csp_business_upload_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-csp-document-upload"
}

variable "csp_business_upload_cpu" {
  default = 256
}

variable "csp_business_upload_mem" {
  default = 512
}

variable "csp_business_upload_add_owasp10_waf" {
  default = true
}

# all services
variable "cw_log_group_retention_in_days" {}

# review service
variable "csp_review_add_monitoring" {
  default = true
}

variable "csp_review_image_tag" {}

variable "csp_review_desired_container_count" {}
variable "csp_review_min_container_count" {}
variable "csp_review_max_container_count" {}

variable "csp_review_container_port" {
  default = 3000
}

variable "csp_review_name" {
  default = "csp-api-review"
}

variable "csp_review_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-csp-review"
}

variable "csp_review_cpu" {
  default = 256
}

variable "csp_review_mem" {
  default = 512
}

variable "csp_review_add_owasp10_waf" {
  default = true
}

# Batch business service
# batch timezone
variable "batch_default_timezone" {
  default = "America/Chicago"
}

# batch account
variable "batch_business_add_monitoring" {
  default = true
}

variable "batch_business_name" {
  default = "csp-batch-business"
}

variable "batch_business_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-csp-batch-business"
}

variable "batch_business_cpu" {
  default = 256
}

variable "batch_business_mem" {
  default = 512
}

variable "batch_business_image_tag" {}

variable "batch_business_cw_event_scheduled_expression" {}

variable "batch_business_desired_container_count" {}
variable "batch_business_min_container_count" {}
variable "batch_business_max_container_count" {}

variable "grpc_port" {
  default = 3001
}

variable "use_transaction_service" {}
variable "use_banking_service" {}

# batch account_closure

variable "batch_account_closure_add_monitoring" {
  default = true
}

variable "batch_account_closure_name" {
  default = "core-application-batch-account_closure"
}

variable "batch_account_closure_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-batch-account-closure"
}

variable "batch_account_closure_cpu" {
  default = 256
}

variable "batch_account_closure_mem" {
  default = 512
}

variable "batch_account_closure_image_tag" {}

variable "batch_account_closure_desired_container_count" {}
variable "batch_account_closure_min_container_count" {}
variable "batch_account_closure_max_container_count" {}

variable "account_closure_cw_event_scheduled_expression" {
  default = "cron(45 11 * * ? *)"
}
