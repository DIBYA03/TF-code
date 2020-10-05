resource "aws_cloudwatch_log_group" "shopify_order" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-shopify-order"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "/ecs/${module.naming.aws_cloudwatch_log_group}-shopify-order"
    Team        = "${var.team}"
  }
}

data "template_file" "shopify_order" {
  template = "${file("definitions/tasks/shopify_order.json")}"

  vars {
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-shopify-order"
    aws_region    = "${var.aws_region}"

    api_env = "${var.environment}"
    aws_s3_bucket_document = "${data.aws_s3_bucket.documents.id}"
    s3_ach_pull_whitelist  = "${var.s3_ach_pull_list_config_object}"

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

    ecr_image      = "${var.shopify_order_image}"
    ecr_image_tag  = "${var.shopify_order_image_tag}"
    fargate_cpu    = "${var.shopify_order_cpu}"
    fargate_memory = "${var.shopify_order_mem}"
    name           = "${var.shopify_order_name}"

    grpc_port               = "${var.grpc_port}"
    use_transaction_service = "${var.use_transaction_service}"
    use_banking_service     = "${var.use_banking_service}"
    use_invoice_service     = "${var.use_invoice_service}"

    sqs_region = "${var.aws_region}"
    aws_s3_bucket_document = "${data.aws_s3_bucket.documents.id}"

    shopify_order_sqs_url                 = "${data.aws_sqs_queue.shopify_order.id}"
    card_reader_max_money_request_allowed = "${var.card_reader_max_request_amount}"
    card_online_max_money_request_allowed = "${var.card_online_max_request_amount}"

    wise_clearing_account_id  = "${data.aws_ssm_parameter.wise_clearing_account_id.arn}"
    wise_clearing_business_id = "${data.aws_ssm_parameter.wise_clearing_business_id.arn}"
    wise_clearing_user_id     = "${data.aws_ssm_parameter.wise_clearing_user_id.arn}"
    wise_support_email        = "${data.aws_ssm_parameter.wise_support_email_address.arn}"
    wise_support_name         = "${data.aws_ssm_parameter.wise_support_email_name.arn}"
    wise_invoice_email        = "${data.aws_ssm_parameter.wise_invoice_email_address.arn}"

    twilio_account_sid  = "${data.aws_ssm_parameter.twilio_account_sid.arn}"
    twilio_api_sid      = "${data.aws_ssm_parameter.twilio_api_sid.arn}"
    twilio_api_secret   = "${data.aws_ssm_parameter.twilio_api_secret.arn}"
    twilio_sender_phone = "${data.aws_ssm_parameter.twilio_sender_phone.arn}"
    sendgrid_api_key    = "${data.aws_ssm_parameter.sendgrid_api_key.arn}"
    payments_url        = "https://${data.aws_ssm_parameter.payments_url.value}"
  }
}

resource "aws_ecs_task_definition" "shopify_order" {
  family             = "${module.naming.aws_ecs_task_definition}-shopify-order"
  execution_role_arn = "${aws_iam_role.shopify_order_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.shopify_order_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.shopify_order_cpu}"
  memory                   = "${var.shopify_order_mem}"
  container_definitions    = "${data.template_file.shopify_order.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-shopify-order"
    Team        = "${var.team}"
  }
}

resource "aws_ecs_service" "shopify_order" {
  name            = "${module.naming.aws_ecs_service}-shopify-order"
  cluster         = "${aws_ecs_cluster.default.id}"
  task_definition = "${aws_ecs_task_definition.shopify_order.arn}"
  desired_count   = "${var.shopify_order_desired_container_count}"
  launch_type     = "FARGATE"

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = "200"
  deployment_minimum_healthy_percent = "100"

  network_configuration {
    security_groups = ["${aws_security_group.shopify_order_ecs.id}"]
    subnets         = ["${var.app_subnet_ids}"]
  }
}

resource "aws_appautoscaling_target" "shopify_order" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.shopify_order.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${local.ecs_autoscaling_role}"
  min_capacity       = "${var.shopify_order_min_container_count}"
  max_capacity       = "${var.shopify_order_max_container_count}"
}

resource "aws_appautoscaling_policy" "shopify_order_up" {
  name               = "${module.naming.aws_appautoscaling_policy}-shopify-order-up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.shopify_order.name}"
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

  depends_on = ["aws_appautoscaling_target.shopify_order"]
}

resource "aws_appautoscaling_policy" "shopify_order_down" {
  name               = "${module.naming.aws_appautoscaling_policy}-shopify-order-down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.default.name}/${aws_ecs_service.shopify_order.name}"
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

  depends_on = ["aws_appautoscaling_target.shopify_order"]
}

resource "aws_cloudwatch_metric_alarm" "shopify_order_cpu_high" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-shopify-order-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.shopify_order.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.shopify_order_up.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "shopify_order_cpu_low" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-shopify-order-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "10"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.shopify_order.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.shopify_order_down.arn}"]
}

resource "aws_cloudwatch_metric_alarm" "shopify_order_low_container_count" {
  count               = "${var.shopify_order_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-shopify-order-low-ecs-count"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.shopify_order_min_container_count - 1}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.shopify_order.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}

resource "aws_cloudwatch_metric_alarm" "shopify_order_high_container_count" {
  count               = "${var.shopify_order_add_monitoring ? 1 : 0}"
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-shopify-order-high-ecs-count"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "${var.shopify_order_max_container_count}"

  dimensions {
    ClusterName = "${aws_ecs_cluster.default.name}"
    ServiceName = "${aws_ecs_service.shopify_order.name}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
