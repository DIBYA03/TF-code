data "aws_sqs_queue" "stripe_webhook" {
  name = "${module.naming.aws_sqs_queue}-stripe-req-payments"
}

data "aws_sqs_queue" "segment_analytics" {
  name = "${module.naming.aws_sqs_queue}-segment-analytics"
}

data "aws_sqs_queue" "review" {
  name = "${var.csp_environment}-csp-api-review"
}

data "aws_sqs_queue" "internal_banking" {
  name = "${module.naming.aws_sqs_queue}-banking-notifications"
}
