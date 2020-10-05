resource "aws_ecr_repository" "csp_document_upload" {
  name = "${module.naming.aws_ecr_repository}-csp-document-upload"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-csp-document-upload"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "csp_document_upload" {
  repository = "${module.naming.aws_ecr_repository}-csp-document-upload"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_document_upload",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "csp_document_upload" {
  repository = "${module.naming.aws_ecr_repository}-csp-document-upload"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.csp_document_upload",
  ]
}

output "csp_document_repo_name" {
  value = "${aws_ecr_repository.csp_document_upload.name}"
}
