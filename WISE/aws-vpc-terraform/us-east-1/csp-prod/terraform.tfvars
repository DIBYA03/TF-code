aws_profile = "wise-prod"

aws_region = "us-east-1"

environment = "csp-prod"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDGSEL2A6chiQadnuRncHgxTL+UcoCBiDHEN0DcBtJPAfukRtIYFVnxsTelln8mTkixUezUIME7+yeFqCvXnunmsuXLV7kAErdxMMbtGmMMR4zLM9ZsRKP3BLZmZCkZdrGkqV8O4Mz4G1SRN49U4rh28M5SoO9tMIeJdN4s46Qv+pz2CcZmOCoHd5quPl1NnrNp1mCdoa7OqLFN2TlEG7z9A9PcmxGYMUdkPQekQ5Ajw7tnlgsHyJEhnCt85Xt4zeO35hQQcTkDyoiuBB03DlI59SHzrCjXDVfgTS7keytvW70TNz3U73yurGLAmFo5YKTFFp2snCeCoqq3rNDcyP/N"

availability_zones = [
  "us-east-1a",
  "us-east-1b",
  "us-east-1c",
]

vpc_cidr_block = "10.22.0.0/16"

rds_cidr_block = "10.22.0.0/16"

# Route 53 route associations
route53_private_zone_vpc_association_ids = [
  {
    # us-west-2-csp-prod
    "vpc_id" = "vpc-0a3d5ba3cf7256441"
    "region" = "us-west-2"
  },
  {
    # us-west-2-prod
    "vpc_id" = "vpc-0dbc76a72cb520546"
    "region" = "us-west-2"
  },
  {
    # us-east-1-prod
    "vpc_id" = "vpc-0c3996a3cb4a1e002"
    "region" = "us-east-1"
  },
]

public_subnet_cidr_blocks = [
  "10.22.0.0/22",
  "10.22.4.0/22",
  "10.22.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.22.12.0/22",
  "10.22.16.0/22",
  "10.22.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.22.24.0/22",
  "10.22.28.0/22",
  "10.22.32.0/22",
]

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# App tier subnets
# app_subnet_vpc_peering_routes = [
#   {
#     # shared vpc
#     destination_cidr_block    = "10.2.0.0/16"
#     vpc_peering_connection_id = "pcx-0f52f977304f7ceb2"
#   },
#   {
#     # staging vpc
#     destination_cidr_block    = "10.3.0.0/16"
#     vpc_peering_connection_id = "pcx-00dbf97c6204231c8"
#   },
# ]

app_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # ephemeral from qa-prod vpc
    protocol   = "tcp"
    cidr_block = "10.13.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
]

app_subnet_vpc_peering_egress_nacl_rules = [
  {
    # to qa-prod postgres
    protocol   = "tcp"
    cidr_block = "10.13.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
]

pagerduty_non_critical_service_id = "PNMSTMT"

pagerduty_critical_service_id = "P7Q4TN3"
