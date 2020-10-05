resource "aws_ecr_repository" "batch_monthly_interest" {
  name = "${module.naming.aws_ecr_repository}-batch-monthly-interest"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-batch-monthly-interest"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "batch_monthly_interest" {
  repository = "${module.naming.aws_ecr_repository}-batch-monthly-interest"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.batch_monthly_interest",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "batch_monthly_interest" {
  repository = "${module.naming.aws_ecr_repository}-batch-monthly-interest"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.batch_monthly_interest",
  ]
}

output "batch_monthly_interest_repo_name" {
  value = "${aws_ecr_repository.batch_monthly_interest.name}"
}
