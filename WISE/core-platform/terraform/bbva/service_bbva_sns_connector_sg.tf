resource "aws_security_group" "bbva_sns_connector_ecs" {
  name        = "${module.naming.aws_security_group}-bbva-sns-connector-ecs"
  description = "BBVA notifications Security Group"

  vpc_id = "${var.vpc_id}"

  egress {
    description = "HTTPS outbound"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-bbva-sns-connector-ecs"
    Team        = "${var.team}"
  }
}
