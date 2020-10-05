resource "aws_lambda_function" "cognitoauth_pretoken_lambda" {
  function_name = "${module.naming.aws_lambda_function}-cognitoauth-pretoken"
  role          = "${aws_iam_role.cognitoauth_lambda.arn}"

  filename         = "../../../cmd/lambda/csp/cognitoauth/pretoken/lambda.zip"
  source_code_hash = "${base64sha256(file("../../../cmd/lambda/csp/cognitoauth/pretoken/lambda.zip"))}"
  handler          = "main"
  runtime          = "go1.x"
  timeout          = "${var.lambda_timeout}"

  kms_key_arn = "${aws_kms_alias.cognito_lambda.target_key_arn}"

  vpc_config = {
    subnet_ids         = ["${var.app_subnet_ids}"]
    security_group_ids = ["${aws_security_group.lambda_default.id}"]
  }

  environment {
    variables = {
      API_ENV = "${var.environment_name}"

      CSP_DB_WRITE_URL  = "${data.aws_ssm_parameter.csp_rds_master_endpoint.value}"
      CSP_DB_READ_URL   = "${data.aws_ssm_parameter.csp_rds_read_endpoint.value}"
      CSP_DB_WRITE_PORT = "${data.aws_ssm_parameter.csp_rds_port.value}"
      CSP_DB_READ_PORT  = "${data.aws_ssm_parameter.csp_rds_port.value}"
      CSP_DB_NAME       = "${data.aws_ssm_parameter.csp_rds_db_name.value}"
      CSP_DB_USER       = "${data.aws_ssm_parameter.csp_rds_username.value}"
      CSP_DB_PASSWD     = "${data.aws_ssm_parameter.csp_rds_password.value}"

      CSP_NOTIFICATION_SLACK_CHANNEL = "${data.aws_ssm_parameter.csp_notification_slack_channel.value}"
      CSP_NOTIFICATION_SLACK_URL     = "${data.aws_ssm_parameter.csp_notification_slack_url.value}"
    }
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_lambda_function}-cognitoauth-pretoken"
    Team        = "${var.team}"
  }
}

# Give permissions for Cognito to access the lambda function
resource "aws_lambda_permission" "cognitoauth_pretoken_lambda" {
  statement_id = "AllowExecutionFromCognito"

  function_name = "${module.naming.aws_lambda_function}-cognitoauth-pretoken"
  source_arn    = "${aws_cognito_user_pool.default.arn}"
  action        = "lambda:InvokeFunction"
  principal     = "cognito-idp.amazonaws.com"
}

resource "aws_cloudwatch_metric_alarm" "cognitoauth_pretoken_lambda_non_critical_errors" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognitoauth-pretoken-non-crit-errors"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    FunctionName = "${aws_lambda_function.cognitoauth_pretoken_lambda.function_name}"
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
    Name        = "${module.naming.aws_lambda_function}-cognitoauth-pretoken-non-crit-errors"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognitoauth_pretoken_lambda_critical_errors" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognitoauth-pretoken-crit-errors"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_critical_alarm_error_count}"

  dimensions {
    FunctionName = "${aws_lambda_function.cognitoauth_pretoken_lambda.function_name}"
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
    Name        = "${module.naming.aws_lambda_function}-cognitoauth-pretoken-crit-errors"
    Team        = "${var.team}"
  }
}
