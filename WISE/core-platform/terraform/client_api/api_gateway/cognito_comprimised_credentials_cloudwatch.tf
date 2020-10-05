resource "aws_cloudwatch_metric_alarm" "cognito_compromised_credentials_password_change_medium_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-password-medium-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CompromisedCredentialsRisk"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-password-medium-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_compromised_credentials_password_change_high_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-password-high-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CompromisedCredentialsRisk"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-password-high-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_compromised_credentials_signin_change_medium_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-signin-medium-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CompromisedCredentialsRisk"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-signin-medium-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_compromised_credentials_signin_change_high_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-signin-high-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CompromisedCredentialsRisk"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-signin-high-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_compromised_credentials_signup_change_medium_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-signup-medium-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CompromisedCredentialsRisk"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-signup-medium-risk"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "cognito_compromised_credentials_signup_change_high_risk" {
  alarm_name          = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-signup-high-risk"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CompromisedCredentialsRisk"
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
    Name        = "${module.naming.aws_lambda_function}-cognito-compromised-credentials-signup-high-risk"
    Team        = "${var.team}"
  }
}
