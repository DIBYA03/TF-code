resource "aws_security_group" "services" {
  name        = "${module.naming.aws_security_group}-services"
  description = "security group for ${var.environment_name} services"

  vpc_id = "${var.vpc_id}"

  ingress {
    description = "http inbound for redirect to https"
    from_port   = 80
    to_port     = 80
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "https inbound"
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS inbound from stripe"
    from_port   = 4433
    to_port     = 4433
    protocol    = "TCP"

    cidr_blocks = ["${var.stripe_webhook_ip_list}"]
  }

  ingress {
    description = "HTTPS inbound from hello sign"
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"

    cidr_blocks = ["${var.hello_sign_webhook_ip_list}"]
  }

  egress {
    from_port   = "${var.services_container_port}"
    to_port     = "${var.services_container_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_security_group}-services"
    Team        = "${var.team}"
  }
}

resource "aws_alb" "services" {
  name               = "${module.naming.aws_alb}-services"
  load_balancer_type = "application"

  enable_cross_zone_load_balancing = true
  subnets                          = ["${var.public_subnet_ids}"]
  security_groups                  = ["${aws_security_group.services.id}"]

  tags {
    Application   = "${var.application}"
    Component     = "${var.component}"
    Environment   = "${var.environment_name}"
    Name          = "${module.naming.aws_alb}-services"
    Team          = "${var.team}"
    AddOWASP10WAF = "${var.merchant_logo_add_owasp10_waf}"
  }
}

resource "aws_alb_listener" "services_http" {
  load_balancer_arn = "${aws_alb.services.arn}"
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_alb_listener" "services_https" {
  load_balancer_arn = "${aws_alb.services.arn}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${module.services_cert.arn}"

  # This is because rules will be added at the service that will be using this
  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      status_code  = "403"
    }
  }
}

resource "aws_alb_listener" "services_stripe_https" {
  load_balancer_arn = "${aws_alb.services.arn}"
  port              = "4433"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${module.services_cert.arn}"

  # This is because rules will be added at the service that will be using this
  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      status_code  = "403"
    }
  }
}

module "services_cert" {
  source = "./modules/terraform-module-acm-certificate"

  application            = "${var.application}"
  aws_profile            = "${var.aws_profile}"
  route53_aws_profile    = "${var.aws_profile}"
  aws_region             = "${var.aws_region}"
  component              = "${var.component}"
  domain_name            = "*.wise.us"
  environment            = "${var.environment}"
  route53_aws_profile    = "${var.public_route53_account_profile}"
  route53_hosted_zone_id = "${var.public_hosted_zone_id}"
  team                   = "${var.team}"
}
