#!/bin/bash

read --silent --prompt-str 'Enter Hostname: ' rds_hostname
read --silent --prompt-str 'Enter Username: ' rds_username
read --silent --prompt-str 'Enter Password: ' rds_password
read --silent --prompt-str 'Enter Database Name: ' rds_db_name

set -gx RDS_HOSTNAME $rds_hostname
set -gx RDS_USERNAME $rds_username
set -gx RDS_PASSWORD $rds_password
set -gx RDS_DB_NAME $rds_db_name

echo
