resource "aws_cloudwatch_log_group" "merchant_logo_default" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-merchant-logo"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-merchant-logo"
    Team        = "${var.team}"
  }
}

data "template_file" "merchant_logo" {
  template = "${file("definitions/tasks/merchant_logo.json")}"

  vars {
    name = "${var.merchant_logo_name}"

    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-merchant-logo"
    aws_region    = "${var.aws_region}"

    ecr_image      = "${var.merchant_logo_image}"
    ecr_image_tag  = "${var.merchant_logo_image_tag}"
    fargate_cpu    = "${var.merchant_logo_cpu}"
    fargate_memory = "${var.merchant_logo_mem}"
    container_port = "${var.services_container_port}"
    host_port      = "${var.services_container_port}"

    clear_bit_api_key = "${data.aws_ssm_parameter.clear_bit_api_key.arn}"
  }
}

resource "aws_ecs_task_definition" "merchant_logo" {
  family             = "${module.naming.aws_ecs_task_definition}-merchant-logo"
  execution_role_arn = "${aws_iam_role.merchant_logo_execution.arn}"
  task_role_arn      = "${aws_iam_role.merchant_logo_execution.arn}"

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.merchant_logo_cpu}"
  memory                   = "${var.merchant_logo_mem}"
  container_definitions    = "${data.template_file.merchant_logo.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-merchant-logo"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "merchant_logo" {
  name            = "${module.naming.aws_ecs_service}-merchant-logo"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.merchant_logo.arn}"
  desired_count   = "${var.merchant_logo_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.merchant_logo_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.merchant_logo.id}"
    container_name   = "${var.merchant_logo_name}"
    container_port   = "${var.services_container_port}"
  }

  depends_on = [
    "aws_alb_listener.services_http",
    "aws_alb_listener.services_https",
  ]
}

resource "aws_appautoscaling_target" "merchant_logo" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.merchant_logo.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.merchant_logo_min_container_count}"
  max_capacity       = "${var.merchant_logo_max_container_count}"
}

resource "aws_appautoscaling_policy" "merchant_logo_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-merchant-logo-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.merchant_logo.name}"
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

  depends_on = ["aws_appautoscaling_target.merchant_logo"]
}

resource "aws_appautoscaling_policy" "merchant_logo_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-merchant-logo-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.merchant_logo.name}"
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

  depends_on = ["aws_appautoscaling_target.merchant_logo"]
}

resource "aws_cloudwatch_metric_alarm" "merchant_logo_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-merchant-logo-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.merchant_logo.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.merchant_logo_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "merchant_logo_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-merchant-logo-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.merchant_logo.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.merchant_logo_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "merchant_logo_low_container_count" {
  count               = "${var.merchant_logo_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-merchant-logo-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.merchant_logo_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.merchant_logo.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "merchant_logo_high_container_count" {
  count               = "${var.merchant_logo_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-merchant-logo-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.merchant_logo_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.merchant_logo.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
