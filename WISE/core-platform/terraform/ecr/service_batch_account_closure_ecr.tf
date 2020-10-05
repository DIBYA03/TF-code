resource "aws_ecr_repository" "account_closure" {
  name = "${module.naming.aws_ecr_repository}-batch-account-closure"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-batch-account-closure"
    Team        = "${var.team}"
  }

}


resource "aws_ecr_lifecycle_policy" "account_closure" {
  repository = "${module.naming.aws_ecr_repository}-batch-account-closure"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.account_closure",
  ]
}

# allow access from only wise analytics
resource "aws_ecr_repository_policy" "account_closure" {
  repository = "${module.naming.aws_ecr_repository}-batch-account-closure"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.account_closure",
  ]
}

output "account_closure_repo_name" {
  value = "${aws_ecr_repository.account_closure.name}"
}
