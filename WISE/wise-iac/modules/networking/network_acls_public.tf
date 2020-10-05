resource "aws_network_acl" "public" {
  vpc_id = "${aws_vpc.main.id}"

  subnet_ids = [
    "${aws_subnet.public_subnets.*.id}",
  ]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.public.acl"
    Team        = "${var.team}"
  }
}

# Allow IPV4 http traffic into public subnet
resource "aws_network_acl_rule" "public_ipv4_http_ingress_allow" {
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "100"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 80
  to_port        = 80
}

# Allow IPV6 http traffic into public subnet
resource "aws_network_acl_rule" "public_ipv6_http_ingress_allow" {
  network_acl_id  = "${aws_network_acl.public.id}"
  rule_number     = "200"
  egress          = false
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 80
  to_port         = 80
}

# Allow IPV4 https traffic into public subnet
resource "aws_network_acl_rule" "public_ipv4_https_ingress_allow" {
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "300"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 443
  to_port        = 443
}

# Allow IPV6 https traffic into public subnet
resource "aws_network_acl_rule" "public_ipv6_https_ingress_allow" {
  network_acl_id  = "${aws_network_acl.public.id}"
  rule_number     = "400"
  egress          = false
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 443
  to_port         = 443
}

# Allow return ephemeral traffic from public subnet
resource "aws_network_acl_rule" "ephemeral_to_public_ipv4_ingress_allow" {
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "500"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 1024
  to_port        = 65535
}

# Allow return ephemeral traffic from public subnet
resource "aws_network_acl_rule" "ephemeral_to_public_ipv6_ingress_allow" {
  network_acl_id  = "${aws_network_acl.public.id}"
  rule_number     = "600"
  egress          = false
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 1024
  to_port         = 65535
}

# Public custom NACL rules from tfvars file
resource "aws_network_acl_rule" "public_custom_ingress_nacls" {
  count          = "${length(var.public_subnet_vpc_peering_ingress_nacl_rules)}"
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "${700 + (count.index * var.custom_nacl_multiplier)}"
  egress         = false
  protocol       = "${lookup(var.public_subnet_vpc_peering_ingress_nacl_rules[count.index], "protocol")}"
  rule_action    = "allow"
  cidr_block     = "${lookup(var.public_subnet_vpc_peering_ingress_nacl_rules[count.index], "cidr_block")}"
  from_port      = "${lookup(var.public_subnet_vpc_peering_ingress_nacl_rules[count.index], "from_port")}"
  to_port        = "${lookup(var.public_subnet_vpc_peering_ingress_nacl_rules[count.index], "to_port")}"
}

# Allow outbound traffic from public subnet to everywhere
# Since NAT gateways are in this subnet tier, it's needed
resource "aws_network_acl_rule" "public_ipv4_http_egress_allow" {
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "100"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 80
  to_port        = 80
}

# Allow outbound traffic from public subnet to everywhere
# Since NAT gateways are in this subnet tier, it's needed
resource "aws_network_acl_rule" "public_ipv6_http_egress_allow" {
  network_acl_id  = "${aws_network_acl.public.id}"
  rule_number     = "200"
  egress          = true
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 80
  to_port         = 80
}

# Allow outbound traffic from public subnet to everywhere
# Since NAT gateways are in this subnet tier, it's needed
resource "aws_network_acl_rule" "public_ipv4_https_egress_allow" {
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "300"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 443
  to_port        = 443
}

# Allow outbound traffic from public subnet to everywhere
# Since NAT gateways are in this subnet tier, it's needed
resource "aws_network_acl_rule" "public_ipv6_https_egress_allow" {
  network_acl_id  = "${aws_network_acl.public.id}"
  rule_number     = "400"
  egress          = true
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 443
  to_port         = 443
}

# Allow the public subnet to SSH to app tier
resource "aws_network_acl_rule" "public_to_app_ssh_egress_allow" {
  count          = "${length(var.app_subnet_cidr_blocks)}"
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "${500 + (count.index * var.custom_nacl_multiplier)}"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${element(var.app_subnet_cidr_blocks, count.index)}"
  from_port      = 22
  to_port        = 22
}

# Allow ephemeral outbound traffic from public subnet
resource "aws_network_acl_rule" "ephemeral_ipv4_public_egress_allow" {
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "600"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 1024
  to_port        = 65535
}

# Allow ephemeral outbound traffic from public subnet
resource "aws_network_acl_rule" "ephemeral_ipv6_public_egress_allow" {
  network_acl_id  = "${aws_network_acl.public.id}"
  rule_number     = "700"
  egress          = true
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 1024
  to_port         = 65535
}

# Public custom NACL rules from tfvars file
resource "aws_network_acl_rule" "public_custom_egress_nacls" {
  count          = "${length(var.public_subnet_vpc_peering_egress_nacl_rules)}"
  network_acl_id = "${aws_network_acl.public.id}"
  rule_number    = "${800 + (count.index * var.custom_nacl_multiplier)}"
  egress         = true
  protocol       = "${lookup(var.public_subnet_vpc_peering_egress_nacl_rules[count.index], "protocol")}"
  rule_action    = "allow"
  cidr_block     = "${lookup(var.public_subnet_vpc_peering_egress_nacl_rules[count.index], "cidr_block")}"
  from_port      = "${lookup(var.public_subnet_vpc_peering_egress_nacl_rules[count.index], "from_port")}"
  to_port        = "${lookup(var.public_subnet_vpc_peering_egress_nacl_rules[count.index], "to_port")}"
}
