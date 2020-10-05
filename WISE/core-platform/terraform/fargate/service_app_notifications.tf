resource "aws_cloudwatch_log_group" "app_notifications" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-app-notifications"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-app-notifications"
    Team        = "${var.team}"
  }
}

data "template_file" "app_notifications" {
  template = "${file("definitions/tasks/app_notification.json")}"

  vars {
    name    = "${var.app_notification_name}"
    api_env = "${var.environment}"

    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-app-notifications"
    aws_region    = "${var.aws_region}"

    aws_s3_bucket_document = "${data.aws_s3_bucket.documents.id}"

    bbva_app_env    = "${data.aws_ssm_parameter.bbva_app_env.arn}"
    bbva_app_id     = "${data.aws_ssm_parameter.bbva_app_id.arn}"
    bbva_app_name   = "${data.aws_ssm_parameter.bbva_app_name.arn}"
    bbva_app_secret = "${data.aws_ssm_parameter.bbva_app_secret.arn}"

    core_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    core_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    core_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    core_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    core_db_name       = "${data.aws_ssm_parameter.core_rds_db_name.arn}"
    core_db_user       = "${data.aws_ssm_parameter.core_rds_user_name.arn}"
    core_db_passwd     = "${data.aws_ssm_parameter.core_rds_password.arn}"

    bank_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    bank_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    bank_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    bank_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    bank_db_name       = "${data.aws_ssm_parameter.bank_rds_db_name.arn}"
    bank_db_user       = "${data.aws_ssm_parameter.bank_rds_user_name.arn}"
    bank_db_passwd     = "${data.aws_ssm_parameter.bank_rds_password.arn}"

    identity_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    identity_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    identity_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    identity_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    identity_db_name       = "${data.aws_ssm_parameter.identity_rds_db_name.arn}"
    identity_db_user       = "${data.aws_ssm_parameter.identity_rds_user_name.arn}"
    identity_db_passwd     = "${data.aws_ssm_parameter.identity_rds_password.arn}"

    txn_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    txn_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    txn_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    txn_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    txn_db_name       = "${data.aws_ssm_parameter.txn_rds_db_name.arn}"
    txn_db_user       = "${data.aws_ssm_parameter.txn_rds_user_name.arn}"
    txn_db_passwd     = "${data.aws_ssm_parameter.txn_rds_password.arn}"

    ecr_image      = "${var.app_notification_image}"
    ecr_image_tag  = "${var.app_notification_image_tag}"
    fargate_cpu    = "${var.app_notification_cpu}"
    fargate_memory = "${var.app_notification_mem}"

    grpc_port               = "${var.grpc_port}"
    use_transaction_service = "${var.use_transaction_service}"
    use_banking_service     = "${var.use_banking_service}"
    use_invoice_service     = "${var.use_invoice_service}"

    firebase_config  = "${data.aws_ssm_parameter.firebase_config.arn}"
    sendgrid_api_key = "${data.aws_ssm_parameter.sendgrid_api_key.arn}"

    kinesis_trx_name   = "${var.txn_kinesis_name}"
    kinesis_trx_region = "${var.txn_kinesis_region}"

    csp_review_sqs_region = "${var.aws_region}"
    csp_review_sqs_url    = "${local.csp_review_sqs_url}"

    sqs_region     = "${var.aws_region}"
    sqs_url        = "${data.aws_sqs_queue.internal_banking.url}"
    sqs_app_region = "${var.aws_region}"

    segment_sqs_url = "${data.aws_sqs_queue.segment_analytics.url}"

    wise_clearing_account_id  = "${data.aws_ssm_parameter.wise_clearing_account_id.arn}"
    wise_clearing_business_id = "${data.aws_ssm_parameter.wise_clearing_business_id.arn}"
    wise_clearing_user_id     = "${data.aws_ssm_parameter.wise_clearing_user_id.arn}"
    wise_support_email        = "${data.aws_ssm_parameter.wise_support_email_address.arn}"
    wise_support_name         = "${data.aws_ssm_parameter.wise_support_email_name.arn}"
    wise_invoice_email        = "${data.aws_ssm_parameter.wise_invoice_email_address.arn}"

    twilio_account_sid  = "${data.aws_ssm_parameter.twilio_account_sid.arn}"
    twilio_api_sid      = "${data.aws_ssm_parameter.twilio_api_sid.arn}"
    twilio_api_secret   = "${data.aws_ssm_parameter.twilio_api_secret.arn}"
    twilio_sender_phone = "${data.aws_ssm_parameter.twilio_sender_phone.arn}"
  }
}

resource "aws_ecs_task_definition" "app_notifications" {
  family             = "${module.naming.aws_ecs_task_definition}-app-notifications"
  execution_role_arn = "${aws_iam_role.app_notifications_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.app_notifications_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.app_notification_cpu}"
  memory                   = "${var.app_notification_mem}"
  container_definitions    = "${data.template_file.app_notifications.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-app-notifications"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "app_notifications" {
  name            = "${module.naming.aws_ecs_service}-app-notifications"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.app_notifications.arn}"
  desired_count   = "${var.app_notification_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.app_notifications_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }
}

resource "aws_appautoscaling_target" "app_notifications" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.app_notifications.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.app_notification_min_container_count}"
  max_capacity       = "${var.app_notification_max_container_count}"
}

resource "aws_appautoscaling_policy" "app_notifications_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-app-notifications-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.app_notifications.name}"
  scalable_dimension = "ecs:service:DesiredCount"

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Maximum"

    step_adjustment {
      metric_interval_lower_bound = 0
      scaling_adjustment          = 1
    }
  }

  depends_on = ["aws_appautoscaling_target.app_notifications"]
}

resource "aws_appautoscaling_policy" "app_notifications_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-app-notifications-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.app_notifications.name}"
  scalable_dimension = "ecs:service:DesiredCount"

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Maximum"

    step_adjustment {
      metric_interval_lower_bound = 0
      scaling_adjustment          = -1
    }
  }

  depends_on = ["aws_appautoscaling_target.app_notifications"]
}

resource "aws_cloudwatch_metric_alarm" "app_notifications_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-app-notifications-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.app_notifications.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.app_notifications_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "app_notifications_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-app-notifications-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.app_notifications.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.app_notifications_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "app_notifications_low_container_count" {
  count               = "${var.app_notification_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-app-notifications-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.app_notification_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.app_notifications.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "app_notifications_high_container_count" {
  count               = "${var.app_notification_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-app-notifications-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.app_notification_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.app_notifications.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
