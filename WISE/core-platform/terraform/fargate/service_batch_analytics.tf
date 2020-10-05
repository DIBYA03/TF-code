resource "aws_cloudwatch_log_group" "batch_analytics" {
  name              = "/ecs/${module.naming.aws_cloudwatch_log_group}-batch-analytics"
  retention_in_days = "${var.cw_log_group_retention_in_days}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "/ecs/${module.naming.aws_cloudwatch_log_group}-batch-analytics"
    Team        = "${var.team}"
  }
}

data "template_file" "batch_analytics" {
  template = "${file("definitions/tasks/batch_analytics.json")}"

  vars {
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-batch-analytics"
    aws_region    = "${var.aws_region}"

    api_env = "${var.environment}"

    batch_tz = "${var.batch_default_timezone}"

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

    kinesis_trx_name   = "${var.txn_kinesis_name}"
    kinesis_trx_region = "${var.txn_kinesis_region}"

    ecr_image      = "${var.batch_analytics_image}"
    ecr_image_tag  = "${var.batch_analytics_image_tag}"
    fargate_cpu    = "${var.batch_analytics_cpu}"
    fargate_memory = "${var.batch_analytics_mem}"
    name           = "${var.batch_analytics_name}"

    grpc_port               = "${var.grpc_port}"
    use_transaction_service = "${var.use_transaction_service}"
    use_banking_service     = "${var.use_banking_service}"
    use_invoice_service     = "${var.use_invoice_service}"

    segment_write_key = "${data.aws_ssm_parameter.segment_write_key.arn}"
  }
}

resource "aws_ecs_task_definition" "batch_analytics" {
  family             = "${module.naming.aws_ecs_task_definition}-batch-analytics"
  execution_role_arn = "${aws_iam_role.batch_analytics_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.batch_analytics_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.batch_analytics_cpu}"
  memory                   = "${var.batch_analytics_mem}"
  container_definitions    = "${data.template_file.batch_analytics.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-batch-analytics"
    Team        = "${var.team}"
  }
}
