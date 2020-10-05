resource "aws_ecr_repository" "csp_batch_bank_review_report" {
  name = "${module.naming.aws_ecr_repository}-csp-batch-bank-review-report"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-csp-batch-bank-review-report"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "csp_batch_bank_review_report" {
  repository = "${module.naming.aws_ecr_repository}-csp-batch-bank-review-report"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_batch_bank_review_report",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "csp_batch_bank_review_report" {
  repository = "${module.naming.aws_ecr_repository}-csp-batch-bank-review-report"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_batch_bank_review_report",
  ]
}

output "csp_batch_bank_review_report_repo_name" {
  value = "${aws_ecr_repository.csp_batch_bank_review_report.name}"
}
