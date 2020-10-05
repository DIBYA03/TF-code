resource "aws_sqs_queue" "bbva_dead_letter_queue" {
  name       = "${module.naming.aws_sqs_queue}-dead-letter"
  fifo_queue = "${var.bbva_sqs_fifo_queue}"

  delay_seconds              = "${var.bbva_sqs_delay_seconds}"
  receive_wait_time_seconds  = "${var.bbva_sqs_receive_wait_time_seconds}"
  visibility_timeout_seconds = "${var.bbva_sqs_visibility_timeout_seconds}"
  message_retention_seconds  = "${var.bbva_sqs_dl_message_retention_seconds}"

  max_message_size            = "${var.bbva_sqs_max_message_size}"
  content_based_deduplication = "${var.bbva_sqs_dl_content_based_deduplication}"

  kms_master_key_id                 = "${aws_kms_key.bbva_sqs.key_id}"
  kms_data_key_reuse_period_seconds = "${var.bbva_sqs_kms_data_key_reuse_period_seconds}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_sqs_queue}-dead-letter"
    Team        = "${var.team}"
  }
}
