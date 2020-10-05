data "aws_kms_alias" "env_default" {
  name = "${var.default_kms_alias}"
}

resource "aws_kms_key" "bbva_sqs" {
  description         = "${var.environment_name} KMS Key for BBVA SQS"
  enable_key_rotation = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${var.environment}-bbva-sqs"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role" "bbva_sqs" {
  name = "${module.naming.aws_iam_role}-bbva-sqs"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
        "Action": "sts:AssumeRole",
        "Principal": {
            "Service": [
                "lambda.amazonaws.com",
                "sqs.amazonaws.com"
            ]
        },
        "Effect": "Allow",
        "Sid": ""
        }
    ]
}
EOF
}

resource "aws_kms_alias" "bbva_sqs" {
  name          = "${module.naming.aws_kms_alias}-bbva-sqs"
  target_key_id = "${aws_kms_key.bbva_sqs.key_id}"
}

# This is needed for BBVA to send messages to our encypted SQS queues
resource "aws_kms_grant" "bbva_partner" {
  name              = "${module.naming.aws_kms_grant}-bbva-partner-qs"
  key_id            = "${aws_kms_key.bbva_sqs.key_id}"
  grantee_principal = "arn:aws:iam::341687771823:root"

  operations = [
    "Encrypt",
    "Decrypt",
    "GenerateDataKey",
  ]
}

resource "aws_kms_grant" "bbva_sqs" {
  name              = "${module.naming.aws_kms_grant}-bbva-sqs"
  key_id            = "${aws_kms_key.bbva_sqs.key_id}"
  grantee_principal = "${aws_iam_role.bbva_sqs.arn}"

  operations = [
    "Encrypt",
    "Decrypt",
  ]
}
