resource "aws_ecr_repository" "aws_vpn_auth" {
  name = "${module.naming.aws_ecr_repository}-vpn"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecr_repository}-vpn"
    Team        = "${var.team}"
  }
}

resource "aws_ecr_lifecycle_policy" "aws_vpn_auth" {
  repository = "${module.naming.aws_ecr_repository}-vpn"

  policy = "${data.template_file.default_lifecycle_policy.rendered}"

  depends_on = [
    "aws_ecr_repository.aws_vpn_auth",
  ]
}

output "aws_vpn_auth_deploy_command" {
  value = <<EOF

cd ../../cmd/docker/aws-vpn-auth;
AWS_DEFAULT_PROFILE=${var.aws_profile} \
  AWS_DEFAULT_REGION=${var.aws_region} \
  ECS_ENV=${var.environment} \
  BUILD_TAG=${var.aws_vpn_auth_image_tag} \
  ECR_IMAGE=${aws_ecr_repository.aws_vpn_auth.repository_url} \
  make all;
cd -
EOF
}
