resource "aws_sqs_queue" "review_dead_letter_queue" {
  name       = "${module.naming.aws_sqs_queue}-review-dead-letter"
  fifo_queue = "${var.review_sqs_fifo_queue}"

  delay_seconds              = "${var.review_sqs_delay_seconds}"
  receive_wait_time_seconds  = "${var.review_sqs_receive_wait_time_seconds}"
  visibility_timeout_seconds = "${var.review_sqs_visibility_timeout_seconds}"
  message_retention_seconds  = "${var.review_sqs_dl_message_retention_seconds}"

  max_message_size            = "${var.review_sqs_max_message_size}"
  content_based_deduplication = "${var.review_sqs_dl_content_based_deduplication}"

  kms_master_key_id                 = "${data.aws_kms_alias.env_default.target_key_arn}"
  kms_data_key_reuse_period_seconds = "${var.review_sqs_kms_data_key_reuse_period_seconds}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_sqs_queue}-review-dead-letter"
    Team        = "${var.team}"
  }
}

resource "aws_sqs_queue" "review" {
  name       = "${module.naming.aws_sqs_queue}-review"
  fifo_queue = "${var.review_sqs_fifo_queue}"

  delay_seconds              = "${var.review_sqs_delay_seconds}"
  receive_wait_time_seconds  = "${var.review_sqs_receive_wait_time_seconds}"
  visibility_timeout_seconds = "${var.review_sqs_visibility_timeout_seconds}"
  message_retention_seconds  = "${var.review_sqs_message_retention_seconds}"

  max_message_size            = "${var.review_sqs_max_message_size}"
  content_based_deduplication = "${var.review_sqs_content_based_deduplication}"

  kms_master_key_id                 = "${data.aws_kms_alias.env_default.target_key_arn}"
  kms_data_key_reuse_period_seconds = "${var.review_sqs_kms_data_key_reuse_period_seconds}"

  policy = ""

  redrive_policy = <<EOF
{
  "deadLetterTargetArn": "${aws_sqs_queue.review_dead_letter_queue.arn}",
  "maxReceiveCount": 2
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_sqs_queue}-review"
    Team        = "${var.team}"
  }
}
