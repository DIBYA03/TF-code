module "merchant_logo_route_53" {
  source = "./modules/route53"

  aws_profile            = "${var.public_route53_account_profile}"
  aws_region             = "${var.aws_region}"
  domain_name            = "${var.merchant_logo_domain}"
  route53_hosted_zone_id = "${var.public_hosted_zone_id}"
  resource_alias_name    = "${aws_alb.services.dns_name}"
  resource_alias_zone_id = "${aws_alb.services.zone_id}"
}

resource "aws_alb_target_group" "merchant_logo" {
  name_prefix = "merlog"
  port        = "${var.services_container_port}"
  protocol    = "HTTPS"
  vpc_id      = "${var.vpc_id}"
  target_type = "ip"

  health_check = {
    enabled  = true
    path     = "/healthcheck.html"
    port     = "traffic-port"
    protocol = "HTTPS"
    matcher  = 200
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_alb_target_group}-merchant-logo"
    Team        = "${var.team}"
  }
}

resource "aws_lb_listener_rule" "merchant_logo" {
  listener_arn = "${aws_alb_listener.services_https.arn}"

  action {
    type             = "forward"
    target_group_arn = "${aws_alb_target_group.merchant_logo.arn}"
  }

  condition {
    field  = "host-header"
    values = ["${var.merchant_logo_domain}"]
  }
}
