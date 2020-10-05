resource "aws_security_group" "default" {
  name        = "${module.naming.aws_security_group}"
  description = "${terraform.workspace} default RDS security group"
  vpc_id      = "${var.vpc_id}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_security_group}"
    Team        = "${var.team}"
  }
}

resource "aws_security_group_rule" "allow_vpc" {
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow VPC traffic"
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.vpc_cidr_block}"]
}

resource "aws_security_group_rule" "allow_shared" {
  count             = "${var.shared_vpc_cidr_block == "" ? 0 : 1}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow shared VPC traffic"
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.shared_vpc_cidr_block}"]
}

resource "aws_security_group_rule" "allow_other_cidr_blocks" {
  count             = "${length(var.other_vpc_cidr_block) >= 1 ? 1 : 0}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow other VPC traffic"
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.other_vpc_cidr_block}"]
}

resource "aws_security_group_rule" "allow_csp" {
  count             = "${var.csp_vpc_cidr_block == "" ? 0 : 1}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow CSP VPC traffic"
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.csp_vpc_cidr_block}"]
}

resource "aws_security_group_rule" "allow_peered_vpc_ingress" {
  count             = "${var.peered_vpc_cidr_block == "" ? 0 : 1}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow Peered VPC traffic"
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.peered_vpc_cidr_block}"]
}

resource "aws_security_group_rule" "allow_peered_vpc_egress" {
  count             = "${var.peered_vpc_cidr_block == "" ? 0 : 1}"
  security_group_id = "${aws_security_group.default.id}"
  description       = "Allow Peered VPC traffic"
  type              = "egress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = ["${var.peered_vpc_cidr_block}"]
}

resource "aws_security_group" "lambda_backup" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  name        = "${module.naming.aws_security_group}-lambda-backup"
  description = "${terraform.workspace} lambda backup"
  vpc_id      = "${var.vpc_id}"

  # Needed, since there's no VPC endpoint for RDS yet
  egress {
    description = "Allow access over HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_security_group}-lambda-backup"
    Team        = "${var.team}"
  }
}
