resource "aws_route53_zone" "private_default" {
  name = "${var.environment == "shared" ? "internal.wise.us" : "${var.environment}.${var.aws_region}.internal.wise.us"}"

  vpc {
    vpc_id = "${aws_vpc.main.id}"
  }

  lifecycle {
    ignore_changes = [
      "vpc",
    ]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment == "shared" ? "internal.wise.us" : "${var.environment}.${var.aws_region}.internal.wise.us"}"
    Team        = "${var.team}"
  }
}

resource "aws_route53_zone_association" "wise_vpcs" {
  count = "${length(var.route53_private_zone_vpc_association_ids)}"

  zone_id    = "${aws_route53_zone.private_default.zone_id}"
  vpc_id     = "${lookup(var.route53_private_zone_vpc_association_ids[count.index], "vpc_id")}"
  vpc_region = "${lookup(var.route53_private_zone_vpc_association_ids[count.index], "region")}"
}

resource "aws_route53_record" "vpn" {
  count   = "${var.enable_vpn == "true" ? 1 : 0}"
  zone_id = "${var.vpn_route53_zone_id}"
  name    = "${var.vpn_sub_domain}.${var.vpn_domain}"
  type    = "A"
  ttl     = "900"
  records = ["${join("", aws_eip.vpn_eip.*.public_ip)}"]
}
