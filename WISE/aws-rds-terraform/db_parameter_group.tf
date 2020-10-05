resource "aws_db_parameter_group" "default" {
  name   = "${module.naming.aws_db_parameter_group}"
  family = "${var.rds_parameter_group_family_name}"

  parameter {
    name  = "max_connections"
    value = "${local.rds_max_connections}"

    apply_method = "pending-reboot"
  }

  parameter {
    name  = "rds.force_ssl"
    value = true
  }
}
