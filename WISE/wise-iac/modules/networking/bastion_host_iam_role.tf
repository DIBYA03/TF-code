resource "aws_iam_role" "bastion_host" {
  count = "${var.enable_bastion_host ? 1 : 0}"
  name  = "${module.naming.aws_iam_policy}-bastion-host"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "ec2.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${module.naming.aws_iam_policy}-bastion-host"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "bastion_host" {
  count = "${var.enable_bastion_host ? 1 : 0}"
  name  = "${module.naming.aws_iam_policy}-default"
  role  = "${aws_iam_role.bastion_host.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudwatch:PutMetricData"
      ],
      "Resource": "*"
    },
    {
      "Sid": "ListBackupBucket",
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": [
        "${aws_s3_bucket.backup_bucket.arn}"
      ]
    },
    {
      "Sid": "GetSaltstackZip",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "s3:GetObject",
        "s3:GetObjectTagging",
        "s3:GetObjectVersion",
        "s3:GetObjectVersionTagging",
        "s3:ListBucket"
      ],
      "Resource": [
        "${aws_kms_alias.default.target_key_arn}",
        "${aws_s3_bucket.backup_bucket.arn}/${var.bastion_host_salstack_s3_object_prefix}/${var.bastion_host_salstack_s3_object_name}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "bastion_host" {
  count = "${var.enable_bastion_host ? 1 : 0}"
  name  = "${module.naming.aws_iam_instance_profile}-bastion-host"
  role  = "${aws_iam_role.bastion_host.name}"
}
