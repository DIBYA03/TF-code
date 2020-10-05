data "aws_s3_bucket" "documents" {
  bucket = "${module.naming.aws_s3_bucket}-documents"
}
