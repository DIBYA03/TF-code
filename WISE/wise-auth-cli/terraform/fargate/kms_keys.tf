data "aws_kms_alias" "default" {
  name = "${var.default_kms_alias}"
}
