resource "aws_security_group" "default" {
  count       = "${var.rds_instance_count}"
  name        = "${module.naming.aws_security_group}"
  description = "${terraform.workspace} default RDS security group"
  vpc_id      = "${var.rds_vpc_id}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_security_group}"
    Team        = "${var.team}"
  }

  provider = "aws.${var.provider_name}"
}

resource "aws_security_group_rule" "allow_vpc" {
  count             = "${var.rds_instance_count}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow VPC traffic"
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.rds_vpc_cidr_block}"]

  provider = "aws.${var.provider_name}"
}

resource "aws_security_group_rule" "allow_shared" {
  count             = "${var.shared_vpc_cidr_block != "" && var.rds_instance_count != 0 ? 1 : 0}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow shared VPC traffic"
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.shared_vpc_cidr_block}"]

  provider = "aws.${var.provider_name}"
}

resource "aws_security_group_rule" "allow_peered_ingress" {
  count             = "${var.peering_vpc_cidr_block != "" && var.rds_instance_count != 0 ? 1 : 0}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow Peered VPC traffic"
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.peering_vpc_cidr_block}"]

  provider = "aws.${var.provider_name}"
}

resource "aws_security_group_rule" "allow_peered_egress" {
  count             = "${var.peering_vpc_cidr_block != "" && var.rds_instance_count != 0 ? 1 : 0}"
  count             = "${var.peering_vpc_cidr_block == "" ? 0 : 1}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow Peered VPC traffic"
  type              = "egress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.peering_vpc_cidr_block}"]

  provider = "aws.${var.provider_name}"
}
