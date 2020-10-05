variable "provider_name" {
  default = "fargate-route53"
}

variable "aws_profile" {}
variable "aws_region" {}

variable "domain_name" {}
variable "route53_hosted_zone_id" {}
variable "resource_alias_name" {}
variable "resource_alias_zone_id" {}
