provider "aws" {
  profile = "${var.aws_profile}"
  region  = "${var.aws_region}"

  version = "~> 2.44"
}

provider "local" {
  version = "~> 1.4"
}

provider "null" {
  version = "~> 2.1"
}

provider "random" {
  version = "~> 2.2"
}

provider "template" {
  version = "~> 2.1"
}
