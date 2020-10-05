resource "aws_ecr_repository" "stripe_webhook" {
  name = "${module.naming.aws_ecr_repository}-stripe-webhook"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-stripe-webhook"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "stripe_webhook" {
  repository = "${module.naming.aws_ecr_repository}-stripe-webhook"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.stripe_webhook",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "stripe_webhook" {
  repository = "${module.naming.aws_ecr_repository}-stripe-webhook"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.stripe_webhook",
  ]
}

output "stripe_webhook_repo_name" {
  value = "${aws_ecr_repository.stripe_webhook.name}"
}
