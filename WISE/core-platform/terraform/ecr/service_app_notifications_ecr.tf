resource "aws_ecr_repository" "app_notifications" {
  name = "${module.naming.aws_ecr_repository}-app-notifications"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-app-notifications"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "app_notifications" {
  repository = "${module.naming.aws_ecr_repository}-app-notifications"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.app_notifications",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "app_notifications" {
  repository = "${module.naming.aws_ecr_repository}-app-notifications"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.app_notifications",
  ]
}

output "app_notificaitons_repo_name" {
  value = "${aws_ecr_repository.app_notifications.name}"
}
