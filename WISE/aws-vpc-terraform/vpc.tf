# VPC
resource "aws_vpc" "main" {
  cidr_block           = "${var.vpc_cidr_block}"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.vpc"
    Team        = "${var.team}"
  }
}

resource "aws_internet_gateway" "internet_gw" {
  vpc_id = "${aws_vpc.main.id}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.igw"
    Team        = "${var.team}"
  }
}
