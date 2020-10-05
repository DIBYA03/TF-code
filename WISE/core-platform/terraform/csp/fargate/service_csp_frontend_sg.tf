resource "aws_security_group" "csp_frontend_ecs_alb" {
  name        = "${module.naming.aws_security_group}-csp-frontend-ecs-alb"
  description = "Allow https alb to money request ecs"

  vpc_id = "${var.vpc_id}"

  ingress {
    description = "HTTPS inbound from stripe"
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"

    cidr_blocks = [
      "${var.vpc_cidr_block}",
      "${var.shared_cidr_block}",
    ]
  }

  egress {
    from_port   = "${var.csp_frontend_container_port}"
    to_port     = "${var.csp_frontend_container_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }
 
  egress {
    description = "alb to container"
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-csp-frontend-ecs-alb"
    Team        = "${var.team}"
  }
}

resource "aws_security_group" "csp_frontend_ecs" {
  name        = "${module.naming.aws_security_group}-csp-frontend-ecs"
  description = "Allow only https from csp frontend alb"

  vpc_id = "${var.vpc_id}"

  ingress {
    description     = "alb to container"
    from_port       = "${var.csp_frontend_container_port}"
    to_port         = "${var.csp_frontend_container_port}"
    protocol        = "TCP"
    security_groups = ["${aws_security_group.csp_frontend_ecs_alb.id}"]
  }

  egress {
    description = "https outbound"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-csp-frontend-ecs"
    Team        = "${var.team}"
  }
}
