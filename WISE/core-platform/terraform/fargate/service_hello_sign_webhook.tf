resource "aws_cloudwatch_log_group" "hello_sign_default" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-hello-sign"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-hello-sign"
    Team        = "${var.team}"
  }
}

data "template_file" "hello_sign" {
  template = "${file("definitions/tasks/hello_sign.json")}"

  vars {
    name          = "${var.hello_sign_name}"
    api_env       = "${var.environment}"
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-hello-sign"
    aws_region    = "${var.aws_region}"

    aws_s3_bucket_document = "${data.aws_s3_bucket.documents.id}"

    container_port = "${var.services_container_port}"
    host_port      = "${var.services_container_port}"

    ecr_image      = "${var.hello_sign_image}"
    ecr_image_tag  = "${var.hello_sign_image_tag}"
    fargate_cpu    = "${var.hello_sign_cpu}"
    fargate_memory = "${var.hello_sign_mem}"

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

    signature_sqs_url = "${data.aws_sqs_queue.signature_webhook.url}"

    hellosign_api_key = "${data.aws_ssm_parameter.hellosign_api_key.arn}"
  }
}

resource "aws_ecs_task_definition" "hello_sign" {
  family             = "${module.naming.aws_ecs_task_definition}-hello-sign"
  execution_role_arn = "${aws_iam_role.hello_sign_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.hello_sign_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.hello_sign_cpu}"
  memory                   = "${var.hello_sign_mem}"
  container_definitions    = "${data.template_file.hello_sign.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-hello-sign"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "hello_sign" {
  name            = "${module.naming.aws_ecs_service}-hello-sign"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.hello_sign.arn}"
  desired_count   = "${var.hello_sign_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.hello_sign_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.hello_sign.id}"
    container_name   = "${var.hello_sign_name}"
    container_port   = "${var.services_container_port}"
  }

  depends_on = [
    "aws_alb_listener.services_http",
    "aws_alb_listener.services_https",
  ]
}

resource "aws_appautoscaling_target" "hello_sign" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.hello_sign.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.hello_sign_min_container_count}"
  max_capacity       = "${var.hello_sign_max_container_count}"
}

resource "aws_appautoscaling_policy" "hello_sign_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-hello-sign-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.hello_sign.name}"
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

  depends_on = ["aws_appautoscaling_target.hello_sign"]
}

resource "aws_appautoscaling_policy" "hello_sign_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-hello-sign-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.hello_sign.name}"
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

  depends_on = ["aws_appautoscaling_target.hello_sign"]
}

resource "aws_cloudwatch_metric_alarm" "hello_sign_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-hello-sign-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.hello_sign.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.hello_sign_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "hello_sign_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-hello-sign-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.hello_sign.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.hello_sign_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "hello_sign_low_container_count" {
  count               = "${var.hello_sign_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-hello-sign-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.hello_sign_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.hello_sign.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "hello_sign_high_container_count" {
  count               = "${var.hello_sign_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-hello-sign-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.hello_sign_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.hello_sign.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
