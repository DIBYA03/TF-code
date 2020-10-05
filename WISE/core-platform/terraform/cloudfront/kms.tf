resource "aws_kms_key" "cloudfront_bucket" {
  description         = "${var.environment_name} KMS Key for S3 cloudfront bucket"
  enable_key_rotation = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${var.environment}-s3-cloudfront-bucket"
    Team        = "${var.team}"
  }
}

resource "aws_kms_alias" "cloudfront_bucket" {
  name          = "${module.naming.aws_kms_alias}-s3-cloudfront"
  target_key_id = "${aws_kms_key.cloudfront_bucket.key_id}"
}
