# https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/MonitoringOverview.html#rds-metrics

# CPU
resource "aws_cloudwatch_metric_alarm" "backup_cpu_non_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-cpu-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS CPU utilization on backup"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_cpu_non_critical_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-cpu-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_cpu_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-cpu-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS CPU utilization on backup"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${var.rds_cw_cpu_critical_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-cpu-crit"
    Team        = "${var.team}"
  }
}

# DB Connections
resource "aws_cloudwatch_metric_alarm" "backup_db_connections_non_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-db-conns-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS DB connections on backup"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_non_critical_conn_count_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-db-conns-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_db_connections_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-db-conns-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS DB connections on backup"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_critical_conn_count_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-db-conns-crit"
    Team        = "${var.team}"
  }
}

# Freeable memory
resource "aws_cloudwatch_metric_alarm" "backup_free_mem_non_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-free-mem-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free memory on backup"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeableMemory"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_non_critical_mem_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-free-mem-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_free_mem_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-free-mem-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free memory on backup"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeableMemory"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_critical_mem_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-free-mem-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_free_disk_non_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-free-disk-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free disk on backup"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeStorageSpace"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_non_critical_disk_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-free-disk-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_free_disk_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-free-disk-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS free disk on backup"

  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "FreeStorageSpace"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "${local.rds_critical_disk_limit}"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-free-disk-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_max_transaction_ids_non_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-max-transaction-ids-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS max transaction ids on backup"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "MaximumUsedTransactionIDs"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "500000000"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-max-transaction-ids-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_max_transaction_ids_critical" {
  count = "${var.rds_enable_backup_lambda ? 1 : 0}"

  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-max-transaction-ids-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS max transaction ids on backup"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "MaximumUsedTransactionIDs"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "1000000000"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}--backup-max-transaction-ids-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_replica_lag_non_critical" {
  count             = "${var.rds_enable_backup_lambda ? 1 : 0}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-${count.index}-replica-lag-non-crit"
  alarm_description = "This metric monitors ${title(terraform.workspace)} RDS replica lag on backup"

  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "10"
  metric_name         = "ReplicaLag"
  namespace           = "AWS/RDS"
  period              = "120"
  statistic           = "Average"
  threshold           = "300"
  treat_missing_data  = "breaching"

  dimensions {
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-${count.index}-replica-lag-non-crit"
    Team        = "${var.team}"
  }
}

resource "aws_cloudwatch_metric_alarm" "backup_replica_lag_critical" {
  count             = "${var.rds_enable_backup_lambda ? 1 : 0}"
  alarm_name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-${count.index}-replica-lag-crit"
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
    DBInstanceIdentifier = "${aws_db_instance.backup_replica.id}"
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
    Name        = "${module.naming.aws_cloudwatch_metric_alarm}-backup-${count.index}-replica-lag-crit"
    Team        = "${var.team}"
  }
}
