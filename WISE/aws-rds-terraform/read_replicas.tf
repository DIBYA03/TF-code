resource "aws_db_instance" "read_replica" {
  count = "${var.rds_read_replica_count}"

  # Identity
  identifier          = "${module.naming.aws_db_instance}-read-${count.index}"
  name                = "${var.database_name}"
  replicate_source_db = "${aws_db_instance.master.id}"
  deletion_protection = "${var.rds_deletion_protection}"

  # CA Certificate
  ca_cert_identifier = "${var.rds_ca_cert_identifier}"

  # Backups
  maintenance_window        = "${var.rds_replica_maintenance_window}"
  copy_tags_to_snapshot     = true
  final_snapshot_identifier = "${module.naming.aws_db_instance}-read-${count.index}"

  # Instance info
  instance_class    = "${var.rds_instance_class}"
  allocated_storage = "${var.rds_storage_size}"
  storage_type      = "${var.rds_storage_type}"
  storage_encrypted = "${var.rds_storage_encrypted}"
  kms_key_id        = "${aws_kms_key.rds_default.arn}"

  # Engine
  engine               = "${var.rds_engine}"
  engine_version       = "${var.rds_engine_version}"
  parameter_group_name = "${aws_db_parameter_group.default.id}"

  # Authentication
  iam_database_authentication_enabled = "${var.rds_iam_database_authentication_enabled}"

  # Networking
  multi_az               = "${var.rds_multi_az}"
  vpc_security_group_ids = ["${aws_security_group.default.id}"]
  publicly_accessible    = false

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${terraform.workspace}"
    Name        = "${module.naming.aws_db_instance}-read-${count.index}"
    Team        = "${var.team}"
  }
}
