aws_profile = "wise-prod"

aws_region = "us-east-1"

environment = "prod"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC5EraX5QSUuHw5UpPc2RuyaOda57kSHDRGMKWWxBH7iYt4Iefx2eHs3d7oQjK2QrApf5P7gS7TZ4IelZaxVqOVFbsTmk4oeHTGulrhCl4omVG3J1VIku9gV6PHznslNSfeog7vOZYG12q1uXLODbmBPlZV2XQc12rCwiSJ5oRmkqyAfERcCVnJYhJp8VWlVcZVH979lcuk3ZcSKj6f/PHhWzwoCNJoVxnypJMzvhTJL3P8ett6LzXrsuMDaAF65yw+4JFeLE++5EPwWJQ/Nn1tPeNtNG01CozyM+8V1Y5qiwSYnil5aUGTo/Cwk3EP3aQ4gTvRNlbsvGyNc/s+zo6L"

availability_zones = [
  "us-east-1a",
  "us-east-1b",
  "us-east-1c",
]

vpc_cidr_block = "10.18.0.0/16"

rds_cidr_block = "10.18.0.0/16"

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
    # us-west-2-prod
    "vpc_id" = "vpc-0dbc76a72cb520546"
    "region" = "us-west-2"
  },
]

public_subnet_cidr_blocks = [
  "10.18.0.0/22",
  "10.18.4.0/22",
  "10.18.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.18.12.0/22",
  "10.18.16.0/22",
  "10.18.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.18.24.0/22",
  "10.18.28.0/22",
  "10.18.32.0/22",
]

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# App tier subnets
app_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # Docker container port 3000
    protocol   = "tcp"
    cidr_block = "10.18.0.0/16"
    from_port  = 3000
    to_port    = 3000
  },
]

# DB tier subnets
db_subnet_vpc_peering_routes = [
  {
    # us-west-2 prod
    destination_cidr_block    = "10.17.0.0/16"
    vpc_peering_connection_id = "pcx-03c1a52135facfecc"
  },
]

db_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # postgres from us-west-2 prod vpc
    protocol   = "tcp"
    cidr_block = "10.17.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
]

db_subnet_vpc_peering_egress_nacl_rules = [
  {
    # ephemeral to us-weat-2 prod vpc
    protocol   = "tcp"
    cidr_block = "10.17.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
]

# Pagerduty
pagerduty_non_critical_service_id = "PHRG9GX"

pagerduty_critical_service_id = "PL14NVK"
