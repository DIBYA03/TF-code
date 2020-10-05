module "naming" {
  source = "git@github.com:wiseco/terraform-module-naming.git"

  application = "${var.application}"
  aws_region  = "${var.aws_region}"
  component   = "${var.component}"
  environment = "${var.environment}"
}

module "cert" {
  source = "git@github.com:wiseco/terraform-module-acm-certificate.git"

  application            = "${var.application}"
  aws_profile            = "${var.aws_profile}"
  route53_aws_profile    = "${var.public_route53_account_profile}"
  aws_region             = "us-east-1"
  component              = "${var.component}"
  domain_name            = "${var.cloudfront_domain_name}"
  environment            = "${var.environment}"
  route53_hosted_zone_id = "${var.route53_hosted_zone_id}"
  team                   = "${var.team}"
}
