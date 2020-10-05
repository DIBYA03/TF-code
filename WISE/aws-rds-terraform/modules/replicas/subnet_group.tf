resource "aws_db_subnet_group" "default" {
  count = "${var.rds_instance_count}"
  name  = "${var.application}.${terraform.workspace}.subnet_group"

  subnet_ids = "${var.db_subnet_group_ids}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${var.application}.${terraform.workspace}.subnet_group"
    Team        = "${var.team}"
  }

  provider = "aws.${var.provider_name}"
}
