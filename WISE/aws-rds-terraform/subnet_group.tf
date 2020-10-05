resource "aws_db_subnet_group" "default" {
  name = "${var.application}.${terraform.workspace}.subnet_group"

  subnet_ids = [
    "${var.db_subnet_group_ids}",
  ]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${var.application}.${terraform.workspace}.subnet_group"
    Team        = "${var.team}"
  }
}
