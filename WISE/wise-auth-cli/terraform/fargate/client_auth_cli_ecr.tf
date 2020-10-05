module "client_auth_cli" {
  source = "./modules/cross_account_ecr"

  aws_profile      = "${var.public_route53_account_profile}"
  environment      = "${var.environment}"
  environment_name = "${var.environment_name}"
  team             = "${var.team}"
}

output "client_auth_cli_command" {
  value = "${module.client_auth_cli.aws_auth_cli_wise_deploy_command}"
}
