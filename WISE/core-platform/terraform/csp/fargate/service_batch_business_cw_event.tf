resource "aws_iam_role" "batch_business" {
  name = "${module.naming.aws_iam_role}-batch-business-cw-event"

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

resource "aws_iam_role_policy" "batch_business_run_task_with_any_role" {
  name = "batch_business_run_task_with_any_role"
  role = "${aws_iam_role.batch_business.id}"

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
            "Resource": "${replace(aws_ecs_task_definition.batch_business.arn, "/:\\d+$/", ":*")}"
        }
    ]
}
DOC
}

resource "aws_cloudwatch_event_rule" "batch_business" {
  name                = "${module.naming.aws_cloudwatch_event_rule}-batch-business"
  description         = "run batch business hourly"
  schedule_expression = "${var.batch_business_cw_event_scheduled_expression}"
}

resource "aws_cloudwatch_event_target" "batch_business_ecs_scheduled_task" {
  arn      = "${aws_ecs_cluster.default.arn}"
  rule     = "${aws_cloudwatch_event_rule.batch_business.name}"
  role_arn = "${aws_iam_role.batch_business.arn}"

  ecs_target {
    task_count          = 1
    task_definition_arn = "${aws_ecs_task_definition.batch_business.arn}"
    launch_type         = "FARGATE"

    network_configuration = {
      subnets          = ["${var.app_subnet_ids}"]
      security_groups  = ["${aws_security_group.batch_business_ecs.id}"]
      assign_public_ip = false
    }
  }
}
