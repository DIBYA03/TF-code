locals {
  api_gw_name = "${module.naming.aws_api_gateway_rest_api}"
}

# The API GW deployed from api_gateway tf
data "aws_api_gateway_rest_api" "client" {
  name = "${local.api_gw_name}"
}

# ARN of API gateway
locals {
  api_gw_arn = "arn:aws:execute-api:${var.aws_region}:${data.aws_caller_identity.account.account_id}:${data.aws_api_gateway_rest_api.client.id}"
}

# Deploy the stage
resource "aws_api_gateway_deployment" "client" {
  rest_api_id = "${data.aws_api_gateway_rest_api.client.id}"
  stage_name  = "${var.api_gw_stage}"
}
