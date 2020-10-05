resource "aws_iam_role" "batch_analytics" {
  name = "${module.naming.aws_iam_role}-batch-analytics-cw-event"

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

resource "aws_iam_role_policy" "batch_analytics_run_task_with_any_role" {
  name = "batch_analytics_run_task_with_any_role"
  role = "${aws_iam_role.batch_analytics.id}"

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
            "Resource": "${replace(aws_ecs_task_definition.batch_analytics.arn, "/:\\d+$/", ":*")}"
        }
    ]
}
DOC
}

resource "aws_cloudwatch_event_rule" "batch_analytics" {
  name                = "${module.naming.aws_cloudwatch_event_rule}-batch-analytics"
  description         = "run batch analytics daily"
  schedule_expression = "${var.batch_analytics_cw_event_scheduled_expression}"
}

resource "aws_cloudwatch_event_target" "batch_analytics_ecs_scheduled_task" {
  arn      = "${aws_ecs_cluster.default.arn}"
  rule     = "${aws_cloudwatch_event_rule.batch_analytics.name}"
  role_arn = "${aws_iam_role.batch_analytics.arn}"

  ecs_target {
    task_count          = 1
    task_definition_arn = "${aws_ecs_task_definition.batch_analytics.arn}"
    launch_type         = "FARGATE"

    network_configuration = {
      subnets          = ["${var.app_subnet_ids}"]
      security_groups  = ["${aws_security_group.batch_analytics_ecs.id}"]
      assign_public_ip = false
    }
  }
}
