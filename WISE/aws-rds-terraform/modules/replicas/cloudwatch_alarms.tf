# https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/MonitoringOverview.html#rds-metrics

# CPU
resource "aws_cloudwatch_metric_alarm" "replica_cpu_non_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-cpu-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS CPU utilization on replica"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_cpu_non_critical_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_cpu_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-cpu-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS CPU utilization on replica"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_cpu_critical_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

# DB Connections
resource "aws_cloudwatch_metric_alarm" "replica_db_connections_non_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-db-conns-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS DB connections on replica"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_conn_count_non_critical}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_db_connections_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-db-conns-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS DB connections on replica"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_conn_count_critical}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

# Freeable memory
resource "aws_cloudwatch_metric_alarm" "replica_free_mem_non_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-free-mem-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free memory on replica"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeableMemory"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_free_mem_non_critical}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_free_mem_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-free-mem-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free memory on replica"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeableMemory"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_free_mem_critical}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_free_disk_non_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-free-disk-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free disk on replica"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeStorageSpace"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_free_disk_non_critical}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_free_disk_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-free-disk-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free disk on replica"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeStorageSpace"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_free_disk_critical}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_max_transaction_ids_non_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-max-transaction-ids-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS max transaction ids on replica"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "MaximumUsedTransactionIDs"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "500000000"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_max_transaction_ids_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-max-transaction-ids-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS max transaction ids on replica"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "MaximumUsedTransactionIDs"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "1000000000"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_replica_lag_non_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-replica-lag-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS replica lag on replica"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "ReplicaLag"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "300"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.non_critical_sns_topic}",
  ]

  ok_actions = [
    "${var.non_critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}

resource "aws_cloudwatch_metric_alarm" "replica_replica_lag_critical" {
  count             = "${var.rds_instance_count}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-replica-${count.index}-replica-lag-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS replica lag on replica"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "ReplicaLag"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "600"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.default.*.id[count.index]}"
  }

  alarm_actions = [
    "${var.critical_sns_topic}",
  ]

  ok_actions = [
    "${var.critical_sns_topic}",
  ]

  insufficient_data_actions = []

  provider = "aws.${var.provider_name}"
}
