resource "aws_cloudwatch_log_group" "bbva_notifications" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-bbva-notifications"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-bbva-notifications"
    Team        = "${var.team}"
  }
}

data "template_file" "bbva_notifications" {
  template = "${file("definitions/tasks/bbva_notification.json")}"

  vars {
    name          = "${var.bbva_notification_name}"
    api_env       = "${var.environment}"
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-bbva-notifications"
    aws_region    = "${var.aws_region}"

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

    ecr_image      = "${var.bbva_notification_image}"
    ecr_image_tag  = "${var.bbva_notification_image_tag}"
    fargate_cpu    = "${var.bbva_notification_cpu}"
    fargate_memory = "${var.bbva_notification_mem}"

    kinesis_bank_notif_name   = "${var.ntf_kinesis_name}"
    kinesis_bank_notif_region = "${var.ntf_kinesis_region}"

    sqs_banking_region = "${var.aws_region}"
    sqs_banking_url    = "${data.aws_sqs_queue.internal_banking.url}"
    sqs_bbva_region    = "${var.aws_region}"
    sqs_bbva_url       = "${data.aws_sqs_queue.bbva_notifications.url}"

    grpc_port               = "${var.grpc_port}"
    use_transaction_service = "${var.use_transaction_service}"
    use_banking_service     = "${var.use_banking_service}"
  }
}

resource "aws_ecs_task_definition" "bbva_notifications" {
  family             = "${module.naming.aws_ecs_task_definition}-bbva-notifications"
  execution_role_arn = "${aws_iam_role.bbva_notifications_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.bbva_notifications_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.bbva_notification_cpu}"
  memory                   = "${var.bbva_notification_mem}"
  container_definitions    = "${data.template_file.bbva_notifications.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-bbva-notifications"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "bbva_notifications" {
  name            = "${module.naming.aws_ecs_service}-bbva-notifications"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.bbva_notifications.arn}"
  desired_count   = "${var.bbva_notification_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.bbva_notifications_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }
}

resource "aws_appautoscaling_target" "bbva_notifications" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.bbva_notifications.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.bbva_notification_min_container_count}"
  max_capacity       = "${var.bbva_notification_max_container_count}"
}

resource "aws_appautoscaling_policy" "bbva_notifications_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-bbva-notifications-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.bbva_notifications.name}"
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

  depends_on = ["aws_appautoscaling_target.bbva_notifications"]
}

resource "aws_appautoscaling_policy" "bbva_notifications_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-bbva-notifications-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.bbva_notifications.name}"
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

  depends_on = ["aws_appautoscaling_target.bbva_notifications"]
}

resource "aws_cloudwatch_metric_alarm" "bbva_notifications_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-bbva-notifications-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.bbva_notifications.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.bbva_notifications_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "bbva_notifications_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-bbva-notifications-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.bbva_notifications.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.bbva_notifications_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "bbva_notifications_low_container_count" {
  count               = "${var.bbva_notification_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-bbva-notifications-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.bbva_notification_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.bbva_notifications.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "bbva_notifications_high_container_count" {
  count               = "${var.bbva_notification_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-bbva-notifications-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.bbva_notification_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.bbva_notifications.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
