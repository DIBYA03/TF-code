provider "aws" {
  profile = "${var.aws_profile}"
  region  = "${var.aws_region}"

  version = "~> 2.43"
}

provider "archive" {
  version = "~> 1.3"
}

provider "null" {
  version = "~> 2.1"
}
