resource "aws_ecr_repository" "bbva_notification" {
  name = "${module.naming.aws_ecr_repository}-bbva-notification"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-bbva-notification"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "bbva_notification" {
  repository = "${module.naming.aws_ecr_repository}-bbva-notification"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.bbva_notification",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "bbva_notification" {
  repository = "${module.naming.aws_ecr_repository}-bbva-notification"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.bbva_notification",
  ]
}

output "bbva_notification_repo_name" {
  value = "${aws_ecr_repository.bbva_notification.name}"
}
