resource "aws_vpc_endpoint_service" "aws_vpn_auth" {
  acceptance_required        = true
  network_load_balancer_arns = ["${aws_lb.aws_vpn_auth.arn}"]

  allowed_principals = ["${var.aws_vpn_auth_allowed_account_ids}"]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_alb}-aws"
    Team        = "${var.team}"
  }
}

module "cross_account_endpoint" {
  source = "./modules/cross_account_endpoint"

  aws_profile            = "${var.public_route53_account_profile}"
  environment            = "${var.environment}"
  environment_name       = "${var.environment_name}"
  team                   = "${var.team}"
  vpc_id                 = "${var.endpoint_service_vpc_id}"
  route53_hosted_zone    = "${var.private_route53_hosted_zone}"
  endpoint_subnet_ids    = ["${var.endpoint_service_subnet_ids}"]
  endpoint_incoming_port = "443"
  endpoint_service       = "${aws_vpc_endpoint_service.aws_vpn_auth.service_name}"
  allowed_cidr_blocks    = ["${var.endpoint_service_allowed_cidr_blocks}"]
  endpoint_domain_name   = "${var.aws_vpn_auth_domain}"
}
