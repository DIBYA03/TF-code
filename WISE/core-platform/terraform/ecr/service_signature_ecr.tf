resource "aws_ecr_repository" "signature" {
  name = "${module.naming.aws_ecr_repository}-signature"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-signature"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "signature" {
  repository = "${module.naming.aws_ecr_repository}-signature"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.signature",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "signature" {
  repository = "${module.naming.aws_ecr_repository}-signature"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.signature",
  ]
}

output "signature_repo_name" {
  value = "${aws_ecr_repository.signature.name}"
}
