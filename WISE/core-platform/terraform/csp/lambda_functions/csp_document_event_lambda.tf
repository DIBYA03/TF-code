locals {
  csp_document_event_lambda_env_vars = "${merge(
    local.api_env,
    local.bbva_app_credentials,
    local.core_db_credentials,
    local.csp_db_credentials,
    local.bank_db_credentials,
    local.identity_db_credentials,
    local.business_document_sqs,
    local.segment_sqs,
  )}"
}

resource "aws_lambda_function" "csp_document_event_lambda" {
  function_name = "${module.naming.aws_lambda_function}-doc-evt-${var.api_gw_stage}"
  role          = "${aws_iam_role.csp_document_event_lambda.arn}"
  kms_key_arn   = "${aws_kms_key.lambda_default.arn}"
  timeout       = "${var.lambda_timeout}"

  filename         = "../../../cmd/lambda/csp/documentevent/lambda.zip"
  source_code_hash = "${base64sha256(file("../../../cmd/lambda/csp/documentevent/lambda.zip"))}"

  handler = "main"
  runtime = "go1.x"

  vpc_config = {
    subnet_ids         = ["${var.app_subnet_ids}"]
    security_group_ids = ["${data.aws_security_group.lambda_default.id}"]
  }

  environment {
    variables = "${local.csp_document_event_lambda_env_vars}"
  }

  tags {
    Application = "${var.application}"
    Environment = "${var.environment_name}"
    Component   = "${var.component}"
    Name        = "${module.naming.aws_lambda_function}-doc-evt-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

# Give permission to S3 to trigger this lambda
resource "aws_lambda_permission" "csp_document_event_lambda" {
  statement_id = "AllowExecutionFromS3Bucket"
  action       = "lambda:InvokeFunction"
  principal    = "s3.amazonaws.com"

  function_name = "${aws_lambda_function.csp_document_event_lambda.function_name}"
  source_arn    = "${data.aws_s3_bucket.documents.arn}"
}

resource "aws_s3_bucket_notification" "document_event_lambda" {
  bucket = "${data.aws_s3_bucket.documents.id}"

  lambda_function {
    id                  = "${module.naming.aws_lambda_function}-doc-evt-${var.api_gw_stage}"
    lambda_function_arn = "${aws_lambda_function.csp_document_event_lambda.arn}"
    events              = ["s3:ObjectCreated:*"]
  }
}

resource "aws_cloudwatch_metric_alarm" "csp_document_event_lambda_non_critical_errors" {
  alarm_name          = "${module.naming.aws_lambda_function}-doc-evt-non-crit-errors"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    FunctionName = "${aws_lambda_function.csp_document_event_lambda.function_name}"
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
    Name        = "${module.naming.aws_lambda_function}-doc-evt-non-crit-errors"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "csp_document_event_lambda_critical_errors" {
  alarm_name          = "${module.naming.aws_lambda_function}-doc-evt-crit-errors"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_critical_alarm_error_count}"

  dimensions {
    FunctionName = "${aws_lambda_function.csp_document_event_lambda.function_name}"
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
    Name        = "${module.naming.aws_lambda_function}-doc-evt-crit-errors"
    Team        = "${var.team}"
  }
}
