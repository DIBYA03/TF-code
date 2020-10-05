# https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/MonitoringOverview.html#rds-metrics

# CPU
resource "aws_cloudwatch_metric_alarm" "master_cpu_non_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-cpu-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS CPU utilization on master"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_cpu_non_critical_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-cpu-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "master_cpu_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-cpu-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS CPU utilization on master"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_cpu_critical_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-cpu-crit"
    Team        = "${var.team}"
  }
}

# DB Connections
resource "aws_cloudwatch_metric_alarm" "master_db_connections_non_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-db-conns-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS DB connections on master"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_non_critical_conn_count_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-db-conns-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "master_db_connections_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-db-conns-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS DB connections on master"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_critical_conn_count_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-db-conns-crit"
    Team        = "${var.team}"
  }
}

# Freeable memory
resource "aws_cloudwatch_metric_alarm" "master_free_mem_non_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-free-mem-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free memory on master"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeableMemory"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_non_critical_mem_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-free-mem-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "master_free_mem_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-free-mem-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free memory on master"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeableMemory"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_critical_mem_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-free-mem-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "master_free_disk_non_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-free-disk-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free disk on master"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeStorageSpace"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_non_critical_disk_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-free-disk-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "master_free_disk_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-free-disk-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free disk on master"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeStorageSpace"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_critical_disk_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-master-free-disk-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "master_max_transaction_ids_non_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-max-transaction-ids-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS max transaction ids on master"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "MaximumUsedTransactionIDs"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "500000000"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-max-transaction-ids-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "master_max_transaction_ids_critical" {
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-max-transaction-ids-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS max transaction ids on master"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "MaximumUsedTransactionIDs"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "1000000000"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.master.id}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-max-transaction-ids-crit"
    Team        = "${var.team}"
  }
}
