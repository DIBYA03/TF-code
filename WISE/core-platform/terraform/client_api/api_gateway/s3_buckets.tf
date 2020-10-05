resource "aws_s3_bucket" "documents" {
  bucket = "${module.naming.aws_s3_bucket}-documents"
  acl    = "private"

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }

  # This is needed for s3 presigned URL for uploading documents
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT","GET"]
    allowed_origins = ["*"]
    max_age_seconds = 900
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "${aws_kms_alias.documents_bucket.target_key_arn}"
        sse_algorithm     = "aws:kms"
      }
    }
  }

  replication_configuration {
    role = "${module.documents_s3_replication.replication_iam_role_arn}"

    rules {
      id     = "cross-region-replication"
      status = "Enabled"

      destination {
        bucket             = "${module.documents_s3_replication.s3_bucket_arn}"
        storage_class      = "STANDARD"
        replica_kms_key_id = "${module.documents_s3_replication.kms_key_arn}"
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
    Name        = "${module.naming.aws_s3_bucket}-documents"
    Team        = "${var.team}"
  }
}

resource "aws_s3_bucket_public_access_block" "documents" {
  bucket = "${aws_s3_bucket.documents.id}"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

module "documents_s3_replication" {
  source = "git@github.com:wiseco/terraform-module-s3-replication.git"

  application                         = "${var.application}"
  aws_profile                         = "${var.aws_profile}"
  aws_region                          = "us-east-1"
  aws_profile                         = "${var.aws_profile}"
  component                           = "${var.component}"
  environment                         = "${var.environment}"
  team                                = "${var.team}"
  s3_suffix                           = "documents"
  source_s3_bucket_arn                = "arn:aws:s3:::${module.naming.aws_s3_bucket}-documents"
  source_s3_bucket_encryption_key_arn = "${aws_kms_alias.documents_bucket.target_key_arn}"
}
