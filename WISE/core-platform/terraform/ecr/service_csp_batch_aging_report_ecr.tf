resource "aws_ecr_repository" "csp_batch_aging_report" {
  name = "${module.naming.aws_ecr_repository}-csp-batch-aging-report"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-csp-batch-aging-report"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "csp_batch_aging_report" {
  repository = "${module.naming.aws_ecr_repository}-csp-batch-aging-report"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_batch_aging_report",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "csp_batch_aging_report" {
  repository = "${module.naming.aws_ecr_repository}-csp-batch-aging-report"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_batch_aging_report",
  ]
}

output "csp_batch_aging_report_repo_name" {
  value = "${aws_ecr_repository.csp_batch_aging_report.name}"
}
