aws_profile = "wiseus"

aws_region = "us-west-2"

application = "rds"

component = "core"

team = "cloud-ops"

vpc_id = "vpc-043a00e2bba3b0dd8"

vpc_cidr_block = "10.3.0.0/16"

shared_vpc_cidr_block = "10.2.0.0/16"

csp_vpc_cidr_block = "10.4.0.0/16"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-non-critical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:379379777492:wise-us-staging-noncritical-sns"

# route53
route53_zone_id = "Z1TXQSQLQVDHXJ"

route53_master_domain = "master.db.staging.us-west-2.internal.wise.us."

route53_read_replica_domain = "read.db.staging.us-west-2.internal.wise.us"

app_subnet_ids = [
  "subnet-0e932567d087fdd68",
  "subnet-0b3b84440068c1bf7",
  "subnet-0322b5e9cf73e7eb0",
]

# RDS Subnet Group
db_subnet_group_ids = [
  "subnet-04b619fc8f9ae4c56",
  "subnet-082e9576c0203d034",
  "subnet-0d18b26f6d2de7fa6",
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
database_name = "wise_us_staging" # This is the database name, not RDS

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
