resource "aws_ecr_repository" "batch_transaction" {
  name = "${module.naming.aws_ecr_repository}-batch-transaction"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-batch-transaction"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "batch_transaction" {
  repository = "${module.naming.aws_ecr_repository}-batch-transaction"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.batch_transaction",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "batch_transaction" {
  repository = "${module.naming.aws_ecr_repository}-batch-transaction"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.batch_transaction",
  ]
}

output "batch_transaction_repo_name" {
  value = "${aws_ecr_repository.batch_transaction.name}"
}
