resource "aws_cloudwatch_log_group" "signature" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-signature"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-signature"
    Team        = "${var.team}"
  }
}

data "template_file" "signature" {
  template = "${file("definitions/tasks/signature.json")}"

  vars {
    name    = "${var.signature_name}"
    api_env = "${var.environment}"

    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-signature"
    aws_region    = "${var.aws_region}"

    aws_s3_bucket_document = "${data.aws_s3_bucket.documents.id}"

    ecr_image      = "${var.signature_image}"
    ecr_image_tag  = "${var.signature_image_tag}"
    fargate_cpu    = "${var.signature_cpu}"
    fargate_memory = "${var.signature_mem}"

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

resource "aws_ecs_task_definition" "signature" {
  family             = "${module.naming.aws_ecs_task_definition}-signature"
  execution_role_arn = "${aws_iam_role.signature_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.signature_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.signature_cpu}"
  memory                   = "${var.signature_mem}"
  container_definitions    = "${data.template_file.signature.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-signature"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "signature" {
  name            = "${module.naming.aws_ecs_service}-signature"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.signature.arn}"
  desired_count   = "${var.signature_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.signature_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }
}

resource "aws_appautoscaling_target" "signature" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.signature.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.signature_min_container_count}"
  max_capacity       = "${var.signature_max_container_count}"
}

resource "aws_appautoscaling_policy" "signature_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-signature-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.signature.name}"
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

  depends_on = ["aws_appautoscaling_target.signature"]
}

resource "aws_appautoscaling_policy" "signature_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-signature-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.signature.name}"
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

  depends_on = ["aws_appautoscaling_target.signature"]
}

resource "aws_cloudwatch_metric_alarm" "signature_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-signature-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.signature.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.signature_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "signature_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-signature-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.signature.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.signature_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "signature_low_container_count" {
  count               = "${var.signature_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-signature-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.signature_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.signature.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "signature_high_container_count" {
  count               = "${var.signature_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-signature-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.signature_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.signature.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
