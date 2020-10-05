output "api_gateway_execution_arn" {
  value = "${aws_api_gateway_rest_api.client.execution_arn}"
}

output "api_gateway_id" {
  value = "${aws_api_gateway_rest_api.client.id}"
}

output "cognitoauth_custommessage_lambda" {
  value = "${aws_lambda_function.cognitoauth_custommessage_lambda.arn}"
}

output "cognitoauth_postconfirm_lambda" {
  value = "${aws_lambda_function.cognitoauth_postconfirm_lambda.arn}"
}

output "cognitoauth_presignup_lambda" {
  value = "${aws_lambda_function.cognitoauth_presignup_lambda.arn}"
}

output "cognitoauth_pretoken_lambda" {
  value = "${aws_lambda_function.cognitoauth_pretoken_lambda.arn}"
}

output "aws_iam_role" {
  value = "${aws_iam_role.clientapi_lambda.arn}"
}

output "lambda_security_group" {
  value = "${aws_security_group.lambda_default.id}"
}

output "s3_documents_bucket" {
  value = "${aws_s3_bucket.documents.id}"
}

output "s3_documents_bucket_kms" {
  value = "${aws_kms_key.documents_bucket.key_id}"
}
