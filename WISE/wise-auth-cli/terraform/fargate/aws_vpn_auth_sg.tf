resource "aws_security_group" "aws_vpn_auth_ecs" {
  name        = "${module.naming.aws_security_group}-ecs"
  description = "security group for ${module.naming.aws_security_group} service"

  vpc_id = "${var.vpc_id}"

  ingress {
    description = "allow incoming traffic"
    from_port   = "${var.aws_vpn_auth_container_port}"
    to_port     = "${var.aws_vpn_auth_container_port}"
    protocol    = "TCP"

    cidr_blocks = [
      "${var.vpc_cidr_block}",
    ]
  }

  egress {
    description = "HTTPS outbound to reach google SAML"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-ecs"
    Team        = "${var.team}"
  }
}
