provider "aws" {
  alias   = "${var.provider_name}"
  region  = "${var.aws_region}"
  profile = "${var.aws_profile}"
}

provider "aws" {
  alias   = "${var.route53_provider_name}"
  region  = "${var.aws_region}"
  profile = "${var.route53_aws_profile}"
}
