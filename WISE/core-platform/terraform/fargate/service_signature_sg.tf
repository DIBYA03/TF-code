resource "aws_security_group" "signature_ecs" {
  name        = "${module.naming.aws_security_group}-signature-ecs"
  description = "app notifications Security Group"

  vpc_id = "${var.vpc_id}"

  egress {
    description = "HTTPS outbound"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "postgres outbound"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["${var.client_api_rds_vpc_cidr_block}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-signature-ecs"
    Team        = "${var.team}"
  }
}
