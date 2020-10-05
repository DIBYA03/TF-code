resource "aws_sqs_queue" "shopify_order_dead_letter_queue" {
  name       = "${module.naming.aws_sqs_queue}-shopify-order-dead-letter"
  fifo_queue = "${var.internal_sqs_fifo_queue}"

  delay_seconds              = "${var.internal_sqs_delay_seconds}"
  receive_wait_time_seconds  = "${var.internal_sqs_receive_wait_time_seconds}"
  visibility_timeout_seconds = "${var.internal_sqs_visibility_timeout_seconds}"
  message_retention_seconds  = "${var.internal_sqs_dl_message_retention_seconds}"

  max_message_size            = "${var.internal_sqs_max_message_size}"
  content_based_deduplication = "${var.internal_sqs_dl_content_based_deduplication}"

  kms_master_key_id                 = "${aws_kms_key.internal_sqs.key_id}"
  kms_data_key_reuse_period_seconds = "${var.internal_sqs_kms_data_key_reuse_period_seconds}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_sqs_queue}-shopify-order-dead-letter"
    Team        = "${var.team}"
  }
}

resource "aws_sqs_queue" "shopify_order" {
  name       = "${module.naming.aws_sqs_queue}-shopify-order"
  fifo_queue = "${var.shopify_order_sqs_fifo_queue}"

  delay_seconds              = "${var.shopify_order_sqs_delay_seconds}"
  receive_wait_time_seconds  = "${var.shopify_order_sqs_receive_wait_time_seconds}"
  visibility_timeout_seconds = "${var.shopify_order_sqs_visibility_timeout_seconds}"
  message_retention_seconds  = "${var.shopify_order_sqs_message_retention_seconds}"

  max_message_size            = "${var.shopify_order_sqs_max_message_size}"
  content_based_deduplication = "${var.shopify_order_sqs_content_based_deduplication}"

  kms_master_key_id                 = "${aws_kms_alias.internal_sqs.target_key_arn}"
  kms_data_key_reuse_period_seconds = "${var.shopify_order_sqs_kms_data_key_reuse_period_seconds}"

  policy = ""

  redrive_policy = <<EOF
{
  "deadLetterTargetArn": "${aws_sqs_queue.shopify_order_dead_letter_queue.arn}",
  "maxReceiveCount": 2
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_sqs_queue}-shopify-order"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "shopify_order_sqs_oldest_message_alarm" {
  alarm_name          = "${module.naming.aws_lambda_function}-sqs-shopify-order-oldest-message"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "ApproximateAgeOfOldestMessage"
  namespace           = "AWS/SQS"
  period              = "60"
  statistic           = "Average"
  threshold           = "14400"                                                                     # 4 minutes

  dimensions {
    QueueName = "${module.naming.aws_sqs_queue}-shopify-order"
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
    Name        = "${module.naming.aws_lambda_function}-sqs-shopify-order-oldest-message"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "shopify_order_sqs_dl_receive_alarm" {
  alarm_name          = "${module.naming.aws_lambda_function}-sqs-shopify-order-dl-receive"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "NumberOfMessagesReceived"
  namespace           = "AWS/SQS"
  period              = "60"
  statistic           = "Average"
  threshold           = "1"

  dimensions {
    QueueName = "${module.naming.aws_sqs_queue}-shopify-order-dead-letter"
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
    Name        = "${module.naming.aws_lambda_function}-sqs-shopify-order-dl-receive"
    Team        = "${var.team}"
  }
}
