resource "aws_security_group" "csp_review_ecs" {
  name        = "${module.naming.aws_security_group}-csp-review-ecs"
  description = "Security group for CSP review ECS container"

  vpc_id = "${var.vpc_id}"

    # needed to get ECR image
  egress {
    description = "https outbound"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "ecs to csp postgres"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"

    cidr_blocks = [
      "${var.csp_rds_cidr_block}",
    ]
  }

  egress {
    description = "ecs to core postgres"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"

    cidr_blocks = [
      "${var.core_db_cidr_blocks}",
    ]
  }

  egress {
    description = "beam outbound"
    from_port   = 8888
    to_port     = 8888
    protocol    = "tcp"
    cidr_blocks = ["${var.core_db_cidr_blocks}"]
  }

  egress {
    description = "GRPC outbound"
    from_port   = "${var.grpc_port}"
    to_port     = "${var.grpc_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  egress {
    description = "GRPC outbound"
    from_port   = "${var.grpc_port}"
    to_port     = "${var.grpc_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.core_db_cidr_blocks}"]
  }


  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-csp-review-ecs"
    Team        = "${var.team}"
  }
}
