data "pagerduty_vendor" "cloudwatch" {
  name = "Amazon Cloudwatch"
}

resource "pagerduty_service_integration" "non_critical" {
  name    = "${data.pagerduty_vendor.cloudwatch.name}"
  service = "${var.pagerduty_non_critical_service_id}"
  vendor  = "${data.pagerduty_vendor.cloudwatch.id}"
}

resource "aws_sns_topic_subscription" "pagerduty_non_critical" {
  topic_arn              = "${aws_sns_topic.non_critical.arn}"
  protocol               = "https"
  endpoint               = "https://events.pagerduty.com/integration/${pagerduty_service_integration.non_critical.integration_key}/enqueue"
  endpoint_auto_confirms = true
}

resource "pagerduty_service_integration" "critical" {
  name    = "${data.pagerduty_vendor.cloudwatch.name}"
  service = "${var.pagerduty_critical_service_id}"
  vendor  = "${data.pagerduty_vendor.cloudwatch.id}"
}

resource "aws_sns_topic_subscription" "pagerduty_critical" {
  topic_arn              = "${aws_sns_topic.critical.arn}"
  protocol               = "https"
  endpoint               = "https://events.pagerduty.com/integration/${pagerduty_service_integration.critical.integration_key}/enqueue"
  endpoint_auto_confirms = true
}

data "template_file" "pagerduty_slack_extension_config" {
  count = "${var.enable_pagerduty_slack_integration ? 1 : 0}"

  template = <<EOF
{
  "access_token": "${var.pagerduty_slack_access_token}",
  "bot": {
    "bot_user_id": "ULT1BMAB0"
  },
  "channel": "#pagerduty-incidents",
  "incoming_webhook": {
    "channel": "#pagerduty-incidents",
    "channel_id": "CLSKZB1U7",
    "configuration_url": "${var.pagerduty_slack_configuration_url}",
    "url": "${var.pagerduty_slack_url}"
  },
  "notify_types": {
    "acknowledge": true,
    "annotate": true,
    "assignments": true,
    "resolve": true,
    "ok": true
  },
  "ok": true,
  "referer": "https://wiseco.pagerduty.com/extensions",
  "restrict": "pd-users",
  "scope": "identify,bot,incoming-webhook,channels:read,groups:read,im:read,users:read,users:read.email,chat:write:bot,groups:write",
  "team_id": "TF6L04FUJ",
  "team_name": "Wise",
  "urgency": {
    "high": true,
    "low": true
  }
}
EOF
}

data "pagerduty_extension_schema" "webhook" {
  count = "${var.enable_pagerduty_slack_integration ? 1 : 0}"
  name  = "Slack"
}

resource "pagerduty_extension" "slack_non_critical" {
  count             = "${var.enable_pagerduty_slack_integration ? 1 : 0}"
  name              = "${module.naming.pagerduty_extension}-slack-pagerduty-incidents"
  extension_schema  = "${data.pagerduty_extension_schema.webhook.id}"
  extension_objects = ["${var.pagerduty_non_critical_service_id}"]

  config = "${data.template_file.pagerduty_slack_extension_config.rendered}"
}

resource "pagerduty_extension" "slack_critical" {
  count             = "${var.enable_pagerduty_slack_integration ? 1 : 0}"
  name              = "${module.naming.pagerduty_extension}-slack-pagerduty-incidents"
  extension_schema  = "${data.pagerduty_extension_schema.webhook.id}"
  extension_objects = ["${var.pagerduty_critical_service_id}"]

  config = "${data.template_file.pagerduty_slack_extension_config.rendered}"
}
