data "aws_s3_bucket" "documents" {
  bucket = "wiseus-${var.aws_region}-${var.environment}-client-api-documents"
}
