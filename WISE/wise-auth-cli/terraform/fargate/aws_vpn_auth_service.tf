resource "aws_cloudwatch_log_group" "aws_vpn_auth_default" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-aws"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-aws"
    Team        = "${var.team}"
  }
}

data "template_file" "aws_vpn_auth" {
  template = "${file("definitions/tasks/aws_vpn_auth.json")}"

  vars {
    name          = "${var.aws_vpn_auth_name}"
    aws_log_group = "/ecs/${module.naming.aws_cloudwatch_log_group}-aws"
    aws_region    = "${var.aws_region}"

    container_port = "${var.aws_vpn_auth_container_port}"
    host_port      = "${var.aws_vpn_auth_container_port}"

    ecr_image      = "${aws_ecr_repository.aws_vpn_auth.repository_url}"
    ecr_image_tag  = "${var.aws_vpn_auth_image_tag}"
    fargate_cpu    = "${var.aws_vpn_auth_cpu}"
    fargate_memory = "${var.aws_vpn_auth_mem}"

    google_idp_url_ssm_arn = "${data.aws_ssm_parameter.google_idp_url.arn}"
  }
}

resource "aws_ecs_task_definition" "aws_vpn_auth" {
  family             = "${module.naming.aws_ecs_task_definition}-aws"
  execution_role_arn = "${aws_iam_role.aws_vpn_auth_execution_role.arn}"
  task_role_arn      = "${aws_iam_role.aws_vpn_auth_task_role.arn}"

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.aws_vpn_auth_cpu}"
  memory                   = "${var.aws_vpn_auth_mem}"
  container_definitions    = "${data.template_file.aws_vpn_auth.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-aws"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "aws_vpn_auth" {
  name            = "${module.naming.aws_ecs_service}-aws"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.aws_vpn_auth.arn}"
  desired_count   = "${var.aws_vpn_auth_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.aws_vpn_auth_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }

  load_balancer {
    target_group_arn = "${aws_lb_target_group.aws_vpn_auth.id}"
    container_name   = "${var.aws_vpn_auth_name}"
    container_port   = "${var.aws_vpn_auth_container_port}"
  }
}

resource "aws_appautoscaling_target" "aws_vpn_auth" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.aws_vpn_auth.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.aws_vpn_auth_min_container_count}"
  max_capacity       = "${var.aws_vpn_auth_max_container_count}"
}

resource "aws_appautoscaling_policy" "aws_vpn_auth_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-aws-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.aws_vpn_auth.name}"
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

  depends_on = ["aws_appautoscaling_target.aws_vpn_auth"]
}

resource "aws_appautoscaling_policy" "aws_vpn_auth_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-aws-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.aws_vpn_auth.name}"
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

  depends_on = ["aws_appautoscaling_target.aws_vpn_auth"]
}

resource "aws_cloudwatch_metric_alarm" "aws_vpn_auth_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-aws-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.aws_vpn_auth.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.aws_vpn_auth_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "aws_vpn_auth_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-aws-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.aws_vpn_auth.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.aws_vpn_auth_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "aws_vpn_auth_low_container_count" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-aws-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.aws_vpn_auth_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.aws_vpn_auth.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "aws_vpn_auth_high_container_count" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-aws-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.aws_vpn_auth_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.aws_vpn_auth.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
