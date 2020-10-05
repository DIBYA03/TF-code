resource "aws_ecr_repository" "bbva_sns_connector" {
  name = "${module.naming.aws_ecr_repository}-bbva-sns-connector"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-bbva-sns-connector"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "bbva_sns_connector" {
  repository = "${module.naming.aws_ecr_repository}-bbva-sns-connector"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.bbva_sns_connector",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "bbva_sns_connector" {
  repository = "${module.naming.aws_ecr_repository}-bbva-sns-connector"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.bbva_sns_connector",
  ]
}

output "bbva_sns_connector_repo_name" {
  value = "${aws_ecr_repository.bbva_sns_connector.name}"
}
