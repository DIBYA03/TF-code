output "bbva_sqs_dead_letter_arn" {
  value = "${aws_sqs_queue.bbva_dead_letter_queue.arn}"
}

output "bbva_notifications_sqs_queue" {
  value = "${aws_sqs_queue.bbva_notifications.arn}"
}
