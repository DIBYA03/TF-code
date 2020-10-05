# Route53 entry
module "merchant_logo" {
  source = "./modules/route53"

  aws_profile            = "${var.public_route53_account_profile}"
  aws_region             = "${var.aws_region}"
  domain_name            = "${var.cloudfront_domain_name}"
  route53_hosted_zone_id = "${var.route53_hosted_zone_id}"
  resource_alias_name    = "${aws_cloudfront_distribution.merchant_logo.domain_name}"
  resource_alias_zone_id = "${aws_cloudfront_distribution.merchant_logo.hosted_zone_id}"
}
