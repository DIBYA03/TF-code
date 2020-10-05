output "rds_master_hostname" {
  value = "${aws_db_instance.master.endpoint}"
}

output "rds_read_hostname" {
  value = "${aws_db_instance.read_replica.*.endpoint}"
}
