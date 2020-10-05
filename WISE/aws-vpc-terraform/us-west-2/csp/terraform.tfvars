aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "csp"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCh0vv8/kWPe/B+tMNivW35d9fXAjziu2+7UXZcWEUZ1tUBJABuDIYMRhzyteBdHdQzpY2zK7JjrGGQwHpL0Qmz1qudn91NxtrTT9ixrH0vb1bf7U1QtpRbiz4k+7g602hOYb/HfmZYLCoFdfrhb0WiWcHEGRkSwPxZk7ywQl3NAvqJ8ShSyhE4pcaONVMfFQoZjg0vvhQfp8DvqcXYkZJA/LDIkw38bQZLeb+XHSe0cp3RjFUEB5Q7LwrwnqY5h4pcY9m4JmYeukE2O/ylCRU/SEDolSINejKwJGYYh8gXaKoUdzSsD7gi5DZdVCe2jpDgnLb0Y7D5EkE7Pz+OXwy/"

availability_zones = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]

vpc_cidr_block = "10.4.0.0/16"

rds_cidr_block = "10.4.0.0/16"

public_subnet_cidr_blocks = [
  "10.4.0.0/22",
  "10.4.4.0/22",
  "10.4.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.4.12.0/22",
  "10.4.16.0/22",
  "10.4.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.4.24.0/22",
  "10.4.28.0/22",
  "10.4.32.0/22",
]

# KMS cross-account permissions
kms_grant_cross_account_resources = [
  "arn:aws:iam::178124264531:root",
  "arn:aws:iam::058450407364:root",
]

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# App tier subnets
app_subnet_vpc_peering_routes = [
  {
    # shared vpc
    destination_cidr_block    = "10.2.0.0/16"
    vpc_peering_connection_id = "pcx-0f52f977304f7ceb2"
  },
  {
    # staging vpc
    destination_cidr_block    = "10.3.0.0/16"
    vpc_peering_connection_id = "pcx-00dbf97c6204231c8"
  },
]

# NEED TO FIX THIS
# {
#   # dev vpc
#   destination_cidr_block    = "10.1.0.0/16"
#   vpc_peering_connection_id = "pcx-0a1142f622b44c608"
# },

app_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # postgres from shared vpc
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 443
    to_port    = 443
  },
  {
    # ephemeral to staging vpc
    protocol   = "tcp"
    cidr_block = "10.4.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
  {
    # ephemeral to dev vpc
    protocol   = "tcp"
    cidr_block = "10.1.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
]

app_subnet_vpc_peering_egress_nacl_rules = [
  {
    # ephemeral to shared vpc
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
  {
    # to staging postgres
    protocol   = "tcp"
    cidr_block = "10.3.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
  {
    # to dev postgres
    protocol   = "tcp"
    cidr_block = "10.1.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
]

# App tier subnets
db_subnet_vpc_peering_routes = [
  {
    # shared vpc
    destination_cidr_block    = "10.2.0.0/16"
    vpc_peering_connection_id = "pcx-0f52f977304f7ceb2"
  },
]

db_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # postgres from shared vpc
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
]

db_subnet_vpc_peering_egress_nacl_rules = [
  {
    # ephemeral to shared vpc
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
]

# Pagerduty
enable_pagerduty_slack_integration = "true"

pagerduty_non_critical_service_id = "PWAQ43J"

pagerduty_critical_service_id = "PT123RW"
