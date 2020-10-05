resource "aws_security_group" "default_ecs" {
  name        = "${module.naming.aws_security_group}-default"
  description = "Allow only https for ${var.environment_name} ecs tasks"

  vpc_id = "${var.vpc_id}"

  ingress {
    description = "HTTPS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  egress {
    description = "HTTPS outbound"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "HTTPS outbound"
    from_port   = "${var.grpc_port}"
    to_port     = "${var.grpc_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-default"
    Team        = "${var.team}"
  }
}
