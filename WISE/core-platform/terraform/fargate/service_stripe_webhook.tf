resource "aws_cloudwatch_log_group" "stripe_webhook_default" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-stripe-webhook"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-stripe-webhook"
    Team        = "${var.team}"
  }
}

data "template_file" "stripe_webhook" {
  template = "${file("definitions/tasks/stripe_webhook.json")}"

  vars {
    name          = "${var.stripe_webhook_name}"
    api_env       = "${var.environment}"
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-stripe-webhook"
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

    identity_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    identity_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    identity_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    identity_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    identity_db_name       = "${data.aws_ssm_parameter.identity_rds_db_name.arn}"
    identity_db_user       = "${data.aws_ssm_parameter.identity_rds_user_name.arn}"
    identity_db_passwd     = "${data.aws_ssm_parameter.identity_rds_password.arn}"

    container_port = "${var.services_container_port}"
    host_port      = "${var.services_container_port}"

    wise_support_email = "${data.aws_ssm_parameter.wise_support_email_address.arn}"
    wise_support_name  = "${data.aws_ssm_parameter.wise_support_email_name.arn}"

    ecr_image      = "${var.stripe_webhook_image}"
    ecr_image_tag  = "${var.stripe_webhook_image_tag}"
    fargate_cpu    = "${var.stripe_webhook_cpu}"
    fargate_memory = "${var.stripe_webhook_mem}"

    grpc_port               = "${var.grpc_port}"
    use_transaction_service = "${var.use_transaction_service}"
    use_banking_service     = "${var.use_banking_service}"
    use_invoice_service     = "${var.use_invoice_service}"

    payments_sqs_url = "${data.aws_sqs_queue.stripe_webhook.url}"
  }
}

resource "aws_ecs_task_definition" "stripe_webhook" {
  family             = "${module.naming.aws_ecs_task_definition}-stripe-webhook"
  execution_role_arn = "${aws_iam_role.stripe_webhook_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.stripe_webhook_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.stripe_webhook_cpu}"
  memory                   = "${var.stripe_webhook_mem}"
  container_definitions    = "${data.template_file.stripe_webhook.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-stripe-webhook"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "stripe_webhook" {
  name            = "${module.naming.aws_ecs_service}-stripe-webhook"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.stripe_webhook.arn}"
  desired_count   = "${var.stripe_webhook_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.stripe_webhook_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.stripe_webhook.id}"
    container_name   = "${var.stripe_webhook_name}"
    container_port   = "${var.services_container_port}"
  }

  depends_on = [
    "aws_alb_listener.services_http",
    "aws_alb_listener.services_https",
  ]
}

resource "aws_appautoscaling_target" "stripe_webhook" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.stripe_webhook.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.stripe_webhook_min_container_count}"
  max_capacity       = "${var.stripe_webhook_max_container_count}"
}

resource "aws_appautoscaling_policy" "stripe_webhook_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-stripe-webhook-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.stripe_webhook.name}"
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

  depends_on = ["aws_appautoscaling_target.stripe_webhook"]
}

resource "aws_appautoscaling_policy" "stripe_webhook_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-stripe-webhook-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.stripe_webhook.name}"
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

  depends_on = ["aws_appautoscaling_target.stripe_webhook"]
}

resource "aws_cloudwatch_metric_alarm" "stripe_webhook_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-stripe-webhook-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.stripe_webhook.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.stripe_webhook_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "stripe_webhook_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-stripe-webhook-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.stripe_webhook.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.stripe_webhook_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "stripe_webhook_low_container_count" {
  count               = "${var.stripe_webhook_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-stripe-webhook-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.stripe_webhook_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.stripe_webhook.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "stripe_webhook_high_container_count" {
  count               = "${var.stripe_webhook_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-stripe-webhook-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.stripe_webhook_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.stripe_webhook.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
