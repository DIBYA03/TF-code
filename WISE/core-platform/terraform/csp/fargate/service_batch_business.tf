resource "aws_cloudwatch_log_group" "batch_business" {
  name = "/ecs/${module.naming.aws_cloudwatch_log_group}-batch-business"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "/ecs/${module.naming.aws_cloudwatch_log_group}-batch-business"
    Team        = "${var.team}"
  }
}

data "template_file" "batch_business" {
  template = "${file("definitions/tasks/batch_business.json")}"

  vars {
    aws_log_group = "${module.naming.aws_cloudwatch_log_group}-batch-business"
    aws_region    = "${var.aws_region}"

    batch_tz = "${var.batch_default_timezone}"

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

    csp_db_read_url   = "${data.aws_ssm_parameter.csp_rds_read_endpoint.arn}"
    csp_db_read_port  = "${data.aws_ssm_parameter.csp_rds_port.arn}"
    csp_db_write_url  = "${data.aws_ssm_parameter.csp_rds_master_endpoint.arn}"
    csp_db_write_port = "${data.aws_ssm_parameter.csp_rds_port.arn}"
    csp_db_name       = "${data.aws_ssm_parameter.csp_rds_db_name.arn}"
    csp_db_user       = "${data.aws_ssm_parameter.csp_rds_username.arn}"
    csp_db_passwd     = "${data.aws_ssm_parameter.csp_rds_password.arn}"

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

    sendgrid_api_key = "${data.aws_ssm_parameter.sendgrid_api_key.arn}"

    csp_review_sqs_region = "${var.aws_region}"
    csp_review_sqs_url    = "${data.aws_sqs_queue.review.url}"

    wise_support_email = "${data.aws_ssm_parameter.wise_support_email_address.value}"
    wise_support_name  = "${data.aws_ssm_parameter.wise_support_email_name.value}"

    csp_notification_slack_channel = "${data.aws_ssm_parameter.csp_notification_slack_channel.arn}"
    csp_notification_slack_url     = "${data.aws_ssm_parameter.csp_notification_slack_url.arn}"

    sqs_region      = "${var.aws_region}"
    segment_sqs_url = "${data.aws_sqs_queue.segment_analytics.id}"

    ecr_image      = "${var.batch_business_image}"
    ecr_image_tag  = "${var.batch_business_image_tag}"
    fargate_cpu    = "${var.batch_business_cpu}"
    fargate_memory = "${var.batch_business_mem}"
    name           = "${var.batch_business_name}"
  }
}

resource "aws_ecs_task_definition" "batch_business" {
  family             = "${module.naming.aws_ecs_task_definition}-batch-business"
  execution_role_arn = "${aws_iam_role.batch_business_execution_role.arn}"       # Service management
  task_role_arn      = "${aws_iam_role.batch_business_execution_role.arn}"       # To run AWS services

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "${var.batch_business_cpu}"
  memory                   = "${var.batch_business_mem}"
  container_definitions    = "${data.template_file.batch_business.rendered}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_ecs_task_definition}-batch-business"
    Team        = "${var.team}"
  }
}
