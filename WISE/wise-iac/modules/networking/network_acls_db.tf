resource "aws_network_acl" "db" {
  vpc_id = "${aws_vpc.main.id}"

  subnet_ids = [
    "${aws_subnet.db_subnets.*.id}",
  ]

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.application}.${var.environment}.db.acl"
    Team        = "${var.team}"
  }
}

# Deny traffic from public subnet tier
# Needed to limit rule count and allow from the rest of the VPC
resource "aws_network_acl_rule" "public_to_db_all_ingress_deny" {
  count          = "${length(var.public_subnet_cidr_blocks)}"
  network_acl_id = "${aws_network_acl.db.id}"
  rule_number    = "${100 + (count.index * var.custom_nacl_multiplier)}"
  egress         = false
  protocol       = "-1"
  rule_action    = "deny"
  cidr_block     = "${element(var.public_subnet_cidr_blocks, count.index)}"
  from_port      = 0
  to_port        = 0
}

# Allows traffic over postgresql from the rest of VPC
# the dbs need to talk to each other and app tier needs access
resource "aws_network_acl_rule" "vpc_to_db_postgresql_ingress_allow" {
  network_acl_id = "${aws_network_acl.db.id}"
  rule_number    = "200"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${var.vpc_cidr_block}"
  from_port      = 5432
  to_port        = 5432
}

# Allows traffic over redis from the rest of VPC
# the dbs need to talk to each other and app tier needs access
resource "aws_network_acl_rule" "vpc_to_db_redis_ingress_allow" {
  network_acl_id = "${aws_network_acl.db.id}"
  rule_number    = "300"
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${var.vpc_cidr_block}"
  from_port      = 6379
  to_port        = 6379
}

# DB custom NACL rules from tfvars file
resource "aws_network_acl_rule" "db_custom_ingress_nacls" {
  count          = "${length(var.db_subnet_vpc_peering_ingress_nacl_rules)}"
  network_acl_id = "${aws_network_acl.db.id}"
  rule_number    = "${400 + (count.index * var.custom_nacl_multiplier)}"
  egress         = false
  protocol       = "${lookup(var.db_subnet_vpc_peering_ingress_nacl_rules[count.index], "protocol")}"
  rule_action    = "allow"
  cidr_block     = "${lookup(var.db_subnet_vpc_peering_ingress_nacl_rules[count.index], "cidr_block")}"
  from_port      = "${lookup(var.db_subnet_vpc_peering_ingress_nacl_rules[count.index], "from_port")}"
  to_port        = "${lookup(var.db_subnet_vpc_peering_ingress_nacl_rules[count.index], "to_port")}"
}

# Return ephemeral ports to VPC
resource "aws_network_acl_rule" "db_to_vpc_ephemeral_egress" {
  network_acl_id = "${aws_network_acl.db.id}"
  rule_number    = "100"
  egress         = true
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = "${var.vpc_cidr_block}"
  from_port      = 1024
  to_port        = 65535
}

# DB custom NACL rules from tfvars file
resource "aws_network_acl_rule" "db_custom_egress_nacls" {
  count          = "${length(var.db_subnet_vpc_peering_egress_nacl_rules)}"
  network_acl_id = "${aws_network_acl.db.id}"
  rule_number    = "${200 + (count.index * var.custom_nacl_multiplier)}"
  egress         = true
  protocol       = "${lookup(var.db_subnet_vpc_peering_egress_nacl_rules[count.index], "protocol")}"
  rule_action    = "allow"
  cidr_block     = "${lookup(var.db_subnet_vpc_peering_egress_nacl_rules[count.index], "cidr_block")}"
  from_port      = "${lookup(var.db_subnet_vpc_peering_egress_nacl_rules[count.index], "from_port")}"
  to_port        = "${lookup(var.db_subnet_vpc_peering_egress_nacl_rules[count.index], "to_port")}"
}
