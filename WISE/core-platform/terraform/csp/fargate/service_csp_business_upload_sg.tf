resource "aws_security_group" "csp_business_upload_ecs" {
  name        = "${module.naming.aws_security_group}-csp-business-upload-ecs"
  description = "Security group for CSP business upload ECS container"

  vpc_id = "${var.vpc_id}"

  egress {
    description = "https outbound"
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
    description = "ecs to core postrgres"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"

    cidr_blocks = [
      "${var.core_db_cidr_blocks}",
    ]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-csp-business-upload-ecs"
    Team        = "${var.team}"
  }
}
