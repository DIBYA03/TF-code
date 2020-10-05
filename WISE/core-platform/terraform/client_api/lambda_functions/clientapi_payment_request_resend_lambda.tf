resource "aws_lambda_function" "clientapi_payment_request_resend" {
  function_name = "${module.naming.aws_lambda_function}-pmt-rqt-rsd-${var.api_gw_stage}"
  role          = "${aws_iam_role.clientapi_payment_request_resend_lambda.arn}"
  kms_key_arn   = "${aws_kms_key.api_lambda.arn}"
  timeout       = "${var.lambda_timeout}"

  filename         = "../../../cmd/lambda/clientapi/payment/request/resend/lambda.zip"
  source_code_hash = "${base64sha256(file("../../../cmd/lambda/clientapi/payment/request/resend/lambda.zip"))}"

  runtime = "go1.x"
  handler = "main"

  vpc_config = {
    subnet_ids         = ["${var.app_subnet_ids}"]
    security_group_ids = ["${data.aws_security_group.lambda_default.id}"]
  }

  environment {
    variables = {
      API_ENV                = "${var.environment_name}"
      AWS_S3_BUCKET_DOCUMENT = "${data.aws_s3_bucket.documents.id}"

      BBVA_APP_ENV    = "${data.aws_ssm_parameter.bbva_app_env.value}"
      BBVA_APP_ID     = "${data.aws_ssm_parameter.bbva_app_id.value}"
      BBVA_APP_NAME   = "${data.aws_ssm_parameter.bbva_app_name.value}"
      BBVA_APP_SECRET = "${data.aws_ssm_parameter.bbva_app_secret.value}"

      IDENTITY_DB_WRITE_URL  = "${data.aws_ssm_parameter.rds_master_endpoint.value}"
      IDENTITY_DB_READ_URL   = "${data.aws_ssm_parameter.rds_read_endpoint.value}"
      IDENTITY_DB_WRITE_PORT = "${data.aws_ssm_parameter.rds_port.value}"
      IDENTITY_DB_READ_PORT  = "${data.aws_ssm_parameter.rds_port.value}"
      IDENTITY_DB_NAME       = "${data.aws_ssm_parameter.identity_rds_db_name.value}"
      IDENTITY_DB_USER       = "${data.aws_ssm_parameter.identity_rds_user_name.value}"
      IDENTITY_DB_PASSWD     = "${data.aws_ssm_parameter.identity_rds_password.value}"

      CORE_DB_WRITE_URL  = "${data.aws_ssm_parameter.rds_master_endpoint.value}"
      CORE_DB_READ_URL   = "${data.aws_ssm_parameter.rds_read_endpoint.value}"
      CORE_DB_WRITE_PORT = "${data.aws_ssm_parameter.rds_port.value}"
      CORE_DB_READ_PORT  = "${data.aws_ssm_parameter.rds_port.value}"
      CORE_DB_NAME       = "${data.aws_ssm_parameter.core_rds_db_name.value}"
      CORE_DB_USER       = "${data.aws_ssm_parameter.core_rds_user_name.value}"
      CORE_DB_PASSWD     = "${data.aws_ssm_parameter.core_rds_password.value}"

      BANK_DB_WRITE_URL  = "${data.aws_ssm_parameter.rds_master_endpoint.value}"
      BANK_DB_READ_URL   = "${data.aws_ssm_parameter.rds_read_endpoint.value}"
      BANK_DB_WRITE_PORT = "${data.aws_ssm_parameter.rds_port.value}"
      BANK_DB_READ_PORT  = "${data.aws_ssm_parameter.rds_port.value}"
      BANK_DB_NAME       = "${data.aws_ssm_parameter.bank_rds_db_name.value}"
      BANK_DB_USER       = "${data.aws_ssm_parameter.bank_rds_user_name.value}"
      BANK_DB_PASSWD     = "${data.aws_ssm_parameter.bank_rds_password.value}"

      REDIS_ENDPOINT = "${data.aws_ssm_parameter.redis_endpoint.value}"
      REDIS_PORT     = "${data.aws_ssm_parameter.redis_port.value}"
      REDIS_PASSWORD = "${data.aws_ssm_parameter.redis_password.value}"

      PAYMENTS_URL = "https://${data.aws_ssm_parameter.payments_url.value}"

      PLAID_ENV        = "${data.aws_ssm_parameter.plaid_env.value}"
      PLAID_CLIENT_ID  = "${data.aws_ssm_parameter.plaid_client_id.value}"
      PLAID_SECRET     = "${data.aws_ssm_parameter.plaid_secret.value}"
      PLAID_PUBLIC_KEY = "${data.aws_ssm_parameter.plaid_public_key.value}"

      SENDGRID_API_KEY = "${data.aws_ssm_parameter.sendgrid_api_key.value}"

      SQS_REGION      = "${var.aws_region}"
      SEGMENT_SQS_URL = "${data.aws_sqs_queue.segment_analytics.id}"

      STRIPE_KEY            = "${data.aws_ssm_parameter.stripe_key.value}"
      STRIPE_PUBLISH_KEY    = "${data.aws_ssm_parameter.stripe_publish_key.value}"
      STRIPE_WEBHOOK_SECRET = "${data.aws_ssm_parameter.stripe_webhook_secret.value}"

      WISE_CLEARING_ACCOUNT_ID                = "${data.aws_ssm_parameter.wise_clearing_account_id.value}"
      WISE_CLEARING_BUSINESS_ID               = "${data.aws_ssm_parameter.wise_clearing_business_id.value}"
      WISE_CLEARING_USER_ID                   = "${data.aws_ssm_parameter.wise_clearing_user_id.value}"
      WISE_CLEARING_MAX_MONEY_REQUEST_ALLOWED = "${var.wise_clearing_max_request_amount}"
      CARD_READER_MAX_MONEY_REQUEST_ALLOWED   = "${var.card_reader_max_request_amount}"
      CARD_ONLINE_MAX_MONEY_REQUEST_ALLOWED   = "${var.card_online_max_request_amount}"

      WISE_INVOICE_EMAIL = "${data.aws_ssm_parameter.wise_invoice_email_address.value}"
      WISE_SUPPORT_EMAIL = "${data.aws_ssm_parameter.wise_support_email_address.value}"
      WISE_SUPPORT_NAME  = "${data.aws_ssm_parameter.wise_support_email_name.value}"
      WISE_SUPPORT_PHONE = "${data.aws_ssm_parameter.wise_support_phone.value}"

      GRPC_SERVICE_PORT       = "${var.grpc_port}"
      USE_TRANSACTION_SERVICE = "${var.use_transaction_service}"
      USE_BANKING_SERVICE     = "${var.use_banking_service}"
      USE_INVOICE_SERVICE     = "${var.use_invoice_service}"

      TWILIO_ACCOUNT_SID  = "${data.aws_ssm_parameter.twilio_account_sid.value}"
      TWILIO_API_SID      = "${data.aws_ssm_parameter.twilio_api_sid.value}"
      TWILIO_API_SECRET   = "${data.aws_ssm_parameter.twilio_api_secret.value}"
      TWILIO_SENDER_PHONE = "${data.aws_ssm_parameter.twilio_sender_phone.value}"
    }
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_lambda_function}-pmt-rqt-rsd-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

# Give permissions for the API Gateway to access the lambda function
resource "aws_lambda_permission" "clientapi_payment_request_resend" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${module.naming.aws_lambda_function}-pmt-rqt-rsd-${var.api_gw_stage}"

  principal  = "apigateway.amazonaws.com"
  source_arn = "${local.api_gw_arn}/*/*/*"

  depends_on = [
    "aws_lambda_function.clientapi_payment_request_resend",
  ]
}

resource "aws_cloudwatch_metric_alarm" "clientapi_payment_request_resend_non_critical_errors" {
  alarm_name          = "${module.naming.aws_lambda_function}-pmt-rqt-rsd-non-crit-errors-${var.api_gw_stage}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    FunctionName = "${aws_lambda_function.clientapi_payment_request_resend.function_name}"
  }

  treat_missing_data = "notBreaching"

  alarm_actions = [
    "${var.sns_non_critical_topic}",
  ]

  ok_actions = [
    "${var.sns_non_critical_topic}",
  ]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_lambda_function}-pmt-rqt-rsd-non-crit-errors-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "clientapi_payment_request_resend_critical_errors" {
  alarm_name          = "${module.naming.aws_lambda_function}-pmt-rqt-rsd-crit-errors-${var.api_gw_stage}"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_critical_alarm_error_count}"

  dimensions {
    FunctionName = "${aws_lambda_function.clientapi_payment_request_resend.function_name}"
  }

  treat_missing_data = "notBreaching"

  alarm_actions = [
    "${var.sns_critical_topic}",
  ]

  ok_actions = [
    "${var.sns_critical_topic}",
  ]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_lambda_function}-pmt-rqt-rsd-crit-errors-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}
