variable "aws_profile" {}

variable "aws_region" {}
variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "core-platform"
}

variable "component" {
  default = "ecr"
}

variable "team" {
  default = "cloud-ops"
}

variable "allowed_account_principals" {
  type = "list"
}

variable "tagged_image_count_limit" {
  default = 20
}

variable "untagged_image_count_limit" {
  default = 5
}
