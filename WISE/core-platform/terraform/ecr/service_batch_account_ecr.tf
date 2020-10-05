resource "aws_ecr_repository" "batch_account" {
  name = "${module.naming.aws_ecr_repository}-batch-account"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-batch-account"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "batch_account" {
  repository = "${module.naming.aws_ecr_repository}-batch-account"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.batch_account",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "batch_account" {
  repository = "${module.naming.aws_ecr_repository}-batch-account"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.batch_account",
  ]
}

output "batch_account_repo_name" {
  value = "${aws_ecr_repository.batch_account.name}"
}
