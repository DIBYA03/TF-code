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
    description = "GRPC"
    from_port   = "${var.grpc_port}"
    to_port     = "${var.grpc_port}"
    protocol    = "tcp"

    cidr_blocks = [
      "${var.core_db_cidr_blocks}",
      "${var.vpc_cidr_block}",
    ]
  }

  egress {
    description = "Postgres to CSP RDS VPC"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"

    cidr_blocks = [
      "${var.core_db_cidr_blocks}",
    ]
  }

  egress {
    description = "Postgres to core RDS VPCs"
    from_port   = 5432
    to_port     = 5432
    protocol    = "TCP"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }


  egress {
    description = "Partner Proxy port"
    from_port   = 8080
    to_port     = 8080
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }


  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-lambda"
    Team        = "${var.team}"
  }
}
