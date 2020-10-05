resource "aws_lambda_function" "lambda_warmer_lambda" {
  function_name = "${module.naming.aws_lambda_function}-lambda-warmer"
  role          = "${aws_iam_role.lambda_warmer_lambda.arn}"
  timeout       = "900"                                                # max time to warm lambdas

  filename         = "../../../cmd/lambda/util/lambda_warmer/lambda.zip"
  source_code_hash = "${base64sha256(file("../../../cmd/lambda/util/lambda_warmer/lambda.zip"))}"
  runtime          = "go1.x"
  handler          = "main"

  vpc_config = {
    subnet_ids         = ["${var.app_subnet_ids}"]
    security_group_ids = ["${aws_security_group.lambda_default.id}"]
  }

  environment {
    variables = {
      API_ENV = "${var.environment}"
    }
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_lambda_function}-lambda-warmer"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_event_rule" "lambda_warmer_lambda" {
  name        = "${module.naming.aws_cloudwatch_event_rule}-lambda-warmer"
  description = "trigger lambda warming for ${var.environment_name}"

  schedule_expression = "rate(5 minutes)"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_event_rule}"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_event_target" "lambda_warmer_lambda" {
  rule      = "${aws_cloudwatch_event_rule.lambda_warmer_lambda.name}"
  target_id = "SendToSNS"
  arn       = "${aws_lambda_function.lambda_warmer_lambda.arn}"
}

#Give permissions for CW Event to access the lambda function
resource "aws_lambda_permission" "lambda_warmer_lambda" {
  statement_id  = "AllowExecutionFromCognitoClouwdWatchEvent"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.lambda_warmer_lambda.function_name}"
  principal     = "events.amazonaws.com"
  source_arn    = "${aws_cloudwatch_event_rule.lambda_warmer_lambda.arn}"
}

resource "aws_cloudwatch_metric_alarm" "lambda_warmer_lambda_non_critical_errors" {
  alarm_name          = "${module.naming.aws_lambda_function}-lambda-warmer-non-crit-errors"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "3"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "300"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    FunctionName = "${aws_lambda_function.lambda_warmer_lambda.function_name}"
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
    Name        = "${module.naming.aws_lambda_function}-lambda-warmer-non-crit-errors"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "lambda_warmer_lambda_critical_errors" {
  alarm_name          = "${module.naming.aws_lambda_function}-lambda-warmer-crit-errors"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "3"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "300"
  statistic           = "Average"
  threshold           = "${var.default_lambda_critical_alarm_error_count}"

  dimensions {
    FunctionName = "${aws_lambda_function.lambda_warmer_lambda.function_name}"
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
    Name        = "${module.naming.aws_lambda_function}-lambda-warmer-crit-errors"
    Team        = "${var.team}"
  }
}
