aws_profile = "wise-prod"

aws_region = "us-west-2"

environment = "csp-prod"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDFMbHb9FyrmI1wYZAoW8t6q3h31IsUSOHSQIA7tq8HAS/VsXAXpRbl3RbMWPy+O10wFXF9DjJ8268yjFp+t8DUmP73prQhajg/jj2XTL01VNgo5/1vmfcxlSyPUKBuUPRBP7E/m5IFP9juCzy/GIQ2GXUgX9MnjdgaDVYLiP6H9btLHaUPlwUlhjez4ORTf6gynk6J7EkGKBdOcd+Jhn4QL9MVQbObbSaGDF/nChngbCJ3t1BANWGfVfB3i14FYLau9/h3MDS8ld4EhkSUt3BZBaTKMf+blgqK7BSDuVO/Lk0mkln/ueMU0Y40OUZugPBOi2lvLXZmF/KKb9ihYUiz"

availability_zones = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]

vpc_cidr_block = "10.20.0.0/16"

rds_cidr_block = "10.13.0.0/16"

# Route 53 route associations
route53_private_zone_vpc_association_ids = [
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
  {
    # us-east-1-prod
    "vpc_id" = "vpc-0c3996a3cb4a1e002"
    "region" = "us-east-1"
  },
]

public_subnet_cidr_blocks = [
  "10.20.0.0/22",
  "10.20.4.0/22",
  "10.20.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.20.12.0/22",
  "10.20.16.0/22",
  "10.20.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.20.24.0/22",
  "10.20.28.0/22",
  "10.20.32.0/22",
]

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# App tier subnets
app_subnet_vpc_peering_routes = [
  {
    # prod vpc
    destination_cidr_block    = "10.17.0.0/16"
    vpc_peering_connection_id = "pcx-074b05edca70ab56c"
  },
  {
    # prod vpc
    destination_cidr_block    = "10.2.0.0/16"
    vpc_peering_connection_id = "pcx-0b1d04180f189deb0"
  },
]

app_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # ephemeral from beta-prod vpc
    protocol   = "tcp"
    cidr_block = "10.17.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
  {
    # https from shared vpc
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 443
    to_port    = 443
  },
]

app_subnet_vpc_peering_egress_nacl_rules = [
  {
    # to prod postgres
    protocol   = "tcp"
    cidr_block = "10.17.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
  {
    # ephemeral to shared vpc
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
]

# Bastion Host
enable_bastion_host = "true"

bastion_host_instance_type = "t2.micro"

# Pagerduty
enable_pagerduty_slack_integration = "true"

pagerduty_non_critical_service_id = "PNMSTMT"

pagerduty_critical_service_id = "P7Q4TN3"
