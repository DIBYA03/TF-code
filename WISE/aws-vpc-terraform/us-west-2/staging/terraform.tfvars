aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "staging"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCmzKylTTQhs0B2dBQyUNP644UXLNS5CH+TT2I1R2+s6uWw6PKJpPq81CawEsM5GA7ZfIBqE0Fa+dbz2YvHEhFSDAC6NgNX0yv45euAc8JMIF7H/jhLN5GUrjY6WRgFGZmCDXKWu7ovonNBHrzmEBqnEe0QyoLETHBs8hoOVOMHQ1UK2VnqyRPW8Oj6tzt/05+X6gzoTehv1y/jh2W8YdZP8ljO1piW4AA8qNhtM9sQmtagdkI3piED3rBY0lFN3krHQQnUVxwiBj7PUa7YTSRxSyt+uoNe36XUJ3WAs7eI/y/DbZpb6KKnGIdxClXBunE28Qd0jUgKGx84KVbk406h"

availability_zones = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]

vpc_cidr_block = "10.3.0.0/16"

rds_cidr_block = "10.3.0.0/16"

public_subnet_cidr_blocks = [
  "10.3.0.0/22",
  "10.3.4.0/22",
  "10.3.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.3.12.0/22",
  "10.3.16.0/22",
  "10.3.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.3.24.0/22",
  "10.3.28.0/22",
  "10.3.32.0/22",
]

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# DB tier subnets
# the below is only here for when access is needed fora short period of time
db_subnet_vpc_peering_routes = [
  {
    # shared vpc
    destination_cidr_block    = "10.2.0.0/16"
    vpc_peering_connection_id = "pcx-05f1623d86cd6fbb9"
  },
  {
    # csp vpc
    destination_cidr_block    = "10.4.0.0/16"
    vpc_peering_connection_id = "pcx-00dbf97c6204231c8"
  },
]

db_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # postgres from shared vpc
    protocol   = "tcp"
    cidr_block = "10.4.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
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
    # ephemeral to csp vpc
    egress     = true
    protocol   = "tcp"
    cidr_block = "10.4.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
  {
    # ephemeral to shared vpc
    egress     = true
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
]

# Pagerduty
enable_pagerduty_slack_integration = "true"

pagerduty_non_critical_service_id = "PAL6NIO"

pagerduty_critical_service_id = "PEXUOXB"
