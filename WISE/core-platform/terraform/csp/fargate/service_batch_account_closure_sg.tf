resource "aws_security_group" "batch_account_closure_ecs" {
  name        = "${module.naming.aws_security_group}-batch-account-closure-ecs"
  description = "batch monitor security group"

  vpc_id = "${var.vpc_id}"

  # needed to get ECR image
  egress {
    description = "HTTPS outbound"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "ecs to csp postrgres"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"

    cidr_blocks = [
      "${var.csp_rds_cidr_block}",
    ]
  }

  egress {
    description = "postgres outbound"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["${var.core_db_cidr_blocks}"]
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
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_security_group}-batch-account-closure-ecs"
    Team        = "${var.team}"
  }
}
