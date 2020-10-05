resource "aws_security_group" "lambda_default" {
  name        = "${module.naming.aws_security_group}-lambda"
  description = "Allow TLS inbound traffic to lambdas in ${var.environment_name}"

  vpc_id = "${var.vpc_id}"

  ingress {
    description = "HTTPS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Postgres to VPC"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"

    cidr_blocks = [
      # Sometimes the core DB is in another VPC, like beta-prod and prod
      "${var.core_rds_cidr_block == "" ? var.vpc_cidr_block : var.core_rds_cidr_block}",
    ]
  }

  egress {
    description = "Redis to VPC"
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  egress {
    description = "GRPC"
    from_port   = "${var.grpc_port}"
    to_port     = "${var.grpc_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-lambda"
    Team        = "${var.team}"
  }
}
