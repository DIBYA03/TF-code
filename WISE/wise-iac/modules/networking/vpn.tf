resource "aws_eip" "vpn_eip" {
  count = "${var.enable_vpn == "true" ? 1 : 0}"
  vpc   = true

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-vpn-eip"
    Team        = "${var.team}"
  }
}

resource "aws_ebs_volume" "vpn" {
  count             = "${var.enable_vpn == "true" ? 1 : 0}"
  availability_zone = "${element(var.availability_zones, 1)}"
  size              = 2
  encrypted         = true
  kms_key_id        = "${aws_kms_key.default.arn}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-wise-vpn"
    Team        = "${var.team}"
  }
}

data "template_file" "vpn_userdata" {
  count    = "${var.enable_vpn == "true" ? 1 : 0}"
  template = "${file("./configs/vpn-cloudconfig.yaml")}"

  vars {
    aws_region        = "${var.aws_region}"
    ebs_volume_id     = "${aws_ebs_volume.vpn.id}"
    eip_allocation_id = "${aws_eip.vpn_eip.id}"
  }
}

resource "aws_launch_template" "vpn" {
  count         = "${var.enable_vpn == "true" ? 1 : 0}"
  name_prefix   = "${var.environment}-vpn-"
  image_id      = "${var.vpn_ami}"
  instance_type = "${var.vpn_instance_type}"
  key_name      = "${aws_key_pair.default.key_name}"

  user_data = "${base64encode(data.template_file.vpn_userdata.rendered)}"

  # protect the ec2 instance from being removed
  disable_api_termination = true

  iam_instance_profile = {
    name = "${aws_iam_instance_profile.vpn.name}"
  }

  vpc_security_group_ids = ["${aws_security_group.vpn.id}"]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-wise-vpn"
    Team        = "${var.team}"
  }

  tag_specifications {
    resource_type = "instance"

    tags {
      Application = "${var.application}"
      Component   = "${var.component}"
      Environment = "${var.environment}"
      Name        = "${var.environment}-wise-vpn"
      Team        = "${var.team}"
    }
  }

  tag_specifications {
    resource_type = "volume"

    tags {
      Application = "${var.application}"
      Component   = "${var.component}"
      Environment = "${var.environment}"
      Name        = "${var.environment}-wise-vpn"
      Team        = "${var.team}"
    }
  }
}

resource "aws_autoscaling_group" "vpn" {
  count       = "${var.enable_vpn == "true" ? 1 : 0}"
  name_prefix = "${var.environment}-vpn-"

  min_size                  = 1
  max_size                  = 1
  desired_capacity          = 1
  health_check_type         = "EC2"
  health_check_grace_period = "300"

  vpc_zone_identifier = ["${element(aws_subnet.public_subnets.*.id, 1)}"]

  launch_template = {
    id      = "${aws_launch_template.vpn.id}"
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
      value               = "${var.environment}-wise-vpn"
      propagate_at_launch = true
    },
    {
      key                 = "Team"
      value               = "${var.team}"
      propagate_at_launch = true
    },
  ]
}
