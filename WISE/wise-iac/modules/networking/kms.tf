resource "aws_kms_key" "default" {
  description         = "default kms key for ${var.environment}"
  key_usage           = "ENCRYPT_DECRYPT"
  enable_key_rotation = true

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Id": "key-default-1",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::${data.aws_caller_identity.account.account_id}:root"
      },
      "Action": "kms:*",
      "Resource": "*"
    },
    {
      "Sid": "Allow VPC Flow Logs to use the key",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "delivery.logs.amazonaws.com"
        ]
      },
      "Action": [
       "kms:Encrypt",
       "kms:Decrypt",
       "kms:ReEncrypt*",
       "kms:GenerateDataKey*",
       "kms:DescribeKey"
      ],
      "Resource": "*"
    },
    {
      "Sid": "Allow BastionHost for EBS encryption if bastion host exists",
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "arn:aws:iam::${data.aws_caller_identity.account.account_id}:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
          "arn:aws:sts::${data.aws_caller_identity.account.account_id}:assumed-role/AWSServiceRoleForAutoScaling/AutoScaling"
        ]
      },
      "Action": [
        "kms:Encrypt",
        "kms:Decrypt",
        "kms:ReEncrypt*",
        "kms:GenerateDataKey*",
        "kms:DescribeKey"
      ],
      "Resource": "*"
    },
    {
      "Sid": "Allow attachment of persistent resources",
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "arn:aws:iam::${data.aws_caller_identity.account.account_id}:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
          "arn:aws:sts::${data.aws_caller_identity.account.account_id}:assumed-role/AWSServiceRoleForAutoScaling/AutoScaling"
        ]
      },
      "Action": [
        "kms:CreateGrant"
      ],
      "Resource": "*",
      "Condition": {
        "Bool": {
          "kms:GrantIsForAWSResource": true
        }
      }
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-${var.application}-default"
    Team        = "${var.team}"
  }
}

resource "aws_kms_alias" "default" {
  name          = "${module.naming.aws_kms_alias}"
  target_key_id = "${aws_kms_key.default.key_id}"
}

# This is needed for BBVA to send messages to our encypted SQS queues
resource "aws_kms_grant" "kms_grant_cross_account_resources" {
  count             = "${length(var.kms_grant_cross_account_resources)}"
  name              = "${module.naming.aws_kms_grant}-cross-account-resource-${count.index}"
  key_id            = "${aws_kms_key.default.key_id}"
  grantee_principal = "${var.kms_grant_cross_account_resources[count.index]}"

  operations = [
    "Encrypt",
    "Decrypt",
    "DescribeKey",
    "GenerateDataKey",
  ]
}
