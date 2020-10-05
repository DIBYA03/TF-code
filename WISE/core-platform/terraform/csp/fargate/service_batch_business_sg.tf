resource "aws_security_group" "batch_business_ecs" {
  name        = "${module.naming.aws_security_group}-batch-business-ecs"
  description = "batch account security group"

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

    cidr_blocks = [
      "${var.core_db_cidr_blocks}",
      "${var.csp_rds_cidr_block}",
    ]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-batch-business-ecs"
    Team        = "${var.team}"
  }
}
