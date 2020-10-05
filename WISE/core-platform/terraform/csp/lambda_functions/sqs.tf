data "aws_sqs_queue" "segment_analytics" {
  name = "${var.environment}-client-api-segment-analytics"
}

data "aws_sqs_queue" "business_document_upload" {
  name = "${module.naming.aws_sqs_queue}-document-upload"
}

data "aws_sqs_queue" "review" {
  name = "${module.naming.aws_sqs_queue}-review"
}
