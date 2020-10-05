resource "aws_route_table" "app" {
  count  = "${length((var.availability_zones))}"
  vpc_id = "${aws_vpc.main.id}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.${element(var.availability_zones, count.index)}.app.rt"
    Team        = "${var.team}"
  }
}

resource "aws_route" "app_to_nat_gw" {
  count                  = "${length((var.availability_zones))}"
  route_table_id         = "${aws_route_table.app.*.id[count.index]}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${aws_nat_gateway.public_subnet_gateways.*.id[count.index]}"

  depends_on = ["aws_route_table.app"]
}

resource "aws_route" "app_vpc_peering_routes" {
  count = "${length(var.availability_zones) * length(var.app_subnet_vpc_peering_routes)}"

  route_table_id            = "${element(aws_route_table.app.*.id, count.index % length(var.availability_zones))}"
  destination_cidr_block    = "${lookup(var.app_subnet_vpc_peering_routes[count.index / length(var.availability_zones)], "destination_cidr_block")}"
  vpc_peering_connection_id = "${lookup(var.app_subnet_vpc_peering_routes[count.index / length(var.availability_zones)], "vpc_peering_connection_id")}"
}

resource "aws_route_table_association" "app" {
  count = "${length((var.availability_zones))}"

  subnet_id      = "${aws_subnet.app_subnets.*.id[count.index]}"
  route_table_id = "${aws_route_table.app.*.id[count.index]}"
}
