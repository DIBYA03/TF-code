resource "aws_s3_bucket_object" "bbva_requeue_s3_object" {
  bucket = "${aws_s3_bucket.documents.id}"
  key    = "${var.s3_bbva_requeue_object}"
  source = "./specs/s3_objects/bbva-ready-for-review.pdf"
  etag   = "${filemd5("./specs/s3_objects/bbva-ready-for-review.pdf")}"
}

resource "aws_ssm_parameter" "bbva_requeue_s3_object" {
  name      = "/${var.environment}/dev/bbva/app_env"
  type      = "SecureString"
  key_id    = "${var.default_kms_alias}"
  value     = "${var.s3_bbva_requeue_object}"
  overwrite = true
}
