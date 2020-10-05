data "template_file" "bbva_notifications_sns_policy" {
  template = <<EOF
{
  "Version": "2008-10-17",
  "Id": "BBVAToWiseNotificationsSNS",
  "Statement": [
    {
      "Sid": "SNSOwnerPermissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": [
        "sns:GetTopicAttributes",
        "sns:SetTopicAttributes",
        "sns:AddPermission",
        "sns:RemovePermission",
        "sns:DeleteTopic",
        "sns:Subscribe",
        "sns:ListSubscriptionsByTopic",
        "sns:Publish",
        "sns:Receive"
      ],
      "Resource": "arn:aws:sns:${var.aws_region}:${data.aws_caller_identity.account.account_id}:${module.naming.aws_sns_topic}",
      "Condition": {
        "StringEquals": {
          "AWS:SourceOwner": "${data.aws_caller_identity.account.account_id}"
        }
      }
    },
    {
      "Sid": "SNSWiseSubscribers",
      "Effect": "Allow",
      "Principal": {
        "AWS": ${jsonencode(var.sns_allowed_subscribe_accounts)}
      },
      "Action": [
        "sns:GetTopicAttributes",
        "sns:Subscribe",
        "sns:Receive"
      ],
      "Resource": "arn:aws:sns:${var.aws_region}:${data.aws_caller_identity.account.account_id}:${module.naming.aws_sns_topic}"
    }
  ]
}
EOF
}

resource "aws_sns_topic" "bbva_notifications" {
  name              = "${module.naming.aws_sns_topic}"
  kms_master_key_id = "${aws_kms_alias.bbva_sqs.target_key_arn}"

  policy = "${data.template_file.bbva_notifications_sns_policy.rendered}"

  delivery_policy = <<EOF
{
  "http": {
    "defaultHealthyRetryPolicy": {
      "minDelayTarget": 5,
      "maxDelayTarget": 10,
      "numRetries": 3,
      "numMaxDelayRetries": 0,
      "numNoDelayRetries": 0,
      "numMinDelayRetries": 3,
      "backoffFunction": "linear"
    }
  }
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_sns_topic}"
    Team        = "${var.team}"
  }
}
