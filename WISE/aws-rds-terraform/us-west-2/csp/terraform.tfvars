aws_profile = "wiseus"

aws_region = "us-west-2"

application = "rds"

component = "core"

team = "cloud-ops"

vpc_id = "vpc-0aef975032e418d38"

vpc_cidr_block = "10.4.0.0/16"

shared_vpc_cidr_block = "10.2.0.0/16"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-csp-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-csp-noncritical-sns"

rds_cross_region_non_critical_sns_topic = "arn:aws:sns:us-east-1:379379777492:wise-us-csp-noncritical-sns"

rds_cross_region_critical_sns_topic = "arn:aws:sns:us-east-1:379379777492:wise-us-csp-noncritical-sns"

# route53
route53_zone_id = "Z1LYWCVCY2AD1D"

route53_master_domain = "master.db.csp.internal.wise.us"

route53_read_replica_domain = "read.db.csp.internal.wise.us"

app_subnet_ids = [
  "subnet-0f20d4c37f4242b7b",
  "subnet-051e08c06df7de0f0",
  "subnet-07bd6a1a5586635b0",
]

# RDS Subnet Group
db_subnet_group_ids = [
  "subnet-01f513db6f9693897",
  "subnet-0e987051f1a46e9cc",
  "subnet-03f3c8e82530c0742",
]

# RDS
rds_instance_class = "db.t2.small"

rds_instance_class_mem = 2 # GiB

rds_storage_size = 20

rds_storage_type = "gp2"

rds_engine = "postgres"

rds_engine_version = "11.2"

rds_parameter_group_family_name = "postgres11"

rds_read_replica_count = 1

rds_backup_retention_period = 1

rds_username = "wiseadmin"

# rds_password = ""
rds_multi_az = true

# Database
database_name = "wiseus" # This is the database name, not RDS

enable_cross_region_backup = "false"

# CloudWatch Alarms
rds_cw_conn_count_non_critical = 60 # above percent

rds_cw_conn_count_critical = 70 # above percent

rds_cw_cpu_non_critical_limit = 80 # above percent

rds_cw_cpu_critical_limit = 90 # above percent

rds_cw_free_mem_non_critical = 40 # below percent

rds_cw_free_mem_critical = 30 # below percent

rds_cw_free_disk_non_critical = 20 # below percent

rds_cw_free_disk_critical = 10 # below percent

# Backups
backup_lambda_timeout = 30

rds_enable_backup_lambda = false

# cross-region backups support
enable_cross_region_backup = "false"
