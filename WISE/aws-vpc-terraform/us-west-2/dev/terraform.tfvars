aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "dev"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDcG7bHjoV17oTicC3QStd2SrjaP1oGBF/Zkyeyfku2k5EzchQ9snLlZIOwtVIoW/EDcoBbEYo4NeLUVzLeDm/H9KS6asjEEjltPs4f803W8rnIAAP/N+SE8oPtT/ebeULrKqcBok8lSmzLq5SDl2Gr3xf4vGU2RGNfiw6gb1rpqWQO9AV1bTnj/LlXCaS7prDrZhfM7tNpByXdvdwIqPqoPUzmek+3m5N8z8eMnWa5Yc1OLCrlD8m+wbM+1tpcjeKNyszrZ2ZaJG2leIPFwTGTnX2xQppgNbvehEKxRCBoTHBuuydYI4UfI6Vspryv0HrbrKCSkDsYBmlDfrbOmz/U3K17+/d9FLGifddEw7DObrPko9AXvs8Tl0rk/Ti1K4SK9w2uqUcoGzd6wzPQnXK8U50Oqs2xAzYEvl2xJaTwhjP6LLBXowpUhaKgsGcxir2IV2raOSswJENUw9JAlBTl32stfLnAW73WxTCVJcjMdngDS/Ewt26WfLV7fGIkB43otPn2QHfPNnMeWcTpDhg++1ff7Ep+MTwXwflV+v9ZNoWkir4tAefmkvbz10/fxQH4XNC56jsLlAZw98gehGGRB/XHFnBWkzGtTG6IrpmxVeQsBq6ZMw0csJ7e8BE2PyrO8N173lrOrJCAdtAFSth9M8YPVbs2hwJi704HsGKpBQ=="

availability_zones = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]

vpc_cidr_block = "10.1.0.0/16"

rds_cidr_block = "10.1.0.0/16"

public_subnet_cidr_blocks = [
  "10.1.0.0/22",
  "10.1.4.0/22",
  "10.1.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.1.12.0/22",
  "10.1.16.0/22",
  "10.1.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.1.24.0/22",
  "10.1.28.0/22",
  "10.1.32.0/22",
]

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# App tier subnets
# app_subnet_vpc_peering_routes = [
#   {
#     # shared vpc
#     destination_cidr_block    = "10.2.0.0/16"
#     vpc_peering_connection_id = "pcx-0cbb40f16a2df10a1"
#   },
# ]

app_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # Docker container port 3000
    protocol   = "tcp"
    cidr_block = "10.1.0.0/16"
    from_port  = 3000
    to_port    = 3000
  },
]

# DB tier subnets
db_subnet_vpc_peering_routes = [
  {
    # us-east-1 dev
    destination_cidr_block    = "10.5.0.0/16"
    vpc_peering_connection_id = "pcx-0541e5cce4d51cf99"
  },
  {
    # csp vpc
    destination_cidr_block    = "10.4.0.0/16"
    vpc_peering_connection_id = "pcx-0a1142f622b44c608"
  },
  {
    # shared vpc
    destination_cidr_block    = "10.2.0.0/16"
    vpc_peering_connection_id = "pcx-0cbb40f16a2df10a1"
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
  {
    # postgres from csp vpc
    protocol   = "tcp"
    cidr_block = "10.4.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
  {
    # postgres from us-east-1 dev vpc
    protocol   = "tcp"
    cidr_block = "10.5.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
  {
    # pgbouncer from shared vpc
    protocol   = "tcp"
    cidr_block = "10.2.0.0/16"
    from_port  = 6432
    to_port    = 6433
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
  {
    # ephemeral to csp vpc
    protocol   = "tcp"
    cidr_block = "10.4.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
  {
    # ephemeral to us-east-1 dev vpc
    protocol   = "tcp"
    cidr_block = "10.5.0.0/16"
    from_port  = 1024
    to_port    = 65535
  },
]

pagerduty_non_critical_service_id = "P7Y7IID"

pagerduty_critical_service_id = "PFD0DT9"

enable_pagerduty_slack_integration = "true"
