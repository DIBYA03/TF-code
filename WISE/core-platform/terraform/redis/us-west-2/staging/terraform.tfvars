aws_profile = "master-us-west-2-saml-roles-deployment"

aws_region = "us-west-2"

environment = "staging"
environment_alias = "stg"

# VPC
vpc_id = "vpc-043a00e2bba3b0dd8"

vpc_cidr_block = "10.3.0.0/16"
cidr_block_us_west_2a = "10.3.101.0/22"

total_node = 1

redis_engine_version = "5.0.5"

node_type = "cache.t2.micro"
redis_multi_az = "false"
