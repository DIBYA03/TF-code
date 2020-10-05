aws_profile = "wise-prod"

aws_region = "us-west-2"

application = "rds"

component = "core"

team = "cloud-ops"

vpc_id = "vpc-0dbc76a72cb520546"

vpc_cidr_block = "10.17.0.0/16"

csp_vpc_cidr_block = "10.20.0.0/16"

peered_vpc_cidr_block = "10.18.0.0/16"

# SNS
non_critical_sns_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-noncritical-sns"

critical_sns_topic = "arn:aws:sns:us-west-2:058450407364:wise-us-prod-critical-sns"

rds_cross_region_non_critical_sns_topic = "arn:aws:sns:us-east-1:058450407364:wise-us-prod-noncritical-sns"

rds_cross_region_critical_sns_topic = "arn:aws:sns:us-east-1:058450407364:wise-us-prod-critical-sns"

# route53
route53_zone_id = "ZEQGF56E9VNY4"

route53_master_domain = "master.db.prod.us-west-2.internal.wise.us"

route53_read_replica_domain = "read.db.prod.us-west-2.internal.wise.us"

app_subnet_ids = [
  "subnet-06cdb0c22471f3c1d",
  "subnet-09259e808681d9c5a",
  "subnet-0d75875254efbc0b1",
]

# RDS Subnet Group
db_subnet_group_ids = [
  "subnet-01d92def7950c7cb0",
  "subnet-08fb45a2c2c40eaaa",
  "subnet-0f3a1b792fd1a272f",
]

# RDS
rds_instance_class = "db.t2.large"

rds_instance_class_mem = 8 # GiB

rds_storage_size = 20

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
rds_cross_region_kms_key_arn = "arn:aws:kms:us-east-1:058450407364:key/81b01b74-f88a-4065-929e-e25a5ddb56c7"

rds_cross_region_instance_count = "1"

rds_cross_region_multi_az = true

rds_cross_region_vpc_id = "vpc-0c3996a3cb4a1e002"

rds_cross_region_vpc_cidr_block = "10.18.0.0/16"

rds_cross_region_db_subnet_group_ids = [
  "subnet-022d6401d35f0c597",
  "subnet-05bb246d757e4d065",
  "subnet-0e388759ec9a1ffc4",
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
