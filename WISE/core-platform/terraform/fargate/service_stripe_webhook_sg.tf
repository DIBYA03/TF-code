resource "aws_security_group" "stripe_webhook_ecs" {
  name        = "${module.naming.aws_security_group}-stripe-webhook-ecs"
  description = "Allow only https from money request alb"

  vpc_id = "${var.vpc_id}"

  ingress {
    description     = "alb to container"
    from_port       = "${var.services_container_port}"
    to_port         = "${var.services_container_port}"
    protocol        = "TCP"
    security_groups = ["${aws_security_group.services.id}"]
  }

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
    Name        = "${module.naming.aws_security_group}-stripe-webhook-ecs"
    Team        = "${var.team}"
  }
}
