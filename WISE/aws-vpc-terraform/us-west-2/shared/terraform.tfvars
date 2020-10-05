aws_profile = "wiseus"

aws_region = "us-west-2"

environment = "shared"

application = "wise-us"

component = "vpc"

team = "cloud-ops"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9djRH90f87p60c/QKjPauzG0Dyj3j2R9s9fERrxlYtvPnXjkTi9iAhtGMPDhrAqqVDJlHh9djo5gY6cHHCh7/23XR8eqJVIyJXI2/Rk7pgywZt7rox6xQHnHs0guQMX6t6dcI85gg1dhNp1JpK6VZXDvjnjGfarFgnx/hUnrp8WwnybKwN2o7hkyXbDsqEnQpfvCQxi+dfhCQqgTEITsLZWERtrP2OBVkdP6Ky6/bMAzHdvE04An74kHmaNxXiPydjbjM1tUicWesqgohRpL2ueUusRFWzrcTRVeWRQqzcNLp/oCHLJ7RzUj2rMmgh7/s5gKSCulazS4nDMh4D6PICcv5igNb8SVt0KNzRIpXMx8hFJsnp3fkiQE9a5+f9RZhSv6iB79jlKLpfkGfG5wANsHckSAkG16o3iLfZ2B7h6vnH5XuKalVzFJJ2IA1DjCtNPhxldflKO1Kopztk14q7LROLsxP2UYT1jEVvZoPkF3FxQBriKqmFBoGp9O6WMoZs0KY7eNT1/5lD/DcP09bInx+J1zAauP7nvR23vum+duam5ynwt3TC70bBNBrA2t+FteGAjjVrZERuZAzVONhwX9iYXdUXU7Abvc1iymTWPa0tt3oISU7A7ztwM5izV4gijbsTUUMHu524pAjmonO1nAmxCAT2n28I9w7nUygaQ=="

availability_zones = [
  "us-west-2a",
  "us-west-2b",
]

vpc_cidr_block = "10.2.0.0/16"

rds_cidr_block = "10.2.0.0/16"

public_subnet_cidr_blocks = [
  "10.2.0.0/22",
  "10.2.4.0/22",
]

app_subnet_cidr_blocks = [
  "10.2.12.0/22",
  "10.2.16.0/22",
]

db_subnet_cidr_blocks = [
  "10.2.24.0/22",
  "10.2.28.0/22",
]

# Bastion Host Private Endpoints
# this is used to create a private link between VPCs without having to peer
bastion_host_vpc_endpoint_service_list = [
  {
    # prod bastion
    service = "com.amazonaws.vpce.us-west-2.vpce-svc-0c71a43062d4901f1"
    domain  = "prod-bastion.internal.wise.us"
  },
  {
    # csp prod bastion
    service = "com.amazonaws.vpce.us-west-2.vpce-svc-0bafa0c126600ae70"
    domain  = "csp-prod-bastion.internal.wise.us"
  },
  {
    # security bastion
    service = "com.amazonaws.vpce.us-west-2.vpce-svc-0a87b0c6922251210"
    domain  = "security-bastion.internal.wise.us"
  },
  {
    # security private-ca
    service = "com.amazonaws.vpce.us-west-2.vpce-svc-09b1f1f58d7fc44f5"
    domain  = "ca.internal.wise.us"
  },
]

# VPN
vpn_cidr_block = "10.2.36.0/22" # Noting this here, but used in VPN only

enable_vpn = "true"

vpn_domain = "wise.us"

vpn_sub_domain = "vpn"

vpn_route53_zone_id = "Z3BUXRPXJI78KB"

vpn_ami = "ami-0fd17a62ca7cbc538"

vpn_instance_type = "t3.small"

# VPC Peering
# REMINDER: There is a 20 rule limit on NACLs

# Public tier subnets
public_subnet_vpc_peering_routes = [
  {
    # dev vpc
    destination_cidr_block    = "10.1.0.0/16"
    vpc_peering_connection_id = "pcx-0cbb40f16a2df10a1"
  },
  {
    # staging vpc
    destination_cidr_block    = "10.3.0.0/16"
    vpc_peering_connection_id = "pcx-05f1623d86cd6fbb9"
  },
  {
    # csp vpc
    destination_cidr_block    = "10.4.0.0/16"
    vpc_peering_connection_id = "pcx-0f52f977304f7ceb2"
  },
  {
    # csp prod vpc
    destination_cidr_block    = "10.20.0.0/16"
    vpc_peering_connection_id = "pcx-0b1d04180f189deb0"
  },
  {
    # dev1 vpc
    destination_cidr_block    = "10.24.0.0/16"
    vpc_peering_connection_id = "pcx-097f1ad118933911d"
  },
]

public_subnet_vpc_peering_ingress_nacl_rules = [
  {
    # vpn public in
    protocol   = "udp"
    cidr_block = "0.0.0.0/0"
    from_port  = 1194
    to_port    = 1194
  },
]

public_subnet_vpc_peering_egress_nacl_rules = [
  {
    # dev RDS
    protocol   = "tcp"
    cidr_block = "10.1.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
  {
    # dev redis
    protocol   = "tcp"
    cidr_block = "10.1.0.0/16"
    from_port  = 6379
    to_port    = 6379
  },
  {
    # staging RDS
    protocol   = "tcp"
    cidr_block = "10.3.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
  {
    # csp https
    protocol   = "tcp"
    cidr_block = "10.4.0.0/16"
    from_port  = 443
    to_port    = 443
  },
  {
    # csp RDS
    protocol   = "tcp"
    cidr_block = "10.4.0.0/16"
    from_port  = 5432
    to_port    = 5432
  },
  {
    # csp prod https
    protocol   = "tcp"
    cidr_block = "10.20.0.0/16"
    from_port  = 443
    to_port    = 443
  },
  {
    # dev1 most ports
    protocol   = "tcp"
    cidr_block = "10.24.0.0/16"
    from_port  = 22
    to_port    = 3000
  },
]

# Pagerduty
enable_pagerduty_slack_integration = "true"

pagerduty_non_critical_service_id = "PL4Y6KP"

pagerduty_critical_service_id = "PHVN67M"
