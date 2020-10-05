resource "aws_s3_bucket" "cloudfront" {
  bucket = "${module.naming.aws_s3_bucket}-cloudfront"
  acl    = "private"

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "${aws_kms_alias.cloudfront_bucket.target_key_arn}"
        sse_algorithm     = "aws:kms"
      }
    }
  }

  replication_configuration {
    role = "${module.cloudfront_s3_replication.replication_iam_role_arn}"

    rules {
      id     = "cross-region-replication"
      status = "Enabled"

      destination {
        bucket             = "${module.cloudfront_s3_replication.s3_bucket_arn}"
        storage_class      = "STANDARD"
        replica_kms_key_id = "${module.cloudfront_s3_replication.kms_key_arn}"
      }

      source_selection_criteria {
        sse_kms_encrypted_objects {
          enabled = true
        }
      }
    }
  }

  tags {
    Application = "${var.application}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_s3_bucket}-cloudfront"
    Team        = "${var.team}"
  }
}

resource "aws_s3_bucket_public_access_block" "cloudfront" {
  bucket = "${aws_s3_bucket.cloudfront.id}"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

module "cloudfront_s3_replication" {
  source = "git@github.com:wiseco/terraform-module-s3-replication.git"

  application                         = "${var.application}"
  aws_region                          = "us-east-1"
  aws_profile                         = "${var.aws_profile}"
  component                           = "${var.component}"
  environment                         = "${var.environment}"
  team                                = "${var.team}"
  s3_suffix                           = "cloudfront"
  source_s3_bucket_arn                = "arn:aws:s3:::${module.naming.aws_s3_bucket}-cloudfront"
  source_s3_bucket_encryption_key_arn = "${aws_kms_alias.cloudfront_bucket.target_key_arn}"
}
