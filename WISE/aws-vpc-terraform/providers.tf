provider "aws" {
  region  = "${var.aws_region}"
  profile = "${var.aws_profile}"

  version = "~> 2.43"
}

provider "pagerduty" {
  token = "${var.pagerduty_token}"

  version = "~> 1.4"
}

provider "template" {
  version = "~> 2.1"
}
