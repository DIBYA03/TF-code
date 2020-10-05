variable "provider_name" {
  default = "terraform-module-rds"
}

variable "aws_profile" {}
variable "aws_region" {}
variable "environment" {}
variable "application" {}
variable "component" {}
variable "team" {}

variable "rds_vpc_id" {}
variable "rds_vpc_cidr_block" {}
variable "peering_vpc_cidr_block" {}
variable "shared_vpc_cidr_block" {}

# SNS
variable "non_critical_sns_topic" {}

variable "critical_sns_topic" {}

# RDS
variable "rds_replicate_source_db" {
  default = ""
}

variable "rds_kms_key_arn" {
  default = ""
}

variable "iam_database_authentication_enabled" {
  default = false
}

variable "rds_instance_count" {}
variable "rds_database_name" {}

variable "rds_deletion_protection" {
  default = "true"
}

variable "rds_engine" {}
variable "rds_engine_version" {}
variable "rds_parameter_group_family_name" {}

variable "rds_max_connections" {}

variable "rds_instance_class" {}
variable "rds_storage_size" {}
variable "rds_storage_type" {}

variable "rds_skip_final_snapshot" {}

variable "db_subnet_group_ids" {
  type = "list"
}

variable "rds_storage_encrypted" {
  default = true
}

variable "rds_ca_cert_identifier" {}

variable "rds_multi_az" {}

variable "rds_publicly_accessible" {
  default = false
}

# CloudWatch Alarms

variable "rds_cw_cpu_non_critical_limit" {}
variable "rds_cw_cpu_critical_limit" {}
variable "rds_cw_conn_count_non_critical" {}
variable "rds_cw_conn_count_critical" {}
variable "rds_cw_free_mem_non_critical" {}
variable "rds_cw_free_mem_critical" {}
variable "rds_cw_free_disk_non_critical" {}
variable "rds_cw_free_disk_critical" {}
