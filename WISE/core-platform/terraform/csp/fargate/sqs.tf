data "aws_sqs_queue" "business_document_upload" {
  name = "${var.environment}-csp-api-document-upload"
}

data "aws_sqs_queue" "review" {
  name = "${var.environment}-csp-api-review"
}

data "aws_sqs_queue" "internal_banking" {
  name = "${var.environment}-client-api-banking-notifications"
}

data "aws_sqs_queue" "segment_analytics" {
  name = "${var.environment}-client-api-segment-analytics"
}
