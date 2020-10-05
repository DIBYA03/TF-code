resource "aws_ecr_repository" "merchant_logo" {
  name = "${module.naming.aws_ecr_repository}-merchant-logo"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-merchant-logo"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "merchant_logo" {
  repository = "${module.naming.aws_ecr_repository}-merchant-logo"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.merchant_logo",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "merchant_logo" {
  repository = "${module.naming.aws_ecr_repository}-merchant-logo"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.merchant_logo",
  ]
}

output "merchant_logo_repo_name" {
  value = "${aws_ecr_repository.merchant_logo.name}"
}
