data "aws_sqs_queue" "segment_analytics" {
  name = "${var.environment}-client-api-segment-analytics"
}
