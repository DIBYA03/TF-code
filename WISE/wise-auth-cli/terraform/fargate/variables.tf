variable "aws_profile" {}

variable "aws_region" {}
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

# Route 53
variable "public_route53_account_profile" {
  default = "master-us-west-2-saml-roles-admin"
}

variable "public_route53_hosted_zone" {}
variable "private_route53_hosted_zone" {}

# SNS
variable "non_critical_sns_topic" {}

variable "critical_sns_topic" {}

# VPC
variable "vpc_id" {}

variable "vpc_cidr_block" {}

variable "app_subnet_ids" {
  type = "list"
}

# kms
variable "default_kms_alias" {}

# ecr
variable "tagged_image_count_limit" {
  default = 10
}

variable "untagged_image_count_limit" {
  default = 5
}

# cloudwatch
variable "cw_log_group_retention_in_days" {}

# aws_vpn_auth service

variable "aws_vpn_auth_domain" {}

variable "aws_vpn_auth_allowed_account_ids" {
  type = "list"
}

variable "aws_vpn_nat_gateway_ip_cidr_blocks" {
  type = "list"

  default = [
    "34.213.158.195/32",
    "35.160.105.202/32",
  ]
}

variable "aws_vpn_auth_image_tag" {}

variable "aws_vpn_auth_desired_container_count" {}

variable "aws_vpn_auth_min_container_count" {}

variable "aws_vpn_auth_max_container_count" {}

variable "aws_vpn_auth_container_port" {
  default = 3000
}

variable "aws_vpn_auth_name" {
  default = "aws-cli-auth"
}

variable "aws_vpn_auth_cpu" {
  default = 256
}

variable "aws_vpn_auth_mem" {
  default = 512
}

# wise aws auth cli
variable "aws_vpn_auth_wise_image_tag" {}

# endpoint service
variable "endpoint_service_vpc_id" {}

variable "endpoint_service_subnet_ids" {
  type = "list"
}

variable "endpoint_service_allowed_cidr_blocks" {
  type = "list"
}
