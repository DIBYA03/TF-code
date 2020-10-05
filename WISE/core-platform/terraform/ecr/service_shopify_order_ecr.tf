resource "aws_ecr_repository" "shopify_order" {
  name = "${module.naming.aws_ecr_repository}-shopify-order"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_ecr_repository}-shopify-order"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "shopify_order" {
  repository = "${module.naming.aws_ecr_repository}-shopify-order"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.shopify_order",
  ]
}

# allow access from only wise account
resource "aws_ecr_repository_policy" "shopify_order" {
  repository = "${module.naming.aws_ecr_repository}-shopify-order"

  policy = "${data.template_file.default_repo_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.shopify_order",
  ]
}

output "shopify_order_repo_name" {
  value = "${aws_ecr_repository.shopify_order.name}"
}
