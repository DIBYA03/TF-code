output "lambda_arn" {
  value = "${aws_lambda_function.hsts.arn}"
}

output "lambda_qualified_arn" {
  value = "${aws_lambda_function.hsts.qualified_arn}"
}
