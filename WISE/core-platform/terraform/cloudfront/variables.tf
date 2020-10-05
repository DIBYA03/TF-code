variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "client"
}

variable "component" {
  default = "api"
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

variable "cloudfront_logging_bucket" {
  default = "global-account-resources-logging.s3.amazonaws.com"
}

variable "clodfront_country_restriction_type" {}

variable "clodfront_country_allowed_country_codes" {
  default = ["US"]
  type    = "list"
}

variable "cloudfront_price_class" {}

variable "cloudfront_add_waf" {
  default = "false"
}
