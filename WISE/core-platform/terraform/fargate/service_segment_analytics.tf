resource "aws_cloudwatch_log_group" "segment_analytics" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-segment-analytics"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "/ecs/${module.naming.aws_cloudwatch_log_group}-segment-analytics"
    Team        = "${var.team}"
  }
}

data "template_file" "segment_analytics" {
  template = "${file("definitions/tasks/segment_analytics.json")}"

  vars {
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-segment-analytics"
    aws_region    = "${var.aws_region}"

    api_env = "${var.environment}"

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

    ecr_image      = "${var.segment_analytics_image}"
    ecr_image_tag  = "${var.segment_analytics_image_tag}"
    fargate_cpu    = "${var.segment_analytics_cpu}"
    fargate_memory = "${var.segment_analytics_mem}"
    name           = "${var.segment_analytics_name}"

    grpc_port               = "${var.grpc_port}"
    use_transaction_service = "${var.use_transaction_service}"
    use_banking_service     = "${var.use_banking_service}"

    sqs_region = "${var.aws_region}"

    segment_sqs_url   = "${data.aws_sqs_queue.segment_analytics.id}"
    segment_write_key = "${data.aws_ssm_parameter.segment_write_key.arn}"
  }
}

resource "aws_ecs_task_definition" "segment_analytics" {
  family             = "${module.naming.aws_ecs_task_definition}-segment-analytics"
  execution_role_arn = "${aws_iam_role.segment_analytics_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.segment_analytics_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.segment_analytics_cpu}"
  memory                   = "${var.segment_analytics_mem}"
  container_definitions    = "${data.template_file.segment_analytics.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-segment-analytics"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "segment_analytics" {
  name            = "${module.naming.aws_ecs_service}-segment-analytics"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.segment_analytics.arn}"
  desired_count   = "${var.segment_analytics_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.segment_analytics_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }
}

resource "aws_appautoscaling_target" "segment_analytics" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.segment_analytics.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.segment_analytics_min_container_count}"
  max_capacity       = "${var.segment_analytics_max_container_count}"
}

resource "aws_appautoscaling_policy" "segment_analytics_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-segment-analytics-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.segment_analytics.name}"
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

  depends_on = ["aws_appautoscaling_target.segment_analytics"]
}

resource "aws_appautoscaling_policy" "segment_analytics_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-segment-analytics-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.segment_analytics.name}"
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

  depends_on = ["aws_appautoscaling_target.segment_analytics"]
}

resource "aws_cloudwatch_metric_alarm" "segment_analytics_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-segment-analytics-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.segment_analytics.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.segment_analytics_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "segment_analytics_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-segment-analytics-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.segment_analytics.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.segment_analytics_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "segment_analytics_low_container_count" {
  count               = "${var.segment_analytics_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-segment-analytics-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.segment_analytics_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.segment_analytics.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "segment_analytics_high_container_count" {
  count               = "${var.segment_analytics_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-segment-analytics-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.segment_analytics_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.segment_analytics.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
