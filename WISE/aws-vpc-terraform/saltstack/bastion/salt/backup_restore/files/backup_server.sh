#!/bin/bash
set -x

# timestamp for directory
export TIMESTAMP=$(date +%s)

# make sure this is mounted first
mount -t nfs -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport {{ efs_dns_name }}:/ /mnt/efs/

# backup the /home folder
mkdir -p /mnt/efs/$TIMESTAMP/home
rsync -a /home/ /mnt/efs/$TIMESTAMP/home

# remove old backups in days
find /mnt/efs -maxdepth 1 -type d -mtime +{{ days_to_backup }} -print -exec rm -rf {} \;

# TODO
# add restore script
