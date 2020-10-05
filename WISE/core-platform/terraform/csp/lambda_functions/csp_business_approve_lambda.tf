locals {
  csp_business_approve_lambda_env_vars = "${merge(
    local.api_env,
    local.bbva_app_credentials,
    local.core_db_credentials,
    local.csp_db_credentials,
    local.bank_db_credentials,
    local.identity_db_credentials,
    local.business_document_sqs,
    local.segment_sqs,
    local.review_sqs,
  )}"
}

resource "aws_lambda_function" "csp_business_approve_lambda" {
  function_name = "${module.naming.aws_lambda_function}-bus-apr-${var.api_gw_stage}"
  role          = "${aws_iam_role.csp_business_approve_lambda.arn}"
  kms_key_arn   = "${aws_kms_key.lambda_default.arn}"
  timeout       = "${var.lambda_timeout}"

  filename         = "../../../cmd/lambda/csp/review/business/approve/lambda.zip"
  source_code_hash = "${base64sha256(file("../../../cmd/lambda/csp/review/business/approve/lambda.zip"))}"

  handler = "main"
  runtime = "go1.x"

  vpc_config = {
    subnet_ids         = ["${var.app_subnet_ids}"]
    security_group_ids = ["${data.aws_security_group.lambda_default.id}"]
  }

  environment {
    variables = "${local.csp_business_approve_lambda_env_vars}"
  }

  tags {
    Application = "${var.application}"
    Environment = "${var.environment_name}"
    Component   = "${var.component}"
    Name        = "${module.naming.aws_lambda_function}-bus-apr-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

# Give permissions for the API Gateway to access the lambda function
resource "aws_lambda_permission" "csp_business_approve_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.csp_business_approve_lambda.function_name}"

  principal  = "apigateway.amazonaws.com"
  source_arn = "${local.api_gw_arn}/*/*/*"

  depends_on = [
    "aws_lambda_function.csp_business_approve_lambda",
  ]
}
