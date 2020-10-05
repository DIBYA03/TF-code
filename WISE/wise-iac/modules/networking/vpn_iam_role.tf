resource "aws_iam_role" "vpn" {
  count = "${var.enable_vpn == "true" ? 1 : 0}"
  name  = "${module.naming.aws_iam_policy}-vpn"

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
    Name        = "${module.naming.aws_iam_policy}-vpn"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "vpn" {
  count = "${var.enable_vpn == "true" ? 1 : 0}"

  name = "${module.naming.aws_iam_policy}-vpn"
  role = "${aws_iam_role.vpn.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudwatch:PutMetricData",
        "ec2:DescribeTags"
      ],
      "Resource": "*"
    },
    {
      "Sid": "AllowEc2ToAttachPublicIP",
      "Effect": "Allow",
      "Action": [
        "ec2:AssociateAddress",
        "ec2:DescribeAddresses",
        "ec2:ModifyInstanceAttribute"
      ],
      "Resource": "*"
    },
    {
      "Sid": "AllowEC2AttachEncryptedVolume",
      "Effect": "Allow",
      "Action": [
        "ec2:AttachVolume",
        "kms:CreateGrant",
        "kms:Decrypt",
        "kms:Describe*",
        "kms:Encrypt",
        "kms:GenerateDataKey*",
        "kms:ListKeys",
        "kms:ReEncrypt*"
      ],
      "Resource": [
        "arn:aws:ec2:${var.aws_region}:*:instance/*",
        "${aws_ebs_volume.vpn.arn}",
        "${aws_kms_key.default.arn}",
        "${aws_kms_alias.default.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "vpn" {
  count = "${var.enable_vpn == "true" ? 1 : 0}"
  name  = "${module.naming.aws_iam_instance_profile}-vpn"
  role  = "${aws_iam_role.vpn.name}"
}
