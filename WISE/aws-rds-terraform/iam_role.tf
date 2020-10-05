resource "aws_iam_role" "rds_backup_lambda" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  name = "${module.naming.aws_iam_role}-backup-lambda"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "lambda.amazonaws.com"
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
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_iam_role}-backup-lambda"
    Team        = "${var.team}"
  }
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "rds_backup_lambda_cw_logs" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  role       = "${aws_iam_role.rds_backup_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "rds_backup_lambda_vpc_access" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  role       = "${aws_iam_role.rds_backup_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "rds_backup_lambda_create_snapshot" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  name = "${module.naming.aws_iam_role_policy}-create-snapshot"
  role = "${aws_iam_role.rds_backup_lambda.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "rds:CreateDBSnapshot",
        "rds:DescribeDBSnapshots",
        "rds:DeleteDBSnapshot"
      ],
      "Effect": "Allow",
      "Resource": [
        "${aws_db_instance.backup_replica.arn}",
        "arn:aws:rds:${var.aws_region}:${data.aws_caller_identity.account.account_id}:snapshot:${aws_db_instance.backup_replica.id}*"
      ]
    }
  ]
}
EOF
}
