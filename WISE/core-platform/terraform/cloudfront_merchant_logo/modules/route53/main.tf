resource "aws_route53_record" "default" {
  zone_id = "${var.route53_hosted_zone_id}"
  name    = "${var.domain_name}"
  type    = "A"

  alias {
    name                   = "${var.resource_alias_name}"
    zone_id                = "${var.resource_alias_zone_id}"
    evaluate_target_health = true
  }

  provider = "aws.${var.provider_name}"
}
