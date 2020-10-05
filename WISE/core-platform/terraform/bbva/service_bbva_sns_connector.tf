resource "aws_cloudwatch_log_group" "bbva_sns_connector" {
  name = "/ecs/${module.naming.aws_cloudwatch_log_group}-bbva-sns-connector"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-bbva-sns-connector"
    Team        = "${var.team}"
  }
}

data "template_file" "bbva_sns_connector" {
  template = "${file("definitions/tasks/bbva_sns_connector.json")}"

  vars {
    name          = "${var.bbva_sns_connector_name}"
    api_env       = "${var.environment}"
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-bbva-sns-connector"
    aws_region    = "${var.aws_region}"

    ecr_image      = "${var.bbva_sns_connector_image}"
    ecr_image_tag  = "${var.bbva_sns_connector_image_tag}"
    fargate_cpu    = "${var.bbva_sns_connector_cpu}"
    fargate_memory = "${var.bbva_sns_connector_mem}"

    sns_bbva_arn    = "${aws_sns_topic.bbva_notifications.arn}"
    sns_bbva_region = "${var.aws_region}"

    sqs_bbva_region = "${var.aws_region}"
    sqs_bbva_url    = "${aws_sqs_queue.bbva_notifications.id}"
  }
}

resource "aws_ecs_task_definition" "bbva_sns_connector" {
  family             = "${module.naming.aws_ecs_task_definition}-bbva-sns-connector"
  execution_role_arn = "${aws_iam_role.bbva_sns_connector_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.bbva_sns_connector_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.bbva_sns_connector_cpu}"
  memory                   = "${var.bbva_sns_connector_mem}"
  container_definitions    = "${data.template_file.bbva_sns_connector.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-bbva-sns-connector"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "bbva_sns_connector" {
  name            = "${module.naming.aws_ecs_service}-bbva-sns-connector"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.bbva_sns_connector.arn}"
  desired_count   = "${var.bbva_sns_connector_desired_container_count}"
  launch_type     = "FARGATE"

  network_configuration {
    security_groups = ["${aws_security_group.bbva_sns_connector_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }
}

resource "aws_appautoscaling_target" "bbva_sns_connector" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.bbva_sns_connector.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "arn:aws:iam::${data.aws_caller_identity.account.account_id}:role/aws-service-role/ecs.application-autoscaling.amazonaws.com/AWSServiceRoleForApplicationAutoScaling_ECSService"
  min_capacity       = "${var.bbva_sns_connector_min_container_count}"
  max_capacity       = "${var.bbva_sns_connector_max_container_count}"
}

resource "aws_appautoscaling_policy" "bbva_sns_connector_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-bbva-sns-connector-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.bbva_sns_connector.name}"
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

  depends_on = ["aws_appautoscaling_target.bbva_sns_connector"]
}

resource "aws_appautoscaling_policy" "bbva_sns_connector_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-bbva-sns-connector-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.bbva_sns_connector.name}"
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

  depends_on = ["aws_appautoscaling_target.bbva_sns_connector"]
}

resource "aws_cloudwatch_metric_alarm" "bbva_sns_connector_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-bbva-sns-connector-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.bbva_sns_connector.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.bbva_sns_connector_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "bbva_sns_connector_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-bbva-sns-connector-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.bbva_sns_connector.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.bbva_sns_connector_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "bbva_sns_connector_low_container_count" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-bbva-sns-connector-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.bbva_sns_connector_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.bbva_sns_connector.name}"
  }

  alarm_actions = [
    "${var.sns_critical_topic}",
  ]

  ok_actions = [
    "${var.sns_critical_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "bbva_sns_connector_high_container_count" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-bbva-sns-connector-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.bbva_sns_connector_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.bbva_sns_connector.name}"
  }

  alarm_actions = [
    "${var.sns_critical_topic}",
  ]

  ok_actions = [
    "${var.sns_critical_topic}",
  ]
}
