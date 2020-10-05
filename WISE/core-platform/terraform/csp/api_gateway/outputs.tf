output "api_gateway_execution_arn" {
  value = "${aws_api_gateway_rest_api.csp.execution_arn}"
}

output "api_gateway_id" {
  value = "${aws_api_gateway_rest_api.csp.id}"
}

output "aws_iam_role" {
  value = "${aws_iam_role.csp_lambda.arn}"
}

output "lambda_security_group" {
  value = "${aws_security_group.lambda_default.id}"
}
