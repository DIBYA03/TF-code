resource "aws_iam_role" "batch_monthly_interest" {
  name = "${module.naming.aws_iam_role}-batch-monthly-interest-cw-event"

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

resource "aws_iam_role_policy" "batch_monthly_interest_run_task_with_any_role" {
  name = "batch_monthly_interest_run_task_with_any_role"
  role = "${aws_iam_role.batch_monthly_interest.id}"

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
      "Resource": "${aws_sfn_state_machine.batch_monthly_interest.id}"
    }
  ]
}
DOC
}

resource "aws_cloudwatch_event_rule" "batch_monthly_interest" {
  name                = "${module.naming.aws_cloudwatch_event_rule}-batch-monthly-interest"
  description         = "run batch account, batch transaction, and batch interest 1st of every month"
  schedule_expression = "${var.batch_monthly_interest_cw_event_scheduled_expression}"
}

resource "aws_cloudwatch_event_target" "batch_monthly_interest_ecs_scheduled_task" {
  arn      = "${aws_sfn_state_machine.batch_monthly_interest.id}"
  rule     = "${aws_cloudwatch_event_rule.batch_monthly_interest.name}"
  role_arn = "${aws_iam_role.batch_monthly_interest.arn}"
}
