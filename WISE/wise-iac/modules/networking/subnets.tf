resource "aws_subnet" "public_subnets" {
  count = "${length((var.availability_zones))}"

  vpc_id            = "${aws_vpc.main.id}"
  cidr_block        = "${element(var.public_subnet_cidr_blocks, count.index)}"
  availability_zone = "${element(var.availability_zones, count.index)}"

  map_public_ip_on_launch = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.${element(var.availability_zones, count.index)}.public"
    Team        = "${var.team}"
    Tier        = "Public"
  }
}

resource "aws_subnet" "app_subnets" {
  count = "${length((var.availability_zones))}"

  vpc_id            = "${aws_vpc.main.id}"
  cidr_block        = "${element(var.app_subnet_cidr_blocks, count.index)}"
  availability_zone = "${element(var.availability_zones, count.index)}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.${element(var.availability_zones, count.index)}.app"
    Team        = "${var.team}"
    Tier        = "app"
  }
}

resource "aws_subnet" "db_subnets" {
  count = "${length((var.availability_zones))}"

  vpc_id            = "${aws_vpc.main.id}"
  cidr_block        = "${element(var.db_subnet_cidr_blocks, count.index)}"
  availability_zone = "${element(var.availability_zones, count.index)}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.${element(var.availability_zones, count.index)}.db"
    Team        = "${var.team}"
    Tier        = "db"
  }
}
