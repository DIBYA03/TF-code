resource "aws_security_group" "batch_monitor_ecs" {
  name        = "${module.naming.aws_security_group}-batch-monitor-ecs"
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
    description = "postgres outbound"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["${var.client_api_rds_vpc_cidr_block}"]
  }

  egress {
    description = "beam outbound"
    from_port   = 8888
    to_port     = 8888
    protocol    = "tcp"
    cidr_blocks = ["${var.client_api_rds_vpc_cidr_block}"]
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
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_security_group}-batch-monitor-ecs"
    Team        = "${var.team}"
  }
}
