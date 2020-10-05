aws_profile = "wise-dev"

aws_region = "us-west-2"

application = "rds"

component = "core"

team = "cloud-ops"

vpc_id = "vpc-02d22afaa5a6a4d8a"

vpc_cidr_block = "10.24.0.0/16"

shared_vpc_cidr_block = "10.2.0.0/16"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:152334605517:wise-us-dev1-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:152334605517:wise-us-dev1-noncritical-sns"

# route53
route53_zone_id = "Z053620016VI135DJHSAI"

route53_master_domain = "master.db.dev1.us-west-2.internal.wise.us"

route53_read_replica_domain = "read.db.dev1.us-west-2.internal.wise.us"

app_subnet_ids = [
  "subnet-04b1df6771405bab4",
  "subnet-0faedb9d63ded80a0",
  "subnet-02d9938e32c97b20a",
]

# RDS Subnet Group
db_subnet_group_ids = [
  "subnet-007e4605f149a1d04",
  "subnet-0b2634646b5fb4a95",
  "subnet-0339682fdb486f634",
]

# RDS
rds_instance_class = "db.t2.small"

rds_instance_class_mem = 2 # GiB

rds_storage_size = 20

rds_storage_type = "gp2"

rds_engine = "postgres"

rds_engine_version = "11.5"

rds_parameter_group_family_name = "postgres11"

rds_read_replica_count = 1

rds_backup_retention_period = 1

rds_username = "wiseadmin"

# rds_password = ""
rds_multi_az = false

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
