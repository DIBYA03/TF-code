resource "aws_s3_bucket" "backup_bucket" {
  bucket = "${module.naming.aws_s3_bucket}-backups"
  acl    = "private"

  # NEED TO STILL ADD BUCKET REPLICATION

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "${aws_kms_key.default.arn}"
        sse_algorithm     = "aws:kms"
      }
    }
  }
  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${module.naming.aws_s3_bucket}-backups"
    Team        = "${var.team}"
  }
}

resource "aws_s3_bucket_public_access_block" "backup_bucket" {
  bucket = "${aws_s3_bucket.backup_bucket.id}"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
