resource "aws_network_acl" "app" {
  vpc_id = "${aws_vpc.main.id}"

  subnet_ids = [
    "${aws_subnet.app_subnets.*.id}",
  ]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.app.acl"
    Team        = "${var.team}"
  }
}

# Allow http traffic coming in from vpc
resource "aws_network_acl_rule" "vpc_to_app_http_ingress_allow" {
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "100"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${var.vpc_cidr_block}"
  from_port      = 80
  to_port        = 80
}

# Allow https traffic coming in from vpc
resource "aws_network_acl_rule" "vpc_to_app_https_ingress_allow" {
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "200"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${var.vpc_cidr_block}"
  from_port      = 443
  to_port        = 443
}

# Allow traffic coming in from public subnet over 22
resource "aws_network_acl_rule" "public_to_app_ssh_ingress_allow" {
  count = "${length(var.public_subnet_cidr_blocks)}"

  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "${300 + (count.index * var.custom_nacl_multiplier)}"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${element(var.public_subnet_cidr_blocks, count.index)}"
  from_port      = 22
  to_port        = 22
}

# So traffic can return back to the subnet
resource "aws_network_acl_rule" "all_to_app_ipv4_ephemeral_egress_allow" {
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "400"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 1024
  to_port        = 65535
}

# So traffic can return back to the subnet
resource "aws_network_acl_rule" "all_to_app_ipv6_ephemeral_egress_allow" {
  network_acl_id  = "${aws_network_acl.app.id}"
  rule_number     = "500"
  egress          = false
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 1024
  to_port         = 65535
}

# App custom NACL rules from tfvars file
resource "aws_network_acl_rule" "app_custom_ingress_nacls" {
  count          = "${length(var.app_subnet_vpc_peering_ingress_nacl_rules)}"
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "${600 + (count.index * var.custom_nacl_multiplier)}"
  egress         = false
  protocol       = "${lookup(var.app_subnet_vpc_peering_ingress_nacl_rules[count.index], "protocol")}"
  rule_action    = "allow"
  cidr_block     = "${lookup(var.app_subnet_vpc_peering_ingress_nacl_rules[count.index], "cidr_block")}"
  from_port      = "${lookup(var.app_subnet_vpc_peering_ingress_nacl_rules[count.index], "from_port")}"
  to_port        = "${lookup(var.app_subnet_vpc_peering_ingress_nacl_rules[count.index], "to_port")}"
}

# Allow http outbound traffic
# If this isn't here, some AMIs can't update correctly
resource "aws_network_acl_rule" "app_to_all_ipv4_http_egress_allow" {
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "100"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 80
  to_port        = 80
}

# Allow http outbound traffic
# If this isn't here, some AMIs can't update correctly
resource "aws_network_acl_rule" "app_to_all_ipv6_http_egress_allow" {
  network_acl_id  = "${aws_network_acl.app.id}"
  rule_number     = "200"
  egress          = true
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 80
  to_port         = 80
}

# Allow https outbound traffic
resource "aws_network_acl_rule" "app_to_all_ipv4_https_egress_allow" {
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "300"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 443
  to_port        = 443
}

# Allow https outbound traffic
resource "aws_network_acl_rule" "app_to_all_ipv6_https_egress_allow" {
  network_acl_id  = "${aws_network_acl.app.id}"
  rule_number     = "400"
  egress          = true
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 443
  to_port         = 443
}

# So app tier can use postgresql
resource "aws_network_acl_rule" "app_to_db_postgres_egress_allow" {
  count          = "${length(var.db_subnet_cidr_blocks)}"
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "5${count.index}0"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${element(var.db_subnet_cidr_blocks, count.index)}"
  from_port      = 5432
  to_port        = 5432
}

# So traffic can happen when going out
resource "aws_network_acl_rule" "app_to_all_ipv4_ephemeral_egress_allow" {
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "600"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
  from_port      = 1024
  to_port        = 65535
}

# So traffic can happen when going out
resource "aws_network_acl_rule" "app_to_all_ipv6_ephemeral_egress_allow" {
  network_acl_id  = "${aws_network_acl.app.id}"
  rule_number     = "700"
  egress          = true
  protocol        = "tcp"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
  from_port       = 1024
  to_port         = 65535
}

# So app tier can use redis
resource "aws_network_acl_rule" "app_to_db_redis_egress_allow" {
  count          = "${length(var.db_subnet_cidr_blocks)}"
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "${800 + (count.index * var.custom_nacl_multiplier)}"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${element(var.db_subnet_cidr_blocks, count.index)}"
  from_port      = 6379
  to_port        = 6379
}

# App custom NACL rules from tfvars file
resource "aws_network_acl_rule" "app_custom_egress_nacls" {
  count          = "${length(var.app_subnet_vpc_peering_egress_nacl_rules)}"
  network_acl_id = "${aws_network_acl.app.id}"
  rule_number    = "${900 + (count.index * var.custom_nacl_multiplier)}"
  egress         = true
  protocol       = "${lookup(var.app_subnet_vpc_peering_egress_nacl_rules[count.index], "protocol")}"
  rule_action    = "allow"
  cidr_block     = "${lookup(var.app_subnet_vpc_peering_egress_nacl_rules[count.index], "cidr_block")}"
  from_port      = "${lookup(var.app_subnet_vpc_peering_egress_nacl_rules[count.index], "from_port")}"
  to_port        = "${lookup(var.app_subnet_vpc_peering_egress_nacl_rules[count.index], "to_port")}"
}
