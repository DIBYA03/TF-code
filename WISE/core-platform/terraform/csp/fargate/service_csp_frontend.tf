resource "aws_cloudwatch_log_group" "csp_frontend_default" {
  name = "/ecs/${module.naming.aws_cloudwatch_log_group}-csp-frontend"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-csp-frontend"
    Team        = "${var.team}"
  }
}

data "template_file" "csp_frontend" {
  template = "${file("definitions/tasks/csp_frontend.json")}"

  vars {
    aws_log_group    = "${module.naming.aws_cloudwatch_log_group}-csp-frontend"
    aws_region       = "${var.aws_region}"
    ecr_image        = "${var.csp_frontend_image}"
    ecr_image_tag    = "${var.csp_frontend_image_tag}"
    fargate_cpu      = "${var.csp_frontend_cpu}"
    fargate_memory   = "${var.csp_frontend_mem}"
    name             = "${var.csp_frontend_name}"
    target_group_arn = "${aws_alb.csp_frontend.arn}"
    container_port   = "${var.csp_frontend_container_port}"
    host_port        = "${var.csp_frontend_container_port}"
  }
}

resource "aws_ecs_task_definition" "csp_frontend" {
  family             = "${module.naming.aws_ecs_task_definition}-csp-frontend"
  execution_role_arn = "${aws_iam_role.csp_frontend_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.csp_frontend_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.csp_frontend_cpu}"
  memory                   = "${var.csp_frontend_mem}"
  container_definitions    = "${data.template_file.csp_frontend.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-csp-frontend"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "csp_frontend" {
  name            = "${module.naming.aws_ecs_service}-csp-frontend"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.csp_frontend.arn}"
  desired_count   = "${var.csp_frontend_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.csp_frontend_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.csp_frontend.id}"
    container_name   = "${var.csp_frontend_name}"
    container_port   = "${var.csp_frontend_container_port}"
  }

  depends_on = [
    "aws_alb_listener.csp_frontend_https",
  ]
}

resource "aws_appautoscaling_target" "csp_frontend" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.csp_frontend.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.csp_frontend_min_container_count}"
  max_capacity       = "${var.csp_frontend_max_container_count}"
}

resource "aws_appautoscaling_policy" "csp_frontend_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-csp-frontend-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.csp_frontend.name}"
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

  depends_on = ["aws_appautoscaling_target.csp_frontend"]
}

resource "aws_appautoscaling_policy" "csp_frontend_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-csp-frontend-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.csp_frontend.name}"
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

  depends_on = ["aws_appautoscaling_target.csp_frontend"]
}

resource "aws_cloudwatch_metric_alarm" "csp_frontend_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-csp-frontend-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.csp_frontend.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.csp_frontend_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "csp_frontend_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-csp-frontend-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.csp_frontend.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.csp_frontend_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "csp_frontend_low_container_count" {
  count               = "${var.csp_frontend_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-csp-frontend-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.csp_frontend_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.csp_frontend.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "csp_frontend_high_container_count" {
  count               = "${var.csp_frontend_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-csp-frontend-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.csp_frontend_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.csp_frontend.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
