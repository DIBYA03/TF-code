resource "aws_db_parameter_group" "default" {
  count  = "${var.rds_instance_count}"
  name   = "${module.naming.aws_db_parameter_group}"
  family = "${var.rds_parameter_group_family_name}"

  parameter {
    name  = "max_connections"
    value = "${var.rds_max_connections}"

    apply_method = "pending-reboot"
  }

  provider = "aws.${var.provider_name}"
}
