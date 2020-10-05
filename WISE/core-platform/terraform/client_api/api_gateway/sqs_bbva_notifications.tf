resource "aws_sqs_queue" "bbva_dead_letter_queue" {
  name       = "${module.naming.aws_sqs_queue}-bbva-notifications-dead-letter"
  fifo_queue = "${var.bbva_sqs_fifo_queue}"

  delay_seconds              = "${var.bbva_sqs_delay_seconds}"
  receive_wait_time_seconds  = "${var.bbva_sqs_receive_wait_time_seconds}"
  visibility_timeout_seconds = "${var.bbva_sqs_visibility_timeout_seconds}"
  message_retention_seconds  = "${var.bbva_sqs_dl_message_retention_seconds}"

  max_message_size            = "${var.bbva_sqs_max_message_size}"
  content_based_deduplication = "${var.bbva_sqs_dl_content_based_deduplication}"

  kms_master_key_id                 = "${aws_kms_alias.internal_sqs.target_key_arn}"
  kms_data_key_reuse_period_seconds = "${var.bbva_sqs_kms_data_key_reuse_period_seconds}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_sqs_queue}-bbva-notifications-dead-letter"
    Team        = "${var.team}"
  }
}

resource "aws_sqs_queue" "bbva_notifications" {
  name       = "${module.naming.aws_sqs_queue}-bbva-notifications"
  fifo_queue = "${var.bbva_sqs_fifo_queue}"

  delay_seconds              = "${var.bbva_sqs_delay_seconds}"
  receive_wait_time_seconds  = "${var.bbva_sqs_receive_wait_time_seconds}"
  visibility_timeout_seconds = "${var.bbva_sqs_visibility_timeout_seconds}"
  message_retention_seconds  = "${var.bbva_sqs_message_retention_seconds}"

  max_message_size            = "${var.bbva_sqs_max_message_size}"
  content_based_deduplication = "${var.bbva_sqs_content_based_deduplication}"

  kms_master_key_id                 = "${aws_kms_alias.internal_sqs.target_key_arn}"
  kms_data_key_reuse_period_seconds = "${var.bbva_sqs_kms_data_key_reuse_period_seconds}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Id": "arn:aws:sqs:${var.aws_region}:${data.aws_caller_identity.account.account_id}:${module.naming.aws_sqs_queue}-bbva-notifications/SQSDefaultPolicy",
  "Statement": [
    {
      "Sid": "AllowSNS",
      "Effect": "Allow",
      "Principal": "*",
      "Action": [
        "SQS:SendMessage",
        "SQS:ReceiveMessage",
        "SQS:DeleteMessage"
      ],
      "Resource": "arn:aws:sqs:${var.aws_region}:${data.aws_caller_identity.account.account_id}:${module.naming.aws_sqs_queue}-bbva-notifications",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "arn:aws:sns:${var.aws_region}:${data.aws_caller_identity.account.account_id}:${var.bbva_notifications_env}-bbva-ntf"
        }
      }
    }
  ]
}
EOF

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
    Name        = "${module.naming.aws_sqs_queue}-bbva-notifications"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "bbva_notifications_sqs_oldest_message_alarm" {
  alarm_name          = "${module.naming.aws_lambda_function}-sqs-bbva-notifications-oldest-message"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "ApproximateAgeOfOldestMessage"
  namespace           = "AWS/SQS"
  period              = "60"
  statistic           = "Average"
  threshold           = "14400"                                                                      # 4 minutes

  dimensions {
    QueueName = "${module.naming.aws_sqs_queue}-bbva-notifications"
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
    Name        = "${module.naming.aws_lambda_function}-sqs-bbva-notifications-oldest-message"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "bbva_notifications_sqs_dl_receive_alarm" {
  alarm_name          = "${module.naming.aws_lambda_function}-sqs-bbva-notifications-dl-receive"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "NumberOfMessagesReceived"
  namespace           = "AWS/SQS"
  period              = "60"
  statistic           = "Average"
  threshold           = "1"

  dimensions {
    QueueName = "${module.naming.aws_sqs_queue}-bbva-notifications-dead-letter"
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
    Name        = "${module.naming.aws_lambda_function}-sqs-bbva-notifications-dl-receive"
    Team        = "${var.team}"
  }
}

module "bbva_sns" {
  source = "./modules/bbva_sns"

  aws_profile = "${var.bbva_wise_profile}"
  aws_region  = "${var.aws_region}"
  environment = "${var.bbva_sqs_environment}"
}

resource "aws_sns_topic_subscription" "bbva_notifications" {
  topic_arn            = "${module.bbva_sns.bbva_sns_arn}"
  protocol             = "sqs"
  endpoint             = "${aws_sqs_queue.bbva_notifications.arn}"
  raw_message_delivery = true
}
