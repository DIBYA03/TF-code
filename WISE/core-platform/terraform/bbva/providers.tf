provider "aws" {
  profile = "${var.aws_profile}"
  region  = "${var.aws_region}"

  version = "~> 2.13"
}

provider "template" {
  version = "~> 2.1"
}
