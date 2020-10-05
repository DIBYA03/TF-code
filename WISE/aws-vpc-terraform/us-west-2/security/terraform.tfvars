aws_profile = "wise-security"

aws_region = "us-west-2"

environment = "security"

application = "wise-us"

component = "vpc"

team = "security"

default_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDAF5evCUaovWDDGalPKPJL3XtBz0dphfrIiqpz7tUEvIaPCgdo123V7767h7M6FYjWnY/l/nO+9NkQkMdT3RlwrUkO53J0gY2mzfv+M2VaSwgLHvqWvwnd7dBhQ+TuuoNnm0gEN24IblzXbfGeyFuCvcDxo9zjivdMj9IO0var3zEpUEqkCMnhMt1Swl2Syvra9kRDBmqEsyiAQLRRFlbpGCaG7xMJaiGbnPQH737FqpjJRqvMxcI4p4sxYbBqizrtAwgUcHvEp42qdLge3D8o9vLPDf2TZbOCWxcStm1EpoETAPNxzGH+OHFeYmvhjqYvVdAFcG4sCM9lGFIl8IA+Z4M0UxcG9vMyV/ipUjGuslu4NSHWpr/QqOZ9imNhPFjHPKA49eTM/8viQK8MDWpoOA1olySytdBg+zQnX8ZjZg4rc/ODKvYby5KWCpdF1CLH10/zIwtRFNpujGw36cIk2SBeJ38sJfSbcrd44ixF0nBhmWtJbyP4tZFmKiX6Q89wSrCqCj1JZLf21kxk4ilCgvR9Np9FERrCqDyuG/CCsGeBp0mYHqWlym+xe3wYEHeDoZBLV2YKi1PwJsI2GqiPcdwnf6C7ck5Gv6EdbeIYHO6gHzb1tOA1I27K6GBEsRgViSnFTVddCGczgVnfXqEj6NIwIkpFo8oOM5hOnxYGVQ=="

availability_zones = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]

vpc_cidr_block = "10.23.0.0/16"

rds_cidr_block = "10.23.0.0/16"

# Route 53 route associations
route53_private_zone_vpc_association_ids = []

public_subnet_cidr_blocks = [
  "10.23.0.0/22",
  "10.23.4.0/22",
  "10.23.8.0/22",
]

app_subnet_cidr_blocks = [
  "10.23.12.0/22",
  "10.23.16.0/22",
  "10.23.20.0/22",
]

db_subnet_cidr_blocks = [
  "10.23.24.0/22",
  "10.23.28.0/22",
  "10.23.32.0/22",
]

# Bastion Host
enable_bastion_host = true

bastion_host_hostname = ""

bastion_host_instance_type = "t3a.micro"

# Pagerduty
enable_pagerduty_slack_integration = "true"

pagerduty_non_critical_service_id = "PYHAHOR"

pagerduty_critical_service_id = "PCTLI0E"
