# Route53 entry
module "client_api" {
  source = "./modules/route53"

  aws_profile            = "${var.public_route53_account_profile}"
  aws_region             = "${var.aws_region}"
  domain_name            = "${var.cloudfront_domain_name}"
  route53_hosted_zone_id = "${var.route53_hosted_zone_id}"
  resource_alias_name    = "${aws_cloudfront_distribution.client_api.domain_name}"
  resource_alias_zone_id = "${aws_cloudfront_distribution.client_api.hosted_zone_id}"
}
