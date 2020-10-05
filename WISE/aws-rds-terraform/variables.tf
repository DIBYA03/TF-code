variable "aws_profile" {}
variable "aws_region" {}
variable "application" {}
variable "component" {}
variable "team" {}

variable "vpc_id" {}
variable "vpc_cidr_block" {}

variable "shared_vpc_cidr_block" {
  default = ""
}

variable "other_vpc_cidr_block" {
  type    = "list"
  default = []
}

variable "csp_vpc_cidr_block" {
  default = ""
}

variable "peered_vpc_cidr_block" {
  default = ""
}

# SNS
variable "non_critical_sns_topic" {}

variable "critical_sns_topic" {}

variable "rds_cross_region_non_critical_sns_topic" {
  default = ""
}

variable "rds_cross_region_critical_sns_topic" {
  default = ""
}

# route53
variable "route53_zone_id" {}

variable "route53_master_domain" {}
variable "route53_read_replica_domain" {}

# RDS Subnet Group
variable "app_subnet_ids" {
  type = "list"
}

variable "db_subnet_group_ids" {
  type = "list"
}

# RDS
variable "rds_deletion_protection" {
  default = true
}

variable "rds_instance_class" {}

variable "rds_instance_class_mem" {
  description = "in GiB and found here: https://aws.amazon.com/rds/instance-types/"
}

variable "rds_instance_max_connection_divider" {
  default = 5000000
}

variable "rds_storage_size" {}
variable "rds_storage_type" {}

variable "rds_storage_encrypted" {
  default = true
}

variable "rds_iam_database_authentication_enabled" {
  default = false
}

variable "rds_ca_cert_identifier" {
  default = "rds-ca-2019"
}

variable "rds_engine" {}
variable "rds_engine_version" {}
variable "rds_parameter_group_family_name" {}
variable "rds_read_replica_count" {}
variable "rds_username" {}
variable "rds_password" {}
variable "rds_multi_az" {}

variable "database_name" {
  default = "wiseus"
}

variable "rds_backup_retention_period" {}

variable "backup_db_snapshot_count_limit" {
  default = 24
}

variable "rds_maintenance_window" {
  default = "Tue:00:00-Tue:03:00"
}

variable "rds_replica_maintenance_window" {
  default = "Tue:03:00-Tue:05:00"
}

# Cross-Region replication
variable "rds_cross_region_kms_key_arn" {
  default = ""
}

variable "rds_cross_region_multi_az" {
  default = ""
}

variable "rds_cross_region_vpc_id" {
  default = ""
}

variable "rds_cross_region_vpc_cidr_block" {
  default = ""
}

variable "rds_cross_region_db_subnet_group_ids" {
  type    = "list"
  default = []
}

# CloudWatch Alarms

variable "rds_cw_conn_count_non_critical" {}

variable "rds_cw_conn_count_critical" {}
variable "rds_cw_cpu_non_critical_limit" {}
variable "rds_cw_cpu_critical_limit" {}
variable "rds_cw_free_mem_non_critical" {}
variable "rds_cw_free_mem_critical" {}
variable "rds_cw_free_disk_non_critical" {}
variable "rds_cw_free_disk_critical" {}

# Backups
variable "enable_cross_region_backup" {
  default = "true"
}

variable "backup_lambda_timeout" {}

variable "rds_enable_backup_lambda" {}
