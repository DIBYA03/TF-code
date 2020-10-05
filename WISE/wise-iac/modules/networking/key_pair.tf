resource "aws_key_pair" "default" {
  key_name   = "wise-${var.aws_region}-${var.environment}"
  public_key = "${var.default_ssh_key}"
}
