resource "aws_security_group" "vpc_endpoints" {
  name        = "${module.naming.aws_security_group}-vpc-endpoints"
  description = "SG for the VPC endpoints in the ${var.environment} VPC"
  vpc_id      = "${aws_vpc.main.id}"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-vpc-endpoints"
    Team        = "${var.team}"
  }
}

resource "aws_security_group" "vpn" {
  count = "${var.enable_vpn == "true" ? 1 : 0}"

  name        = "${module.naming.aws_security_group}-vpn"
  description = "SG for the VPN in the ${var.environment} VPC"
  vpc_id      = "${aws_vpc.main.id}"

  # Allow SSH only from the VPC
  ingress {
    description = "${var.environment} ssh access"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  # Allow UI access for all
  # Seems this is needed to connect to the VPN
  # Needs investigation on fixing this
  ingress {
    description = "all https access"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow everyone to be able to access the VPN
  ingress {
    from_port   = 1194
    to_port     = 1194
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow outbound to the VPCs
  egress {
    description = "Outbound access to all possible Wise VPCs"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["10.0.0.0/8"]
  }

  # ubuntu apt servers
  # This is how I\to got the list
  # dig +short $(grep -Pho '^\s*[^#].*?https?://\K[^/]+(?=.*updates)' \
  #             /etc/apt/sources.list /etc/apt/sources.list.d/*.list | sort -u) | sort -u
  egress {
    description = "ubuntu apt server"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"

    cidr_blocks = [
      "34.210.25.51/32",
      "34.212.136.213/32",
      "54.190.18.91/32",
      "54.191.55.41/32",
      "54.191.70.203/32",
      "54.218.137.160/32",
      "91.189.88.162/32",
    ]
  }

  # # OpenVPN licensing servers
  # egress {
  #   description = "openvpn license server"
  #   from_port   = 443
  #   to_port     = 443
  #   protocol    = "tcp"
  #
  #   cidr_blocks = [
  #     "107.191.99.82/32",
  #     "107.161.19.201/32",
  #   ]
  # }
  egress {
    description = "allow outbound traffic for VPN users over port 443"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-vpn"
    Team        = "${var.team}"
  }
}
