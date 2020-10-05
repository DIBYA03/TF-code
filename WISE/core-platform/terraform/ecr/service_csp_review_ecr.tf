resource "aws_ecr_repository" "csp_review" {
  name = "${module.naming.aws_ecr_repository}-csp-review"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-csp-review"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "csp_review" {
  repository = "${module.naming.aws_ecr_repository}-csp-review"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_review",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "csp_review" {
  repository = "${module.naming.aws_ecr_repository}-csp-review"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_review",
  ]
}

output "csp_review_repo_name" {
  value = "${aws_ecr_repository.csp_review.name}"
}
