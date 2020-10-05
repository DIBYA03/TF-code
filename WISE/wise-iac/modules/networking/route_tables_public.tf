resource "aws_route_table" "public" {
  vpc_id = "${aws_vpc.main.id}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.public.rt"
    Team        = "${var.team}"
  }
}

resource "aws_route" "public_to_igw" {
  route_table_id         = "${aws_route_table.public.id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.internet_gw.id}"

  depends_on = ["aws_route_table.public"]
}

resource "aws_route" "public_vpc_peering_routes" {
  count                     = "${length(var.public_subnet_vpc_peering_routes)}"
  route_table_id            = "${aws_route_table.public.id}"
  destination_cidr_block    = "${lookup(var.public_subnet_vpc_peering_routes[count.index], "destination_cidr_block")}"
  vpc_peering_connection_id = "${lookup(var.public_subnet_vpc_peering_routes[count.index], "vpc_peering_connection_id")}"
}


resource "aws_route_table_association" "public" {
  count = "${length((var.availability_zones))}"

  subnet_id      = "${aws_subnet.public_subnets.*.id[count.index]}"
  route_table_id = "${aws_route_table.public.id}"
}
