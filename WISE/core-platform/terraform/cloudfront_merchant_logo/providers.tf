provider "aws" {
  profile = "${var.aws_profile}"
  region  = "${var.aws_region}"

  version = "~> 2.44"
}

provider "archive" {
  version = "~> 1.3"
}
