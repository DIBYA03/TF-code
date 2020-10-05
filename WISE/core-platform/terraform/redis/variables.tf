variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}
variable "vpc_id" {}
variable "vpc_cidr_block" {}

variable "cidr_block_us_west_2a" {}

variable "application" {
  default = "core-platform"
}

variable "component" {
  default = "redis"
}

variable "team" {
  default = "cloud-ops"
}


variable "total_node" {
  default = 1
}

variable "redis_engine_version" {
  default = "5.0.6"
}

variable "node_type" {
  default = "cache.t2.micro"
}

variable "redis_multi_az" {
  default = "false"
}