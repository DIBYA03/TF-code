resource "aws_ecs_cluster" "default" {
  name = "${module.naming.aws_ecs_cluster}-default"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_cluster}-default"
    Team        = "${var.team}"
  }
}
