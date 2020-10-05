resource "aws_eip" "nat_eip" {
  count = "${length((var.availability_zones))}"

  vpc        = true
  depends_on = ["aws_internet_gateway.internet_gw"]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.nat_gw.${element(var.availability_zones, count.index)}"
    Team        = "${var.team}"
  }
}

resource "aws_nat_gateway" "public_subnet_gateways" {
  count = "${length((var.availability_zones))}"

  allocation_id = "${aws_eip.nat_eip.*.id[count.index]}"
  subnet_id     = "${aws_subnet.public_subnets.*.id[count.index]}"
  depends_on    = ["aws_internet_gateway.internet_gw"]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.${element(var.availability_zones, count.index)}"
    Team        = "${var.team}"
  }
}
