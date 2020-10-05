variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "ssm_environment" {}
variable "environment_name" {}

variable "application" {
  default = "bbva"
}

variable "component" {
  default = "ntf"
}

variable "team" {
  default = "cloud-ops"
}

variable "vpc_id" {}
variable "vpc_cidr_block" {}

# variable "public_subnet_ids" {
#   type = "list"
# }

variable "app_subnet_ids" {
  type = "list"
}

# SNS
variable "sns_non_critical_topic" {}

variable "sns_critical_topic" {}

variable "sns_allowed_subscribe_accounts" {
  type    = "list"
  default = []
}

# KMS
variable "default_kms_alias" {}

# SQS
variable "bbva_sqs_visibility_timeout_seconds" {
  default = 30
}

variable "bbva_sqs_message_retention_seconds" {
  default = 604800
}

variable "bbva_sqs_dl_message_retention_seconds" {
  default = 1209600
}

variable "bbva_sqs_content_based_deduplication" {
  default = false
}

variable "bbva_sqs_dl_content_based_deduplication" {
  default = false
}

variable "bbva_sqs_max_message_size" {
  default = 262144
}

variable "bbva_sqs_delay_seconds" {
  default = 0
}

variable "bbva_sqs_receive_wait_time_seconds" {
  default = 0
}

variable "bbva_sqs_fifo_queue" {
  default = false
}

variable "bbva_sqs_kms_data_key_reuse_period_seconds" {
  default = 300
}

# BBVA
variable "bbva_iam_role_env" {}

# Subscribe scripts
variable "bbva_subscribe_account_transactions_script" {
  default = "../../cmd/app/bbva/subscribe/account_transactions/main.go"
}

variable "bbva_subscribe_card_transactions_script" {
  default = "../../cmd/app/bbva/subscribe/card_transactions/main.go"
}

variable "bbva_subscribe_move_money_script" {
  default = "../../cmd/app/bbva/subscribe/move_money/main.go"
}

variable "bbva_subscribe_other_notifications_script" {
  default = "../../cmd/app/bbva/subscribe/other_notifications/main.go"
}

# BBVA SNS Connector Task
variable "bbva_sns_conector_env_name" {
  default = "ppd"
}

variable "bbva_sns_connector_name" {
  default = "core-application-bbva-sns-connector"
}

variable "bbva_sns_connector_image" {
  default = "379379777492.dkr.ecr.us-west-2.amazonaws.com/master-core-platform-ecr-bbva-sns-connector"
}

variable "bbva_sns_connector_cpu" {
  default = 256
}

variable "bbva_sns_connector_mem" {
  default = 512
}

variable "bbva_sns_connector_image_tag" {
  default = "build1"
}

variable "bbva_sns_connector_desired_container_count" {
  default = 1
}

variable "bbva_sns_connector_min_container_count" {
  default = 1
}

variable "bbva_sns_connector_max_container_count" {
  default = 3
}
