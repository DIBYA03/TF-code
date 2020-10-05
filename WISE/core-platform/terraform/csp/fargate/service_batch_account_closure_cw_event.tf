resource "aws_iam_role" "account_closure" {
  name = "${module.naming.aws_iam_role}-account-closure-cw-event"

  assume_role_policy = <<DOC
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "events.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
DOC
}

resource "aws_iam_role_policy" "account_closure_run_task_with_any_role" {
  name = "account_closure_run_task_with_any_role"
  role = "${aws_iam_role.account_closure.id}"

  policy = <<DOC
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "iam:PassRole",
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "ecs:RunTask",
            "Resource": "${replace(aws_ecs_task_definition.batch_account_closure.arn, "/:\\d+$/", ":*")}"
        }
    ]
}
DOC
}

resource "aws_cloudwatch_event_rule" "account_closure" {
  name                = "${module.naming.aws_cloudwatch_event_rule}-account-closure"
  description         = "run batch analytics daily"
  schedule_expression = "${var.account_closure_cw_event_scheduled_expression}"
}

resource "aws_cloudwatch_event_target" "account_closure_ecs_scheduled_task" {
  arn      = "${aws_ecs_cluster.default.arn}"
  rule     = "${aws_cloudwatch_event_rule.account_closure.name}"
  role_arn = "${aws_iam_role.account_closure.arn}"

  ecs_target {
    task_count          = 1
    task_definition_arn = "${aws_ecs_task_definition.batch_account_closure.arn}"
    launch_type         = "FARGATE"

    network_configuration = {
      subnets          = ["${var.app_subnet_ids}"]
      security_groups  = ["${aws_security_group.batch_account_closure_ecs.id}"]
      assign_public_ip = false
    }
  }
}
