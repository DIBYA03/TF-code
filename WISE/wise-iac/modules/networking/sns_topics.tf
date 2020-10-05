resource "aws_sns_topic" "non_critical" {
  name = "${var.application}-${var.environment}-noncritical-sns"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}-${var.environment}-noncritical-sns"
    Team        = "${var.team}"
  }
}

resource "aws_sns_topic" "critical" {
  name = "${var.application}-${var.environment}-critical-sns"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}-${var.environment}-critical-sns"
    Team        = "${var.team}"
  }
}
