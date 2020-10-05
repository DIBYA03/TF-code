resource "aws_ecr_repository" "csp_batch_business" {
  name = "${module.naming.aws_ecr_repository}-csp-batch-business"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-csp-batch-business"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "csp_batch_business" {
  repository = "${module.naming.aws_ecr_repository}-csp-batch-business"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_batch_business",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "csp_batch_business" {
  repository = "${module.naming.aws_ecr_repository}-csp-batch-business"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_batch_business",
  ]
}

output "csp_batch_business_repo_name" {
  value = "${aws_ecr_repository.csp_batch_business.name}"
}
