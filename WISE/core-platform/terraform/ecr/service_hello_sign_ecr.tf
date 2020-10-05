resource "aws_ecr_repository" "hello_sign" {
  name = "${module.naming.aws_ecr_repository}-hello-sign"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-hello-sign"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "hello_sign" {
  repository = "${module.naming.aws_ecr_repository}-hello-sign"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.hello_sign",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "hello_sign" {
  repository = "${module.naming.aws_ecr_repository}-hello-sign"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.hello_sign",
  ]
}

output "hello_sign_repo_name" {
  value = "${aws_ecr_repository.hello_sign.name}"
}
