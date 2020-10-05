resource "aws_ecr_repository" "app_payments" {
  name = "${module.naming.aws_ecr_repository}-app-payments"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-app-payments"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "app_payments" {
  repository = "${module.naming.aws_ecr_repository}-app-payments"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.app_payments",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "app_payments" {
  repository = "${module.naming.aws_ecr_repository}-app-payments"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.app_payments",
  ]
}

output "app_payments_repo_name" {
  value = "${aws_ecr_repository.app_payments.name}"
}
