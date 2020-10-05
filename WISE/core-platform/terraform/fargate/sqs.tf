data "aws_sqs_queue" "stripe_webhook" {
  name = "${var.environment}-client-api-stripe-req-payments"
}

data "aws_sqs_queue" "signature_webhook" {
  name = "${var.environment}-client-api-signature"
}

data "aws_sqs_queue" "internal_banking" {
  name = "${var.environment}-client-api-banking-notifications"
}

data "aws_sqs_queue" "bbva_notifications" {
  name = "${var.environment}-client-api-bbva-notifications"
}

data "aws_sqs_queue" "segment_analytics" {
  name = "${var.environment}-client-api-segment-analytics"
}

data "aws_sqs_queue" "shopify_order" {
  name = "${var.environment}-client-api-shopify-order"
}

data "aws_sqs_queue" "csp_review" {
  name = "${var.environment}-csp-api-review"
}
