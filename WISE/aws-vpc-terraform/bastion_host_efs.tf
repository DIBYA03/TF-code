resource "aws_efs_file_system" "bastion_host" {
  count          = "${var.enable_bastion_host ? 1 : 0}"
  creation_token = "${var.environment}-bastion-host"

  encrypted  = true
  kms_key_id = "${aws_kms_alias.default.target_key_arn}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-bst-hst"
    Team        = "${var.team}"
  }
}

locals {
  bastion_host_count           = "${var.enable_bastion_host ? 1 : 0}"
  bastion_host_efs_mount_count = "${local.bastion_host_count * length(var.availability_zones)}"
}

resource "aws_efs_mount_target" "bastion_host" {
  count          = "${local.bastion_host_efs_mount_count}"
  file_system_id = "${aws_efs_file_system.bastion_host.id}"

  subnet_id       = "${element(aws_subnet.app_subnets.*.id, count.index)}"
  security_groups = ["${aws_security_group.bastion_host_efs.id}"]
}

resource "aws_security_group" "bastion_host_efs" {
  count       = "${var.enable_bastion_host ? 1 : 0}"
  name        = "${module.naming.aws_security_group}-bst-hst-efs"
  description = "security group for ${var.environment} bastion host efs"
  vpc_id      = "${aws_vpc.main.id}"

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = ["${aws_security_group.bastion_host.id}"]
  }

  ingress {
    from_port       = 2049
    to_port         = 2049
    protocol        = "tcp"
    security_groups = ["${aws_security_group.bastion_host.id}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-bst-hst-efs"
    Team        = "${var.team}"
  }
}
