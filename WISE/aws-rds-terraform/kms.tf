resource "aws_kms_key" "rds_default" {
  description         = "KMS Key for ${terraform.workspace} RDS"
  enable_key_rotation = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${terraform.workspace}-rds-kms-key"
    Team        = "${var.team}"
  }
}

resource "aws_kms_alias" "rds_default" {
  name          = "${module.naming.aws_kms_alias}-default"
  target_key_id = "${aws_kms_key.rds_default.key_id}"
}

# resource "aws_iam_role" "rds_default" {
#   name = "${module.naming.aws_iam_role}-default"


#   assume_role_policy = <<EOF
# {
#     "Version": "2012-10-17",
#     "Statement": [
#         {
#         "Action": "sts:AssumeRole",
#         "Principal": {
#             "Service": "rds.amazonaws.com"
#         },
#         "Effect": "Allow",
#         "Sid": ""
#         }
#     ]
# }
# EOF
# }


# resource "aws_kms_grant" "rds_default" {
#   name              = "${terraform.workspace}-iam-lambdas-to-s3-documents"
#   key_id            = "${aws_kms_key.rds_default.key_id}"
#   grantee_principal = "${aws_iam_role.rds_default.arn}"


#   operations = [
#     "Encrypt",
#     "Decrypt",
#   ]
# }

