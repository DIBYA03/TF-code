resource "aws_security_group" "batch_monthly_interest_ecs" {
  name        = "${module.naming.aws_security_group}-batch-monthly-interest-ecs"
  description = "batch monthly interest security group"

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
    Name        = "${module.naming.aws_security_group}-batch-monthly-interest-ecs"
    Team        = "${var.team}"
  }
}
