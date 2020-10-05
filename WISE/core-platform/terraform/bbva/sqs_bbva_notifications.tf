resource "aws_sqs_queue" "bbva_notifications" {
  name       = "${module.naming.aws_sqs_queue}-primary"
  fifo_queue = "${var.bbva_sqs_fifo_queue}"

  delay_seconds              = "${var.bbva_sqs_delay_seconds}"
  receive_wait_time_seconds  = "${var.bbva_sqs_receive_wait_time_seconds}"
  visibility_timeout_seconds = "${var.bbva_sqs_visibility_timeout_seconds}"
  message_retention_seconds  = "${var.bbva_sqs_message_retention_seconds}"

  max_message_size            = "${var.bbva_sqs_max_message_size}"
  content_based_deduplication = "${var.bbva_sqs_content_based_deduplication}"

  kms_master_key_id                 = "${aws_kms_alias.bbva_sqs.target_key_arn}"
  kms_data_key_reuse_period_seconds = "${var.bbva_sqs_kms_data_key_reuse_period_seconds}"

  policy = ""

  redrive_policy = <<EOF
{
  "deadLetterTargetArn": "${aws_sqs_queue.bbva_dead_letter_queue.arn}",
  "maxReceiveCount": 2
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_sqs_queue}-primary"
    Team        = "${var.team}"
  }
}

resource "aws_sqs_queue_policy" "bbva_notifications" {
  queue_url = "${aws_sqs_queue.bbva_notifications.id}"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement": [
    {
      "Sid": "First",
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "arn:aws:iam::341687771823:role/p-${var.bbva_iam_role_env}-us-east-1-accounts-role",
          "arn:aws:iam::341687771823:role/p-${var.bbva_iam_role_env}-us-east-1-cards-role",
          "arn:aws:iam::341687771823:role/p-${var.bbva_iam_role_env}-us-east-1-events-role"
        ]
      },
      "Action": [
        "sqs:GetQueueAttributes",
        "sqs:GetQueueUrl",
        "sqs:SendMessage"
      ],
      "Resource": "${aws_sqs_queue.bbva_notifications.arn}"
    }
  ]
}
POLICY
}

resource "aws_cloudwatch_metric_alarm" "bbva_notifications_sqs_oldest_message_alarm" {
  alarm_name          = "${module.naming.aws_lambda_function}-sqs-primary-oldest-message"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "ApproximateAgeOfOldestMessage"
  namespace           = "AWS/SQS"
  period              = "60"
  statistic           = "Average"
  threshold           = "14400"                                                           # 4 minutes

  dimensions {
    QueueName = "${module.naming.aws_sqs_queue}-primary"
  }

  alarm_actions = [
    "${var.sns_critical_topic}",
  ]

  ok_actions = [
    "${var.sns_critical_topic}",
  ]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_lambda_function}-sqs-primary-oldest-message"
    Team        = "${var.team}"
  }
}

# resource "null_resource" "register_bbva_notifications" {
#   provisioner "local-exec" {
#     command    = "go run $SUBSCRIBE_APP"
#     on_failure = "fail"


#     environment = {
#       BBVA_APP_ENV                                 = "${data.aws_ssm_parameter.bbva_app_env.value}"
#       BBVA_APP_ID                                  = "${data.aws_ssm_parameter.bbva_app_id.value}"
#       BBVA_APP_NAME                                = "${data.aws_ssm_parameter.bbva_app_name.value}"
#       BBVA_APP_SECRET                              = "${data.aws_ssm_parameter.bbva_app_secret.value}"
#       DB_NAME                                      = "${data.aws_ssm_parameter.rds_db_name.value}"
#       DB_PASSWORD                                  = "${data.aws_ssm_parameter.rds_password.value}"
#       DB_READ_ENDPOINT                             = "${data.aws_ssm_parameter.rds_read_endpoint.value}"
#       DB_READ_PORT                                 = "${data.aws_ssm_parameter.rds_port.value}"
#       DB_USER                                      = "${data.aws_ssm_parameter.rds_user_name.value}"
#       DB_WRITE_ENDPOINT                            = "${data.aws_ssm_parameter.rds_master_endpoint.value}"
#       DB_WRITE_PORT                                = "${data.aws_ssm_parameter.rds_port.value}"
#       KINESIS_FIREHOSE_NAME_BANK_NOTIFICATION_NAME = "${var.environment}-logging-kinesis-firehose"
#       KINESIS_FIREHOSE_REGION_BANK_NOTIFICATION    = "${var.aws_region}"
#       SQS_FIFO_URL_ACTIVITY                        = "https://sqs.${var.aws_region}.amazonaws.com/${data.aws_caller_identity.account.account_id}/${var.environment}-client-api-banking-notifications"
#       SQS_FIFO_URL_BBVA_NOTIFICATION               = "${aws_sqs_queue.bbva_notifications.id}"
#       SUBSCRIBE_APP                                = "${var.bbva_subscribe_card_transactions_script}"
#     }
#   }


#   triggers = {
#     always_run = "${timestamp()}"
#   }


#   depends_on = [
#     "aws_sqs_queue.bbva_notifications",
#   ]
# }

