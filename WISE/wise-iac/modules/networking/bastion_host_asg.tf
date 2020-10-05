data "template_file" "bastion_host_userdata" {
  count    = "${var.enable_bastion_host ? 1 : 0}"
  template = "${file("./configs/bastion-cloudconfig.yaml")}"

  vars {
    bastion_hostname          = "${var.environment}-bastion"
    bastion_host_efs_hostname = "${aws_efs_file_system.bastion_host.dns_name}"
    s3_bucket                 = "${aws_s3_bucket.backup_bucket.id}"
    s3_saltstack_prefix       = "${var.bastion_host_salstack_s3_object_prefix}"
    s3_saltstack_filename     = "${var.bastion_host_salstack_s3_object_name}"
  }
}

resource "aws_s3_bucket_object" "bastion_host_saltstack" {
  count  = "${var.enable_bastion_host ? 1 : 0}"
  bucket = "${aws_s3_bucket.backup_bucket.id}"
  key    = "${var.bastion_host_salstack_s3_object_prefix}/${var.bastion_host_salstack_s3_object_name}"
  source = "./saltstack/bastion/saltstack.zip"
  etag   = "${filemd5("./saltstack/bastion/saltstack.zip")}"
}

resource "aws_launch_template" "bastion_host" {
  count         = "${var.enable_bastion_host ? 1 : 0}"
  name_prefix   = "${var.environment}-bst-hst-"
  image_id      = "${data.aws_ami.ubuntu.id}"
  instance_type = "${var.bastion_host_instance_type}"
  key_name      = "${aws_key_pair.default.key_name}"

  user_data = "${base64encode(data.template_file.bastion_host_userdata.rendered)}"

  iam_instance_profile {
    arn = "${aws_iam_instance_profile.bastion_host.arn}"
  }

  # protect the ec2 instance from being removed
  # disable_api_termination = true

  vpc_security_group_ids = ["${aws_security_group.bastion_host.id}"]
  block_device_mappings {
    device_name = "/dev/sda1"

    ebs {
      delete_on_termination = true
      volume_size           = 20

      encrypted  = true
      kms_key_id = "${aws_kms_alias.default.target_key_arn}"
    }
  }
  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-bst-hst"
    Team        = "${var.team}"
  }
  tag_specifications {
    resource_type = "instance"

    tags {
      Application = "${var.application}"
      Component   = "${var.component}"
      Environment = "${var.environment}"
      Name        = "${var.environment}-bst-hst"
      Team        = "${var.team}"
    }
  }
  tag_specifications {
    resource_type = "volume"

    tags {
      Application = "${var.application}"
      Component   = "${var.component}"
      Environment = "${var.environment}"
      Name        = "${var.environment}-bst-hst"
      Team        = "${var.team}"
    }
  }
}

resource "aws_autoscaling_group" "bastion_host" {
  count       = "${var.enable_bastion_host ? 1 : 0}"
  name_prefix = "${var.environment}-bst-hst-"

  min_size                  = "${var.bastion_host_min_size}"
  max_size                  = "${var.bastion_host_max_size}"
  desired_capacity          = "${var.bastion_host_desired_capacity}"
  health_check_type         = "EC2"
  health_check_grace_period = "300"

  vpc_zone_identifier = ["${aws_subnet.app_subnets.*.id}"]

  target_group_arns = [
    "${aws_lb_target_group.bastion_host.arn}",
  ]

  launch_template = {
    id      = "${aws_launch_template.bastion_host.id}"
    version = "$$Latest"
  }

  tags = [
    {
      key                 = "Application"
      value               = "${var.application}"
      propagate_at_launch = true
    },
    {
      key                 = "Component"
      value               = "${var.component}"
      propagate_at_launch = true
    },
    {
      key                 = "Environment"
      value               = "${var.environment}"
      propagate_at_launch = true
    },
    {
      key                 = "Name"
      value               = "${var.environment}-bastion-host"
      propagate_at_launch = true
    },
    {
      key                 = "Team"
      value               = "${var.team}"
      propagate_at_launch = true
    },
  ]
}
