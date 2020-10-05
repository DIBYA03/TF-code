resource "aws_ecr_repository" "segment_analytics" {
  name = "${module.naming.aws_ecr_repository}-segment-analytics"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-segment-analytics"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "segment_analytics" {
  repository = "${module.naming.aws_ecr_repository}-segment-analytics"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.segment_analytics",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "segment_analytics" {
  repository = "${module.naming.aws_ecr_repository}-segment-analytics"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.segment_analytics",
  ]
}

output "segment_analytics_repo_name" {
  value = "${aws_ecr_repository.segment_analytics.name}"
}
