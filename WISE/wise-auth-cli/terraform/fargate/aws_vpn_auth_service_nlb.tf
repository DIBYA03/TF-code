resource "aws_lb" "aws_vpn_auth" {
  name               = "${module.naming.aws_alb}-aws"
  load_balancer_type = "network"
  internal           = true

  enable_cross_zone_load_balancing = true
  enable_deletion_protection       = true

  subnets = ["${var.app_subnet_ids}"]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_alb}-aws"
    Team        = "${var.team}"
  }
}

resource "aws_lb_target_group" "aws_vpn_auth" {
  name_prefix = "awscli"
  port        = "${var.aws_vpn_auth_container_port}"
  protocol    = "TLS"
  vpc_id      = "${var.vpc_id}"
  target_type = "ip"

  health_check = {
    enabled  = true
    port     = "traffic-port"
    protocol = "HTTPS"
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_alb_target_group}-aws"
    Team        = "${var.team}"
  }
}

resource "aws_lb_listener" "aws_vpn_auth_tls" {
  load_balancer_arn = "${aws_lb.aws_vpn_auth.arn}"
  port              = "443"
  protocol          = "TLS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${module.aws_vpn_auth.arn}"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.aws_vpn_auth.arn}"
  }
}

module "aws_vpn_auth" {
  source = "git@github.com:wiseco/terraform-module-acm-certificate.git"

  application            = "${var.application}"
  aws_profile            = "${var.aws_profile}"
  route53_aws_profile    = "${var.aws_profile}"
  aws_region             = "${var.aws_region}"
  component              = "${var.component}"
  domain_name            = "${var.aws_vpn_auth_domain}"
  environment            = "${var.environment}"
  route53_aws_profile    = "${var.public_route53_account_profile}"
  route53_hosted_zone_id = "${var.public_route53_hosted_zone}"
  team                   = "${var.team}"
}
