module "naming" {
  source = "git@github.com:wiseco/terraform-module-naming.git"

  application = "${var.application}"
  aws_region  = "${var.aws_region}"
  component   = "${var.component}"
  environment = "${var.environment_name}"
}
