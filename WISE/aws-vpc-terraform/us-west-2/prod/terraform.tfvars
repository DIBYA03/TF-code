aws_profile = "wise-prod"

aws_region = "us-west-2"

environment = "prod"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDW6I1mqICyCgSdKJpkWrcmalb3Y1Qvw18xG19je95PvJ64qii8uMzpk1428scyYSwuxv8DFHzPBiI7xr+BvMhJhpZH3S3gUPZ7/VIU8YwKfCZmt20XPoNqrftE9YoD46O2BgmxpntsCmYoeAiYqSG+vJrGYByxpggICl1HCwme/2XGtNyf+Ioisp3ctc/N6RuoVdWlpU72vvj4Vz23DQdFhkk010uz1qcHMI+IitST1XvWGXd3bVveGAa0gdsZh2uUStJ/emOCn08Y1sYbMSHQA2zc7rvnu9YKNa/OiiEr3BJMAQ6QuIY7eegyBdBboDw0VmZ3jbiDgczqRTDkJUwF"

availability_zones = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]

vpc_cidr_block = "10.17.0.0/16"

rds_cidr_block = "10.17.0.0/16"

# Route 53 route associations
route53_private_zone_vpc_association_ids = [
  {
    # us-west-2-csp-prod
    "vpc_id" = "vpc-0a3d5ba3cf7256441"
    "region" = "us-west-2"
  },
  {
    # us-east-1-csp-prod
    "vpc_id" = "vpc-0b0012e20fb5bc05e"
    "region" = "us-east-1"
  },
  {
    # us-east-1-prod
    "vpc_id" = "vpc-0c3996a3cb4a1e002"
    "region" = "us-east-1"
  },
]

public_subnet_cidr_blocks = [
  "10.17.0.0/22",
  "10.17.4.0/22",
  "10.17.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.17.12.0/22",
  "10.17.16.0/22",
  "10.17.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.17.24.0/22",
  "10.17.28.0/22",
  "10.17.32.0/22",
]

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# App tier subnets
app_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # Docker container port 3000
    protocol   = "tcp"
    cidr_block = "10.17.0.0/16"
    from_port  = 3000
    to_port    = 3000
  },
]

# DB tier subnets
db_subnet_vpc_peering_routes = [
  {
    # us-east-1 prod
    destination_cidr_block    = "10.18.0.0/16"
    vpc_peering_connection_id = "pcx-03c1a52135facfecc"
  },
  {
    # csp vpc
    destination_cidr_block    = "10.4.0.0/16"
    vpc_peering_connection_id = "pcx-0390252234d5b1ccc"
  },
  {
    # us-west-2 beta-prod
    destination_cidr_block    = "10.20.0.0/16"
    vpc_peering_connection_id = "pcx-074b05edca70ab56c"
  },
]

db_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # postgres from us-east-1 prod vpc
    protocol   = "tcp"
    cidr_block = "10.18.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
  {
    # postgres from us-west-2-csp-prod
    protocol   = "tcp"
    cidr_block = "10.20.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
]

db_subnet_vpc_peering_egress_nacl_rules = [
  {
    # ephemeral to csp vpc
    protocol   = "tcp"
    cidr_block = "10.4.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
  {
    # ephemeral to us-east-1 prod vpc
    protocol   = "tcp"
    cidr_block = "10.18.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
  {
    # ephemeral to us-west-2 csp-prod
    protocol   = "tcp"
    cidr_block = "10.20.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
]

# Bastion Host
enable_bastion_host = "true"

bastion_host_instance_type = "t2.micro"

# Pagerduty
enable_pagerduty_slack_integration = "true"

pagerduty_non_critical_service_id = "PHRG9GX"

pagerduty_critical_service_id = "PL14NVK"
