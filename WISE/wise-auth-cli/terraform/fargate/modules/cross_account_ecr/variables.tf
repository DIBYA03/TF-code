variable "aws_profile" {}

variable "provider_name" {
  default = "cross-account-ecr"
}

variable "aws_region" {
  default = "us-west-2"
}

variable "environment" {}
variable "environment_name" {}

variable "application" {
  default = "auth"
}

variable "component" {
  default = "cli"
}

variable "team" {
  default = "security"
}

# ecr
variable "tagged_image_count_limit" {
  default = 10
}

variable "untagged_image_count_limit" {
  default = 5
}
