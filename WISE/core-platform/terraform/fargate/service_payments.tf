resource "aws_cloudwatch_log_group" "payments_default" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-payments"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_cloudwatch_log_group}-payments"
    Team        = "${var.team}"
  }
}

data "template_file" "payments" {
  template = "${file("definitions/tasks/payments.json")}"

  vars {
    name = "${var.payments_name}"

    api_env                = "${var.environment}"
    aws_s3_bucket_document = "${data.aws_s3_bucket.documents.id}"
    s3_ach_pull_whitelist  = "${var.s3_ach_pull_list_config_object}"

    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-payments"
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

    txn_db_write_url  = "${data.aws_ssm_parameter.rds_master_endpoint.arn}"
    txn_db_read_url   = "${data.aws_ssm_parameter.rds_read_endpoint.arn}"
    txn_db_write_port = "${data.aws_ssm_parameter.rds_port.arn}"
    txn_db_read_port  = "${data.aws_ssm_parameter.rds_port.arn}"
    txn_db_name       = "${data.aws_ssm_parameter.txn_rds_db_name.arn}"
    txn_db_user       = "${data.aws_ssm_parameter.txn_rds_user_name.arn}"
    txn_db_passwd     = "${data.aws_ssm_parameter.txn_rds_password.arn}"

    ecr_image      = "${var.payments_image}"
    ecr_image_tag  = "${var.payments_image_tag}"
    fargate_cpu    = "${var.payments_cpu}"
    fargate_memory = "${var.payments_mem}"
    container_port = "${var.services_container_port}"
    host_port      = "${var.services_container_port}"

    grpc_port               = "${var.grpc_port}"
    use_transaction_service = "${var.use_transaction_service}"
    use_banking_service     = "${var.use_banking_service}"
    use_invoice_service     = "${var.use_invoice_service}"

    maintenance_enabled = "${var.payments_maintenance_enabled}"
    payments_sqs_url    = "${data.aws_sqs_queue.stripe_webhook.url}"
    payments_url        = "https://${data.aws_ssm_parameter.payments_url.value}"

    sqs_banking_region = "${var.aws_region}"
    sqs_banking_url    = "${data.aws_sqs_queue.internal_banking.url}"

    stripe_key            = "${data.aws_ssm_parameter.stripe_key.arn}"
    stripe_publish_key    = "${data.aws_ssm_parameter.stripe_publish_key.arn}"
    stripe_webhook_secret = "${data.aws_ssm_parameter.stripe_webhook_secret.arn}"

    wise_clearing_account_id  = "${data.aws_ssm_parameter.wise_clearing_account_id.arn}"
    wise_clearing_business_id = "${data.aws_ssm_parameter.wise_clearing_business_id.arn}"
    wise_clearing_user_id     = "${data.aws_ssm_parameter.wise_clearing_user_id.arn}"
    wise_invoice_email        = "${data.aws_ssm_parameter.wise_invoice_email_address.arn}"

    plaid_env        = "${data.aws_ssm_parameter.plaid_env.arn}"
    plaid_public_key = "${data.aws_ssm_parameter.plaid_public_key.arn}"
    plaid_client_id  = "${data.aws_ssm_parameter.plaid_client_id.arn}"
    plaid_secret     = "${data.aws_ssm_parameter.plaid_secret.arn}"

    segment_web_write_key = "${data.aws_ssm_parameter.segment_web_write_key.arn}"
  }
}

resource "aws_ecs_task_definition" "payments" {
  family             = "${module.naming.aws_ecs_task_definition}-payments"
  execution_role_arn = "${aws_iam_role.payments_execution.arn}"
  task_role_arn      = "${aws_iam_role.payments_execution.arn}"

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.payments_cpu}"
  memory                   = "${var.payments_mem}"
  container_definitions    = "${data.template_file.payments.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-payments"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "payments" {
  name            = "${module.naming.aws_ecs_service}-payments"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.payments.arn}"
  desired_count   = "${var.payments_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.payments_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.payments.id}"
    container_name   = "${var.payments_name}"
    container_port   = "${var.services_container_port}"
  }

  depends_on = [
    "aws_alb_listener.services_http",
    "aws_alb_listener.services_https",
  ]
}

resource "aws_appautoscaling_target" "payments" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.payments.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.payments_min_container_count}"
  max_capacity       = "${var.payments_max_container_count}"
}

resource "aws_appautoscaling_policy" "payments_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-payments-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.payments.name}"
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

  depends_on = ["aws_appautoscaling_target.payments"]
}

resource "aws_appautoscaling_policy" "payments_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-payments-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.payments.name}"
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

  depends_on = ["aws_appautoscaling_target.payments"]
}

resource "aws_cloudwatch_metric_alarm" "payments_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-payments-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.payments.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.payments_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "payments_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-payments-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.payments.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.payments_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "payments_low_container_count" {
  count               = "${var.payments_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-payments-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.payments_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.payments.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "payments_high_container_count" {
  count               = "${var.payments_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-payments-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.payments_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.payments.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
