resource "aws_cloudwatch_log_group" "csp_review_default" {
  name = "/ecs/${module.naming.aws_cloudwatch_log_group}-csp-review"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-csp-review"
    Team        = "${var.team}"
  }
}

data "template_file" "csp_review" {
  template = "${file("definitions/tasks/csp_review.json")}"

  vars {
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-csp-review"
    aws_region    = "${var.aws_region}"

    ecr_image      = "${var.csp_review_image}"
    ecr_image_tag  = "${var.csp_review_image_tag}"
    fargate_cpu    = "${var.csp_review_cpu}"
    fargate_memory = "${var.csp_review_mem}"

    name = "${var.csp_review_name}"

    aws_s3_bucket_document = "${data.aws_s3_bucket.documents.id}"
    s3_ach_pull_whitelist  = "${var.s3_ach_pull_list_config_object}"

    bbva_app_env    = "${data.aws_ssm_parameter.bbva_app_env.arn}"
    bbva_app_id     = "${data.aws_ssm_parameter.bbva_app_id.arn}"
    bbva_app_name   = "${data.aws_ssm_parameter.bbva_app_name.arn}"
    bbva_app_secret = "${data.aws_ssm_parameter.bbva_app_secret.arn}"

    bank_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    bank_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    bank_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    bank_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    bank_db_name       = "${data.aws_ssm_parameter.bank_rds_db_name.arn}"
    bank_db_user       = "${data.aws_ssm_parameter.bank_rds_user_name.arn}"
    bank_db_passwd     = "${data.aws_ssm_parameter.bank_rds_password.arn}"

    core_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    core_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    core_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    core_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    core_db_name       = "${data.aws_ssm_parameter.core_rds_db_name.arn}"
    core_db_user       = "${data.aws_ssm_parameter.core_rds_user_name.arn}"
    core_db_passwd     = "${data.aws_ssm_parameter.core_rds_password.arn}"

    csp_db_read_url   = "${data.aws_ssm_parameter.csp_rds_read_endpoint.arn}"
    csp_db_read_port  = "${data.aws_ssm_parameter.csp_rds_port.arn}"
    csp_db_write_url  = "${data.aws_ssm_parameter.csp_rds_master_endpoint.arn}"
    csp_db_write_port = "${data.aws_ssm_parameter.csp_rds_port.arn}"
    csp_db_name       = "${data.aws_ssm_parameter.csp_rds_db_name.arn}"
    csp_db_user       = "${data.aws_ssm_parameter.csp_rds_username.arn}"
    csp_db_passwd     = "${data.aws_ssm_parameter.csp_rds_password.arn}"

    identity_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    identity_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    identity_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    identity_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    identity_db_name       = "${data.aws_ssm_parameter.identity_rds_db_name.arn}"
    identity_db_user       = "${data.aws_ssm_parameter.identity_rds_user_name.arn}"
    identity_db_passwd     = "${data.aws_ssm_parameter.identity_rds_password.arn}"

    csp_review_sqs_region = "${var.aws_region}"
    csp_review_sqs_url    = "${data.aws_sqs_queue.review.url}"

    sqs_region     = "${var.aws_region}"
    sqs_url        = "${data.aws_sqs_queue.internal_banking.url}"
    sqs_app_region = "${var.aws_region}"

    sendgrid_api_key = "${data.aws_ssm_parameter.sendgrid_api_key.arn}"

    csp_notification_slack_channel = "${data.aws_ssm_parameter.csp_notification_slack_channel.arn}"
    csp_notification_slack_url     = "${data.aws_ssm_parameter.csp_notification_slack_url.arn}"

    sqs_region      = "${var.aws_region}"
    segment_sqs_url = "${data.aws_sqs_queue.segment_analytics.id}"

    grpc_port       = "${var.grpc_port}"
  }
}

resource "aws_ecs_task_definition" "csp_review" {
  family             = "${module.naming.aws_ecs_task_definition}-csp-review"
  execution_role_arn = "${aws_iam_role.csp_review_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.csp_review_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.csp_review_cpu}"
  memory                   = "${var.csp_review_mem}"
  container_definitions    = "${data.template_file.csp_review.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-csp-review"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "csp_review" {
  name            = "${module.naming.aws_ecs_service}-csp-review"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.csp_review.arn}"
  desired_count   = "${var.csp_review_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.csp_review_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }
}

resource "aws_appautoscaling_target" "csp_review" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.csp_review.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.csp_review_min_container_count}"
  max_capacity       = "${var.csp_review_max_container_count}"
}

resource "aws_appautoscaling_policy" "csp_review_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-csp-review-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.csp_review.name}"
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

  depends_on = ["aws_appautoscaling_target.csp_review"]
}

resource "aws_appautoscaling_policy" "csp_review_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-csp-review-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.csp_review.name}"
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

  depends_on = ["aws_appautoscaling_target.csp_review"]
}

resource "aws_cloudwatch_metric_alarm" "csp_review_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-csp-review-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.csp_review.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.csp_review_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "csp_review_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-csp-review-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.csp_review.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.csp_review_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "csp_review_low_container_count" {
  count               = "${var.csp_review_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-csp-review-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.csp_review_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.csp_review.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "csp_review_high_container_count" {
  count               = "${var.csp_review_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-csp-review-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.csp_review_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.csp_review.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
