resource "aws_ecr_repository" "batch_analytics" {
  name = "${module.naming.aws_ecr_repository}-batch-analytics"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-batch-analytics"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "batch_analytics" {
  repository = "${module.naming.aws_ecr_repository}-batch-analytics"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.batch_analytics",
  ]
}

# allow access from only wise analytics
resource "aws_ecr_repository_policy" "batch_analytics" {
  repository = "${module.naming.aws_ecr_repository}-batch-analytics"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.batch_analytics",
  ]
}

output "batch_analytics_repo_name" {
  value = "${aws_ecr_repository.batch_analytics.name}"
}
