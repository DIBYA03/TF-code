resource "aws_security_group" "bastion_host" {
  count       = "${var.enable_bastion_host ? 1 : 0}"
  name        = "${module.naming.aws_security_group}-bst-hst"
  description = "security group for ${var.environment} bastion host"
  vpc_id      = "${aws_vpc.main.id}"

  ingress {
    from_port   = "${var.bastion_host_port}"
    to_port     = "${var.bastion_host_port}"
    protocol    = "tcp"
    cidr_blocks = "${var.app_subnet_cidr_blocks}"
  }

  egress {
    from_port = 443
    to_port   = 443
    protocol  = "tcp"

    cidr_blocks = [
      "0.0.0.0/0",
    ]
  }

  egress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"

    cidr_blocks = [
      "${var.github_cidr_blocks}",
      "${var.apt_cidr_blocks}",
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-bst-hst"
    Team        = "${var.team}"
  }
}
