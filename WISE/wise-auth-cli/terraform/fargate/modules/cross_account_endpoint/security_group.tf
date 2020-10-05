resource "aws_security_group" "default" {
  name        = "${module.naming.aws_security_group}-endpoint"
  description = "security group for ${module.naming.aws_security_group} endpoint"

  vpc_id = "${var.vpc_id}"

  ingress {
    description = "allow incoming traffic"
    from_port   = "${var.endpoint_incoming_port}"
    to_port     = "${var.endpoint_incoming_port}"
    protocol    = "TCP"

    cidr_blocks = [
      "${var.allowed_cidr_blocks}",
    ]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-endpoint"
    Team        = "${var.team}"
  }

  provider = "aws.${var.provider_name}"
}
