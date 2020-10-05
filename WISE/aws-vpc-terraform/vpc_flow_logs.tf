resource "aws_s3_bucket" "flow_logs" {
  count  = "${var.enable_flow_logs ? 1 : 0}"
  bucket = "${module.naming.aws_s3_bucket}-flow-logs"
  acl    = "private"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Id": "AWSLogDeliveryWrite20150319",
  "Statement": [
    {
      "Sid": "AWSLogDeliveryAclCheck",
      "Effect": "Allow",
      "Principal": {
        "Service": "delivery.logs.amazonaws.com"
      },
      "Action": "s3:GetBucketAcl",
      "Resource": "arn:aws:s3:::${module.naming.aws_s3_bucket}-flow-logs"
    },
    {
      "Sid": "AWSLogDeliveryWrite",
      "Effect": "Allow",
      "Principal": {
        "Service": "delivery.logs.amazonaws.com"
      },
      "Action": "s3:PutObject",
      "Resource": "arn:aws:s3:::${module.naming.aws_s3_bucket}-flow-logs/AWSLogs/${data.aws_caller_identity.account.account_id}/*",
      "Condition": {
        "StringEquals": {
          "s3:x-amz-acl": "bucket-owner-full-control"
        }
      }
    }
  ]
}
EOF

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
    Name        = "${module.naming.aws_s3_bucket}-flow-logs"
    Team        = "${var.team}"
  }
}

resource "aws_s3_bucket_public_access_block" "flow_logs" {
  count  = "${var.enable_flow_logs ? 1 : 0}"
  bucket = "${aws_s3_bucket.backup_bucket.id}"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_flow_log" "default" {
  count = "${var.enable_flow_logs ? 1 : 0}"

  vpc_id = "${aws_vpc.main.id}"

  log_destination      = "${aws_s3_bucket.flow_logs.arn}"
  log_destination_type = "s3"
  traffic_type         = "ALL"
}
