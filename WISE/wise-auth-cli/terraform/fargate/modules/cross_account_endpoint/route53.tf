resource "aws_route53_record" "default" {
  zone_id = "${var.route53_hosted_zone}"
  name    = "${var.endpoint_domain_name}"
  type    = "CNAME"
  ttl     = "300"
  records = ["${lookup(aws_vpc_endpoint.default.dns_entry[0], "dns_name")}"]

  provider = "aws.${var.provider_name}"
}
