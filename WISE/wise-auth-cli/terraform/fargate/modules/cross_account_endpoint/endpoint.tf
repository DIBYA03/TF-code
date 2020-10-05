resource "aws_vpc_endpoint" "default" {
  vpc_id            = "${var.vpc_id}"
  service_name      = "${var.endpoint_service}"
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    "${aws_security_group.default.id}",
  ]

  subnet_ids          = ["${var.endpoint_subnet_ids}"]
  private_dns_enabled = false

  provider = "aws.${var.provider_name}"
}
