resource "aws_iam_role" "batch_daily" {
  name = "${module.naming.aws_iam_role}-batch-daily-cw-event"

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

resource "aws_iam_role_policy" "batch_daily_run_task_with_any_role" {
  name = "batch_daily_run_task_with_any_role"
  role = "${aws_iam_role.batch_daily.id}"

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
      "Action": "states:StartExecution",
      "Resource": "${aws_sfn_state_machine.batch_daily.id}"
    }
  ]
}
DOC
}

resource "aws_cloudwatch_event_rule" "batch_daily" {
  name                = "${module.naming.aws_cloudwatch_event_rule}-batch-daily"
  description         = "run batch account, transaction, and monitor 2nd to 31st of every month"
  schedule_expression = "${var.batch_daily_cw_event_scheduled_expression}"
}

resource "aws_cloudwatch_event_target" "batch_daily_ecs_scheduled_task" {
  arn      = "${aws_sfn_state_machine.batch_daily.id}"
  rule     = "${aws_cloudwatch_event_rule.batch_daily.name}"
  role_arn = "${aws_iam_role.batch_daily.arn}"
}
