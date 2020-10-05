module "naming" {
  source = "git@github.com:wiseco/terraform-module-naming.git"

  application = "${var.application}"
  aws_region  = "${var.aws_region}"
  component   = "${var.component}"
  environment = "${terraform.workspace}"
}

module "cross_region_read_only" {
  rds_instance_count = "${var.enable_cross_region_backup == "true" ? 1 : 0}"
  source             = "./modules/replicas"

  application = "${var.application}"
  aws_region  = "us-east-1"
  aws_profile = "${var.aws_profile}"
  component   = "${var.component}"
  environment = "${terraform.workspace}"
  team        = "${var.team}"

  rds_kms_key_arn                     = "${var.rds_cross_region_kms_key_arn}"
  iam_database_authentication_enabled = "${var.rds_iam_database_authentication_enabled}"

  rds_database_name       = "${var.database_name}"
  rds_replicate_source_db = "${aws_db_instance.master.arn}"
  rds_multi_az            = "${var.rds_cross_region_multi_az}"
  rds_publicly_accessible = false
  rds_skip_final_snapshot = "false"
  rds_ca_cert_identifier  = "${var.rds_ca_cert_identifier}"

  rds_instance_class              = "${var.rds_instance_class}"
  rds_parameter_group_family_name = "${var.rds_parameter_group_family_name}"
  rds_engine                      = "${var.rds_engine}"
  rds_engine_version              = "${var.rds_engine_version}"

  rds_storage_size      = "${var.rds_storage_size}"
  rds_storage_type      = "${var.rds_storage_type}"
  rds_storage_encrypted = "${var.rds_storage_encrypted}"

  rds_vpc_id             = "${var.rds_cross_region_vpc_id}"
  rds_vpc_cidr_block     = "${var.rds_cross_region_vpc_cidr_block}"
  peering_vpc_cidr_block = "${var.vpc_cidr_block}"
  shared_vpc_cidr_block  = "${var.shared_vpc_cidr_block}"
  db_subnet_group_ids    = "${var.rds_cross_region_db_subnet_group_ids}"

  non_critical_sns_topic = "${var.rds_cross_region_non_critical_sns_topic}"
  critical_sns_topic     = "${var.rds_cross_region_critical_sns_topic}"

  rds_max_connections = "${local.rds_max_connections}"

  rds_cw_conn_count_non_critical = "${local.rds_non_critical_conn_count_limit}"
  rds_cw_conn_count_critical     = "${local.rds_critical_conn_count_limit}"

  rds_cw_free_mem_non_critical = "${local.rds_non_critical_mem_limit}"
  rds_cw_free_mem_critical     = "${local.rds_critical_mem_limit}"

  rds_cw_cpu_non_critical_limit = "${var.rds_cw_cpu_non_critical_limit}"
  rds_cw_cpu_critical_limit     = "${var.rds_cw_cpu_critical_limit}"

  rds_cw_free_disk_non_critical = "${local.rds_non_critical_disk_limit}"
  rds_cw_free_disk_critical     = "${local.rds_critical_disk_limit}"
}
