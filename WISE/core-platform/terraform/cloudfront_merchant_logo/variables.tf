variable "aws_profile" {}
variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "merchant"
}

variable "component" {
  default = "logo"
}

variable "team" {
  default = "cloud-ops"
}

# Route 53
variable "public_route53_account_profile" {
  default = "master-us-west-2-saml-roles-deployment"
}

variable "route53_hosted_zone_id" {
  default = "Z3BUXRPXJI78KB"
}

# CloudFront
variable "cloudfront_domain_name" {}

variable "ecs_merchant_logo_domain" {}

variable "cloudfront_country_restriction_type" {}

variable "cloudfront_country_allowed_country_codes" {
  default = ["US"]
  type    = "list"
}

variable "cloudfront_price_class" {}

variable "cloudfront_add_waf" {
  default = "false"
}
