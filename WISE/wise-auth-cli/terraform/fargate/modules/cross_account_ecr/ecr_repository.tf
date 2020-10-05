resource "aws_ecr_repository" "aws_auth_cli_wise" {
  name = "${module.naming.aws_ecr_repository}-wise"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecr_repository}-wise"
    Team        = "${var.team}"
  }

  provider = "aws.${var.provider_name}"
}

resource "aws_ecr_lifecycle_policy" "aws_auth_cli_wise" {
  repository = "${module.naming.aws_ecr_repository}-wise"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.aws_auth_cli_wise",
  ]

  provider = "aws.${var.provider_name}"
}

output "aws_auth_cli_wise_deploy_command" {
  value = <<EOF

cd ../../cmd/docker/auth-cli-wise;
AWS_DEFAULT_PROFILE=${var.aws_profile} \
  AWS_DEFAULT_REGION=${var.aws_region} \
  ECS_ENV=${var.environment} \
  BUILD_TAG=latest \
  ECR_IMAGE=${aws_ecr_repository.aws_auth_cli_wise.repository_url} \
  make all;
cd -
EOF
}
