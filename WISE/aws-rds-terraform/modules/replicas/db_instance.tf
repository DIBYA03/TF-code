resource "aws_db_instance" "default" {
  count = "${var.rds_instance_count}"

  # Identity
  identifier          = "${module.naming.aws_db_instance}-${count.index}"
  name                = "${var.rds_database_name}"
  deletion_protection = "${var.rds_deletion_protection}"

  # CA Certificate
  ca_cert_identifier = "${var.rds_ca_cert_identifier}"

  # Replication
  replicate_source_db = "${var.rds_replicate_source_db}"
  kms_key_id          = "${var.rds_kms_key_arn}"

  # Instance info
  instance_class    = "${var.rds_instance_class}"
  allocated_storage = "${var.rds_storage_size}"
  storage_type      = "${var.rds_storage_type}"
  storage_encrypted = "${var.rds_storage_encrypted}"

  # Engine
  engine               = "${var.rds_engine}"
  engine_version       = "${var.rds_engine_version}"
  parameter_group_name = "${aws_db_parameter_group.default.id}"

  # Authentication
  iam_database_authentication_enabled = "${var.iam_database_authentication_enabled}"

  skip_final_snapshot       = "${var.rds_skip_final_snapshot}"
  final_snapshot_identifier = "${var.rds_database_name}"

  # Networking
  db_subnet_group_name   = "${aws_db_subnet_group.default.id}"
  multi_az               = "${var.rds_multi_az}"
  vpc_security_group_ids = ["${aws_security_group.default.id}"]
  publicly_accessible    = false

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${module.naming.aws_db_instance}-${count.index}"
    Team        = "${var.team}"
  }

  provider = "aws.${var.provider_name}"
}
