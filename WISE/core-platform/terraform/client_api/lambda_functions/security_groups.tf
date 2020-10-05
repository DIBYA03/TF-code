data "aws_security_group" "lambda_default" {
  name = "${module.naming.aws_security_group}-lambda"
}
