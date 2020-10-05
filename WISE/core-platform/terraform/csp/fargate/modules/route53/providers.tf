provider "aws" {
  alias   = "${var.provider_name}"
  region  = "${var.aws_region}"
  profile = "${var.aws_profile}"
}
