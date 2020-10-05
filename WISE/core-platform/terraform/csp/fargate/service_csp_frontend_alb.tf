module "csp_frontend_route_53" {
  source = "./modules/route53"

  aws_profile            = "${var.public_route53_account_profile}"
  aws_region             = "${var.aws_region}"
  domain_name            = "${var.csp_frontend_domain}"
  route53_hosted_zone_id = "${var.route53_private_hosted_zone_id}"
  resource_alias_name    = "${aws_alb.csp_frontend.dns_name}"
  resource_alias_zone_id = "${aws_alb.csp_frontend.zone_id}"
}

resource "aws_alb" "csp_frontend" {
  name               = "${module.naming.aws_alb}-swh"
  load_balancer_type = "application"
  internal           = true

  enable_cross_zone_load_balancing = true
  subnets                          = ["${var.app_subnet_ids}"]
  security_groups                  = ["${aws_security_group.csp_frontend_ecs_alb.id}"]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_alb}-swh"
    Team        = "${var.team}"
  }
}

resource "aws_alb_target_group" "csp_frontend" {
  name_prefix = "cspfnt"
  port        = "${var.csp_frontend_container_port}"
  protocol    = "HTTP"
  vpc_id      = "${var.vpc_id}"
  target_type = "ip"

  health_check = {
    enabled = true
    path    = "/healthcheck.html"
    port    = "traffic-port"
    matcher = 200
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_alb_target_group}-csp-frontend"
    Team        = "${var.team}"
  }
}

resource "aws_alb_listener" "csp_frontend_https" {
  load_balancer_arn = "${aws_alb.csp_frontend.arn}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${module.csp_frontend.arn}"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_alb_target_group.csp_frontend.arn}"
  }
}
