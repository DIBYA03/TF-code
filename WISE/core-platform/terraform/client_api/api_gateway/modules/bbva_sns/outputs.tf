output "bbva_sns_arn" {
  value = "${data.aws_sns_topic.bbva_notifications.arn}"
}
