aws_profile = "wise-dev"

aws_region = "us-west-2"

environment = "dev1"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCz/BRuvMkbnsQmApOaGNeuTVYtIgK9sQubZ/RIjzJBPZVaYhAE/7x1E6d5IQDgO44v2Looxlfw23Okg1HNxShof8hJQ7CZKpt1WQFoYf5N9UAQEiRy/wvTqZ/RFnCTMySgQBEiNTQkABsYDLTj7af08jFN0GoQm5SV45IGCZuCJnU5Df1uWnBj6MC3H/r1ZImY6ZE40gY89CitdP+QUVvLug+sBceJO0vbWN3XmyAjx4w1Yb1Q/K4w55WuhGMGwpkC4gUvj4J46TY3P5TVdIVERdu0xEp/dDTLx/tAsfboOw9WYp7hqmeUj/7SlygqWdxI8wDMUHyeymhAAeiknd5mlWXGUkdS+hP/r9TnVOMQ7ou/O25MzGyFtFVOGxSDKJroMwtWXfOnzuCBhmxgoftKWoUvgtTS1ErOyNzA+g1SrY2Nr5nVvhdAV1vQ0G3gX2oxzygdma75UopZtLZKL1fNgR/rxkB0KnxRlXbnp1oHRIDrvdy1JBRPeai3XPgX99GhVcTPuV3WxiM3qPGNlW6hPxf3zQwoOk8QioRRDxcoZ43mtEvairUeomnISfj7oIlbakzXSwk+08ixsgYvBmAwU9AfZx5fwgiM/vuOExVTbLoPYEEpnAFTX6fmtg44LvLxRA2WANZ640MwX1Ge9it6qVAiI7z+irWzlgy30SuUAw=="

availability_zones = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]

vpc_cidr_block = "10.24.0.0/16"

rds_cidr_block = "10.24.0.0/16"

public_subnet_cidr_blocks = [
  "10.24.0.0/22",
  "10.24.4.0/22",
  "10.24.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.24.12.0/22",
  "10.24.16.0/22",
  "10.24.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.24.24.0/22",
  "10.24.28.0/22",
  "10.24.32.0/22",
]

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# app tier subnets
app_subnet_vpc_peering_routes = [
  {
    # shared vpc
    destination_cidr_block    = "10.2.0.0/16"
    vpc_peering_connection_id = "pcx-097f1ad118933911d"
  },
]

app_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # ssh  from shared
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 22
    to_port    = 3000
  },
]

# db tier subnets
db_subnet_vpc_peering_routes = [
  {
    # shared vpc
    destination_cidr_block    = "10.2.0.0/16"
    vpc_peering_connection_id = "pcx-097f1ad118933911d"
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
  {
    # redis from shared vpc
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 6379
    to_port    = 6379
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

pagerduty_non_critical_service_id = "P7Y7IID"

pagerduty_critical_service_id = "PFD0DT9"

enable_pagerduty_slack_integration = "true"
