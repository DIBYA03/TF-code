resource "aws_acm_certificate" "default" {
  domain_name       = "${var.domain_name}"
  validation_method = "DNS"

  subject_alternative_names = "${var.subject_alternative_names}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_acm_certificate}"
    Team        = "${var.team}"
  }

  provider = "aws.${var.provider_name}"
}

resource "aws_route53_record" "default" {
  name    = "${aws_acm_certificate.default.domain_validation_options.0.resource_record_name}"
  type    = "${aws_acm_certificate.default.domain_validation_options.0.resource_record_type}"
  zone_id = "${var.route53_hosted_zone_id}"
  records = ["${aws_acm_certificate.default.domain_validation_options.0.resource_record_value}"]
  ttl     = 60

  provider = "aws.${var.route53_provider_name}"
}

# Subject Alternative Names
resource "aws_route53_record" "subject_alternative_names" {
  count   = "${length(var.subject_alternative_names) >= 1 ? 1 : 0}"
  name    = "${aws_acm_certificate.default.domain_validation_options.1.resource_record_name}"
  type    = "${aws_acm_certificate.default.domain_validation_options.1.resource_record_type}"
  zone_id = "${var.route53_hosted_zone_id}"
  records = ["${aws_acm_certificate.default.domain_validation_options.1.resource_record_value}"]
  ttl     = 60

  provider = "aws.${var.route53_provider_name}"
}

# resource "aws_acm_certificate_validation" "default" {
#  certificate_arn = "${aws_acm_certificate.default.arn}"


#  validation_record_fqdns = [
#    "${aws_route53_record.default.fqdn}",
#    "${aws_route53_record.subject_alternative_names.*.fqdn}",
#  ]


#  depends_on = [
#    "aws_route53_record.default",
#  ]


#  provider = "aws.${var.provider_name}"
#}

