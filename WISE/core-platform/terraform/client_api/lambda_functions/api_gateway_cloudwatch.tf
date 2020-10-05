resource "aws_cloudwatch_metric_alarm" "api_gateway_5XX_non_crtical" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-5XX-errors-non-crit"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "5"
  metric_name         = "5XXError"
  namespace           = "AWS/ApiGateway"
  period              = "60"
  statistic           = "Sum"
  threshold           = "${var.api_gw_5XX_error_alarm_non_critical_threshold}"

  dimensions {
    ApiName = "${local.api_gw_name}"
    Stage   = "${var.api_gw_stage}"
  }

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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-5XX-errors-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "api_gateway_5XX_crtical" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-5XX-errors-crit"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "5"
  metric_name         = "5XXError"
  namespace           = "AWS/ApiGateway"
  period              = "60"
  statistic           = "Sum"
  threshold           = "${var.api_gw_5XX_error_alarm_critical_threshold}"

  dimensions {
    ApiName = "${local.api_gw_name}"
    Stage   = "${var.api_gw_stage}"
  }

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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-5XX-errors-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "api_gateway_4XX_non_crtical" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-4XX-errors-non-crit"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "5"
  metric_name         = "4XXError"
  namespace           = "AWS/ApiGateway"
  period              = "60"
  statistic           = "Sum"
  threshold           = "${var.api_gw_4XX_error_alarm_non_critical_threshold}"

  dimensions {
    ApiName = "${local.api_gw_name}"
    Stage   = "${var.api_gw_stage}"
  }

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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-4XX-errors-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "api_gateway_4XX_crtical" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-4XX-errors-crit"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "5"
  metric_name         = "4XXError"
  namespace           = "AWS/ApiGateway"
  period              = "60"
  statistic           = "Sum"
  threshold           = "${var.api_gw_4XX_error_alarm_critical_threshold}"

  dimensions {
    ApiName = "${local.api_gw_name}"
    Stage   = "${var.api_gw_stage}"
  }

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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-4XXX-errors-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "api_gateway_latency" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-high-latency"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "5"
  metric_name         = "Latency"
  namespace           = "AWS/ApiGateway"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.api_gw_latency_alarm_threshold}"

  dimensions {
    ApiName = "${local.api_gw_name}"
    Stage   = "${var.api_gw_stage}"
  }

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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-api-gw-high-latency"
    Team        = "${var.team}"
  }
}
