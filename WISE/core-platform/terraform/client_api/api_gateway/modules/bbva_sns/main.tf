data "aws_sns_topic" "bbva_notifications" {
  name = "${var.environment}-bbva-ntf"

  provider = "aws.${var.provider_name}"
}
