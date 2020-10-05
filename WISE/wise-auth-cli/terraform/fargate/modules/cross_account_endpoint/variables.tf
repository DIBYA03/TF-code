variable "aws_profile" {}

variable "provider_name" {
  default = "cross-account-endpoint"
}

variable "aws_region" {
  default = "us-west-2"
}

variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "auth"
}

variable "component" {
  default = "cli"
}

variable "team" {
  default = "security"
}

variable "vpc_id" {}

variable "endpoint_subnet_ids" {
  type = "list"
}

variable "endpoint_incoming_port" {}

variable "endpoint_service" {}

variable "allowed_cidr_blocks" {
  type = "list"
}

variable "route53_hosted_zone" {}
variable "endpoint_domain_name" {}
