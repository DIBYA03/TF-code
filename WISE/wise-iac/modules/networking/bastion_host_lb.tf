resource "aws_lb" "bastion_host" {
  count              = "${var.enable_bastion_host ? 1 : 0}"
  name               = "${module.naming.aws_lb}-bst-hst"
  internal           = true
  load_balancer_type = "network"

  enable_cross_zone_load_balancing = true
  subnets                          = ["${aws_subnet.app_subnets.*.id}"]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${module.naming.aws_lb}-bst-hst"
    Team        = "${var.team}"
  }
}

resource "aws_lb_target_group" "bastion_host" {
  count = "${var.enable_bastion_host ? 1 : 0}"
  name  = "${module.naming.aws_alb_target_group}-bst-hst"

  vpc_id   = "${aws_vpc.main.id}"
  port     = "${var.bastion_host_port}"
  protocol = "TCP"

  target_type = "instance"

  health_check {
    healthy_threshold   = 3
    unhealthy_threshold = 3
    interval            = 30
    port                = "traffic-port"
    protocol            = "TCP"
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-bst-hst"
    Team        = "${var.team}"
  }
}

resource "aws_lb_listener" "bastion_host" {
  count             = "${var.enable_bastion_host ? 1 : 0}"
  load_balancer_arn = "${aws_lb.bastion_host.arn}"
  port              = "${var.bastion_host_port}"
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.bastion_host.arn}"
  }
}
