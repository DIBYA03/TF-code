# The API GW deployed from api_gateway tf
data "aws_api_gateway_rest_api" "csp" {
  name = "${module.naming.aws_api_gateway_rest_api}"
}

# ARN of API gateway
locals {
  api_gw_arn = "arn:aws:execute-api:${var.aws_region}:${data.aws_caller_identity.account.account_id}:${data.aws_api_gateway_rest_api.csp.id}"
}

# Deploy the stage
resource "aws_api_gateway_deployment" "csp" {
  rest_api_id = "${data.aws_api_gateway_rest_api.csp.id}"
  stage_name  = "${var.api_gw_stage}"
}
