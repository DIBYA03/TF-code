variable "aws_profile" {}
variable "aws_region" {}
variable "environment" {}
variable "application" {}
variable "component" {}
variable "team" {}

variable "default_ssh_key" {}

variable "availability_zones" {
  type = "list"
}

variable "vpc_cidr_block" {}
variable "rds_cidr_block" {}

variable "route53_private_zone_vpc_association_ids" {
  type    = "list"
  default = []
}

# Route 53
variable "internal_route53_zone_id" {
  default = "Z7Q2TZ46NCJYD"
}

variable "public_subnet_cidr_blocks" {
  type = "list"
}

variable "app_subnet_cidr_blocks" {
  type = "list"
}

variable "db_subnet_cidr_blocks" {
  type = "list"
}

variable "custom_nacl_multiplier" {
  default = 5
}

variable "vpn_cidr_block" {
  description = "This isn't used, but states that VPN cider block for connections"
  default     = ""
}

# KMS grant permissioms to cross account
variable "kms_grant_cross_account_resources" {
  type    = "list"
  default = []
}

# VPC Peering
## public
variable "public_subnet_vpc_peering_routes" {
  type    = "list"
  default = []
}

variable "public_subnet_vpc_peering_ingress_nacl_rules" {
  type    = "list"
  default = []
}

variable "public_subnet_vpc_peering_egress_nacl_rules" {
  type    = "list"
  default = []
}

## app
variable "app_subnet_vpc_peering_routes" {
  type    = "list"
  default = []
}

variable "app_subnet_vpc_peering_ingress_nacl_rules" {
  type    = "list"
  default = []
}

variable "app_subnet_vpc_peering_egress_nacl_rules" {
  type    = "list"
  default = []
}

## db
variable "db_subnet_vpc_peering_routes" {
  type    = "list"
  default = []
}

variable "db_subnet_vpc_peering_ingress_nacl_rules" {
  type    = "list"
  default = []
}

variable "db_subnet_vpc_peering_egress_nacl_rules" {
  type    = "list"
  default = []
}

# VPN
variable "enable_vpn" {
  default = "false"
}

variable "vpn_domain" {
  default = ""
}

variable "vpn_sub_domain" {
  default = ""
}

variable "vpn_route53_zone_id" {
  default = ""
}

variable "vpn_ami" {
  default = ""
}

variable "vpn_instance_type" {
  default = ""
}

# Flow logs
variable "enable_flow_logs" {
  default = false
}

# allowed ip lists
# GitHub CIDR blocks
variable "github_cidr_blocks" {
  type = "list"

  default = [
    "192.30.252.0/22",
    "185.199.108.0/22",
    "140.82.112.0/20",
    "192.30.252.0/22",
    "185.199.108.0/22",
    "140.82.112.0/20",
    "13.229.188.59/32",
    "13.250.177.223/32",
    "18.194.104.89/32",
    "18.195.85.27/32",
    "35.159.8.160/32",
    "52.74.223.119/32",
    "192.30.252.153/32",
    "192.30.252.154/32",
    "185.199.108.153/32",
    "185.199.109.153/32",
    "185.199.110.153/32",
    "185.199.111.153/32",
    "54.87.5.173/32",
    "54.166.52.62/32",
    "23.20.92.3/32",
  ]
}

variable "apt_cidr_blocks" {
  type = "list"

  default = [
    "204.145.124.244/32",
    "217.196.149.55/32",
    "34.210.25.51/32",
    "34.212.136.213/32",
    "54.190.18.91/32",
    "54.191.55.41/32",
    "54.191.70.203/32",
    "54.218.137.160/32",
    "72.32.157.246/32",
    "87.238.57.227/32",
    "91.189.88.149/32",
    "91.189.88.162/32",
    "91.189.88.173/32",
    "91.189.88.174/32",
    "91.189.88.24/32",
    "91.189.88.31/32",
    "91.189.91.14/32",
    "91.189.91.23/32",
    "91.189.91.24/32",
    "91.189.91.26/32",
    "91.189.95.83/32",
  ]
}

# Bastion host
variable "enable_bastion_host" {
  default = false
}

variable "bastion_host_hostname" {
  default = ""
}

variable "allowed_principals" {
  type = "list"

  default = [
    # Master account with VPN
    "arn:aws:iam::379379777492:root",
  ]
}

variable "bastion_host_instance_type" {
  default = ""
}

variable "bastion_host_port" {
  default = 22
}

variable "bastion_host_min_size" {
  default = 1
}

variable "bastion_host_max_size" {
  default = 1
}

variable "bastion_host_desired_capacity" {
  default = 1
}

variable "bastion_host_vpc_endpoint_service_list" {
  type    = "list"
  default = []
}

variable "bastion_host_salstack_s3_object_prefix" {
  default = "bastion-host"
}

variable "bastion_host_salstack_s3_object_name" {
  default = "saltstack.zip"
}

# Pagerduty

variable "enable_pagerduty_slack_integration" {
  default = "false"
}

variable "pagerduty_token" {}
variable "pagerduty_slack_access_token" {}
variable "pagerduty_non_critical_service_id" {}
variable "pagerduty_critical_service_id" {}
variable "pagerduty_slack_configuration_url" {}
variable "pagerduty_slack_url" {}
