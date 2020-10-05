variable "provider_name" {
  default = "acm-certificate-module"
}

variable "route53_provider_name" {
  default = "route53-acm-certificate-module"
}

variable "aws_profile" {}
variable "route53_aws_profile" {}

variable "aws_region" {}

variable "environment" {}
variable "application" {}
variable "component" {}
variable "team" {}

variable "domain_name" {}

variable "subject_alternative_names" {
  default = []
  type    = "list"
}

variable "route53_hosted_zone_id" {}
