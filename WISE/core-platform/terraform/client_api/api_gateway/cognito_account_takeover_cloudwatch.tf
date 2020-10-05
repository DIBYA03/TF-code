resource "aws_cloudwatch_metric_alarm" "cognito_account_takeover_password_change_medium_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-account-takeover-password-medium-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "AccountTakeoverRisk"
  namespace           = "AWS/Cognito"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    Operation  = "PasswordChange"
    RiskLevel  = "Medium"
    UserPoolId = "${aws_cognito_user_pool.default.id}"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-account-takeover-password-medium-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_account_takeover_password_change_high_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-account-takeover-password-high-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "AccountTakeoverRisk"
  namespace           = "AWS/Cognito"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    Operation  = "PasswordChange"
    RiskLevel  = "High"
    UserPoolId = "${aws_cognito_user_pool.default.id}"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-account-takeover-password-high-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_account_takeover_signin_change_medium_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-account-takeover-signin-medium-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "AccountTakeoverRisk"
  namespace           = "AWS/Cognito"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    Operation  = "SignIn"
    RiskLevel  = "Medium"
    UserPoolId = "${aws_cognito_user_pool.default.id}"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-account-takeover-signin-medium-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_account_takeover_signin_change_high_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-account-takeover-signin-high-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "AccountTakeoverRisk"
  namespace           = "AWS/Cognito"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    Operation  = "SignIn"
    RiskLevel  = "High"
    UserPoolId = "${aws_cognito_user_pool.default.id}"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-account-takeover-signin-high-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_account_takeover_signup_change_medium_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-account-takeover-signup-medium-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "AccountTakeoverRisk"
  namespace           = "AWS/Cognito"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    Operation  = "SignUp"
    RiskLevel  = "Medium"
    UserPoolId = "${aws_cognito_user_pool.default.id}"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-account-takeover-signup-medium-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_account_takeover_signup_change_high_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-account-takeover-signup-high-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "AccountTakeoverRisk"
  namespace           = "AWS/Cognito"
  period              = "60"
  statistic           = "Average"
  threshold           = "${var.default_lambda_non_critical_alarm_error_count}"

  dimensions {
    Operation  = "SignUp"
    RiskLevel  = "High"
    UserPoolId = "${aws_cognito_user_pool.default.id}"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-account-takeover-signup-high-risk"
    Team        = "${var.team}"
  }
}
