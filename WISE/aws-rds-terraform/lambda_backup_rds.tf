resource "null_resource" "build_lambda" {
  provisioner "local-exec" {
    command    = "env GOOS=linux GOARCH=amd64 go build -o ./lambda/rds_backup/main ./lambda/rds_backup/main.go;"
    on_failure = "fail"
  }

  triggers = {
    always_run = "${timestamp()}"
  }
}

data "archive_file" "rds_backup_lambda" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  type        = "zip"
  source_file = "./lambda/rds_backup/main"
  output_path = "./lambda/rds_backup/lambda.zip"

  depends_on = [
    "null_resource.build_lambda",
  ]
}

resource "aws_lambda_function" "rds_backup" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  function_name = "${module.naming.aws_lambda_function}-backup"
  role          = "${aws_iam_role.rds_backup_lambda.arn}"
  kms_key_arn   = "${aws_kms_key.rds_default.arn}"
  timeout       = "${var.backup_lambda_timeout}"

  filename         = "./lambda/rds_backup/lambda.zip"
  source_code_hash = "${data.archive_file.rds_backup_lambda.output_base64sha256}"

  runtime = "go1.x"
  handler = "main"

  vpc_config = {
    subnet_ids         = ["${var.app_subnet_ids}"]
    security_group_ids = ["${aws_security_group.lambda_backup.id}"]
  }

  environment {
    variables = {
      DB_INSTANCE_IDENTIFIER  = "${aws_db_instance.backup_replica.id}"
      DB_SNAPSHOT_COUNT_LIMIT = "${var.backup_db_snapshot_count_limit}"
    }
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_lambda_function}-backup"
    Team        = "${var.team}"
  }

  depends_on = [
    "data.archive_file.rds_backup_lambda",
  ]
}

# Give permissions for the API Gateway to access the lambda function
resource "aws_lambda_permission" "rds_backup" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  statement_id  = "AllowExecutionFromCWEvents"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.rds_backup.function_name}"

  principal  = "events.amazonaws.com"
  source_arn = "${aws_cloudwatch_event_rule.lambda_backup.arn}"

  depends_on = [
    "aws_lambda_function.rds_backup",
  ]
}

resource "aws_cloudwatch_event_rule" "lambda_backup" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  name                = "${module.naming.aws_cloudwatch_event_rule}"
  description         = "Backup ${terraform.workspace} database every thirty minutes"
  schedule_expression = "rate(30 minutes)"
}

resource "aws_cloudwatch_event_target" "lambda_backup" {
  count     = "${var.rds_enable_backup_lambda ? 1 : 0}"
  rule      = "${aws_cloudwatch_event_rule.lambda_backup.name}"
  target_id = "check_foo"
  arn       = "${aws_lambda_function.rds_backup.arn}"
}
