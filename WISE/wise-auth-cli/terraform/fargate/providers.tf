provider "aws" {
  profile = "${var.aws_profile}"
  region  = "${var.aws_region}"

  version = "~> 2.45"
}

provider "null" {
  version = "~> 2.1"
}

provider "template" {
  version = "~> 2.1"
}
