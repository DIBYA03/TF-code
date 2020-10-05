data "aws_api_gateway_rest_api" "client_api" {
  name = "${module.naming.aws_api_gateway_rest_api}"
}
