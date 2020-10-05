resource "aws_iam_role" "batch_monthly_interest_step_function" {
  name = "${module.naming.aws_iam_role}-batch-monthly-interest-sfn"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "states.${var.aws_region}.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_iam_role}-batch-monthly-interest-sfn"
    Team        = "${var.team}"
  }
}

resource "aws_iam_role_policy_attachment" "role_AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = "${aws_iam_role.batch_monthly_interest_step_function.name}"
}

resource "aws_iam_role_policy" "batch_monthly_interest_step_function" {
  name = "${module.naming.aws_iam_policy}-default-execution"
  role = "${aws_iam_role.batch_monthly_interest_step_function.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecs:RunTask"
      ],
      "Resource": [
        "${aws_ecs_task_definition.batch_account.arn}",
        "${aws_ecs_task_definition.batch_transaction.arn}",
        "${aws_ecs_task_definition.batch_monthly_interest.arn}",
        "${aws_ecs_task_definition.batch_monitor.arn}"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "ecs:StopTask",
        "ecs:DescribeTasks"
      ],
      "Resource": [
        "*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "events:PutTargets",
        "events:PutRule",
        "events:DescribeRule"
      ],
      "Resource": [
        "arn:aws:events:${var.aws_region}:${data.aws_caller_identity.account.account_id}:rule/StepFunctionsGetEventsForECSTaskRule"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "iam:GetRole",
        "iam:PassRole"
      ],
      "Resource": [
        "${aws_iam_role.batch_account_execution_role.arn}",
        "${aws_iam_role.batch_transaction_execution_role.arn}",
        "${aws_iam_role.batch_monthly_interest_execution_role.arn}",
        "${aws_iam_role.batch_monitor_execution_role.arn}"
      ]
    }
  ]
}
EOF
}

resource "aws_sfn_state_machine" "batch_monthly_interest" {
  name     = "${module.naming.aws_sfn_state_machine}-batch-monthly-interest"
  role_arn = "${aws_iam_role.batch_monthly_interest_step_function.arn}"

  definition = <<EOF
{
  "Comment": "run batch account and then run monthly interest services",
  "StartAt": "Run Batch Account",
  "States": {

    "Run Batch Account": {
      "Type": "Task",
      "Resource": "arn:aws:states:::ecs:runTask.sync",
      "Parameters": {
        "LaunchType": "FARGATE",
        "Cluster": "${aws_ecs_cluster.default.arn}",
        "TaskDefinition": "${aws_ecs_task_definition.batch_account.arn}",
        "NetworkConfiguration": {
          "AwsvpcConfiguration": {
            "Subnets": ${jsonencode(var.app_subnet_ids)},
            "AssignPublicIp": "DISABLED"
          }
        }
      },
      "Next": "Run Batch Transaction"
    },

    "Run Batch Transaction": {
      "Type": "Task",
      "Resource": "arn:aws:states:::ecs:runTask.sync",
      "Parameters": {
        "LaunchType": "FARGATE",
        "Cluster": "${aws_ecs_cluster.default.arn}",
        "TaskDefinition": "${aws_ecs_task_definition.batch_transaction.arn}",
        "NetworkConfiguration": {
          "AwsvpcConfiguration": {
            "Subnets": ${jsonencode(var.app_subnet_ids)},
            "AssignPublicIp": "DISABLED"
          }
        }
      },
      "Next": "Wait 5 minutes"
    },

    "Wait 5 minutes": {
      "Comment": "Wait five minutes for read-replica sync",
      "Type": "Wait",
      "Seconds": 300,
      "Next": "Run Batch Monthly Interest"
    },

    "Run Batch Monthly Interest": {
      "Type": "Task",
      "Resource": "arn:aws:states:::ecs:runTask.sync",
      "Parameters": {
        "LaunchType": "FARGATE",
        "Cluster": "${aws_ecs_cluster.default.arn}",
        "TaskDefinition": "${aws_ecs_task_definition.batch_monthly_interest.arn}",
        "NetworkConfiguration": {
          "AwsvpcConfiguration": {
            "Subnets": ${jsonencode(var.app_subnet_ids)},
            "AssignPublicIp": "DISABLED"
          }
        }
      },
      "Next": "Wait 5 minutes monitor"
    },

    "Wait 5 minutes monitor": {
      "Comment": "Wait five minutes for read-replica sync",
      "Type": "Wait",
      "Seconds": 300,
      "Next": "Run Batch Monitor"
    },

    "Run Batch Monitor": {
      "Type": "Task",
      "Resource": "arn:aws:states:::ecs:runTask.sync",
      "Parameters": {
        "LaunchType": "FARGATE",
        "Cluster": "${aws_ecs_cluster.default.arn}",
        "TaskDefinition": "${aws_ecs_task_definition.batch_monitor.arn}",
        "NetworkConfiguration": {
          "AwsvpcConfiguration": {
            "Subnets": ${jsonencode(var.app_subnet_ids)},
            "AssignPublicIp": "DISABLED"
          }
        }
      },
      "End": true
    }
  }
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_sfn_state_machine}-batch-monthly-interest"
    Team        = "${var.team}"
  }

  depends_on = [
    "aws_iam_role.batch_monthly_interest_step_function",
    "aws_iam_role_policy.batch_monthly_interest_step_function",
  ]
}

resource "aws_cloudwatch_metric_alarm" "batch_monthly_interest_step_function_fail" {
  alarm_name          = "${module.naming.aws_cloudwatch_metric_alarm}-batch-monthly-interest-fail"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "ExecutionsFailed"
  namespace           = "AWS/States"
  period              = "3600"                                                                     # 1 hour
  statistic           = "Sum"
  threshold           = "1"
  treat_missing_data  = "notBreaching"

  dimensions {
    StateMachineArn = "${aws_sfn_state_machine.batch_monthly_interest.id}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]
}
