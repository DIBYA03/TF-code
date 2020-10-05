resource "aws_iam_role" "bbva_sns_connector_execution_role" {
  name = "${module.naming.aws_iam_role}-bbva-sns-connector-execution-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
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
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_iam_role}-bbva-sns-connector-execution-role"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy" "bbva_sns_connector_execution_role" {
  name = "${module.naming.aws_iam_policy}-default-execution"
  role = "${aws_iam_role.bbva_sns_connector_execution_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    },
    {
      "Sid": "AllowAccessToBBVASQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "sqs:ChangeMessageVisibility*",
        "sqs:DeleteMessage*",
        "sqs:ReceiveMessage"
      ],
      "Resource": [
        "${aws_kms_alias.bbva_sqs.target_key_arn}",
        "${aws_sqs_queue.bbva_notifications.arn}"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "kms:GenerateDataKey",
        "kms:Decrypt"
      ],
      "Resource": "${aws_kms_alias.bbva_sqs.target_key_arn}"
      }, {
      "Effect": "Allow",
      "Action": [
        "sns:Publish"
      ],
      "Resource": "${aws_sns_topic.bbva_notifications.arn}"
    }
  ]
}
EOF
}
