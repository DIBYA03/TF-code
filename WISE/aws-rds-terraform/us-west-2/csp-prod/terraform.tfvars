aws_profile = "wise-prod"

aws_region = "us-west-2"

application = "rds"

component = "csp"

team = "cloud-ops"

vpc_id = "vpc-0a3d5ba3cf7256441"

vpc_cidr_block = "10.20.0.0/16"

peered_vpc_cidr_block = "10.22.0.0/16"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-csp-prod-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-csp-prod-critical-sns"

rds_cross_region_non_critical_sns_topic = "arn:aws:sns:us-east-1:058450407364:wise-us-csp-prod-noncritical-sns"

rds_cross_region_critical_sns_topic = "arn:aws:sns:us-east-1:058450407364:wise-us-csp-prod-critical-sns"

# route53
route53_zone_id = "Z11QB3IQQ5W9E5"

route53_master_domain = "master.db.csp-prod.us-west-2.internal.wise.us"

route53_read_replica_domain = "read.db.csp-prod.us-west-2.internal.wise.us"

app_subnet_ids = [
  "subnet-0ecaac978d90a9d35",
  "subnet-0746233f88b758a7d",
  "subnet-0c4a66fe32ca4033f",
]

# RDS Subnet Group
db_subnet_group_ids = [
  "subnet-0da2337fce43834f2",
  "subnet-058ef0bbb09d9a416",
  "subnet-06dfd715097618f1b",
]

# RDS
rds_instance_class = "db.t2.medium"

rds_instance_class_mem = 4 # GiB

rds_storage_size = 40

rds_storage_type = "gp2"

rds_engine = "postgres"

rds_engine_version = "11.4"

rds_parameter_group_family_name = "postgres11"

rds_read_replica_count = 1

rds_backup_retention_period = 1

rds_username = "wiseadmin"

# rds_password = ""
rds_multi_az = true

# Database
database_name = "wiseus" # This is the database name, not RDS

# Cross-region replication
rds_cross_region_kms_key_arn = "arn:aws:kms:us-east-1:058450407364:key/1510fd4d-a175-4713-bcea-f29bcfd9da2b"

rds_cross_region_instance_count = "1"

rds_cross_region_multi_az = true

rds_cross_region_vpc_id = "vpc-0b0012e20fb5bc05e"

rds_cross_region_vpc_cidr_block = "10.22.0.0/16"

rds_cross_region_db_subnet_group_ids = [
  "subnet-04a28a3d49bad1eb2",
  "subnet-0a0801396fbcf61e0",
  "subnet-0a9f78e94a8d8171b",
]

# CloudWatch Alarms
rds_cw_conn_count_non_critical = 40 # above percent

rds_cw_conn_count_critical = 50 # above percent

rds_cw_cpu_non_critical_limit = 75 # above percent

rds_cw_cpu_critical_limit = 85 # above percent

rds_cw_free_mem_non_critical = 40 # below percent

rds_cw_free_mem_critical = 30 # below percent

rds_cw_free_disk_non_critical = 50 # below percent

rds_cw_free_disk_critical = 25 # below percent

# Backups
backup_lambda_timeout = 30

rds_enable_backup_lambda = true

# cross-region backups support
enable_cross_region_backup = "true"
