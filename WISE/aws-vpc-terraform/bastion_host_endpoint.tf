# Creates the endpoint service, so other accounts can connect to it
resource "aws_vpc_endpoint_service" "bastion_host" {
  count                      = "${var.enable_bastion_host ? 1 : 0}"
  acceptance_required        = true
  network_load_balancer_arns = ["${aws_lb.bastion_host.arn}"]

  allowed_principals = "${var.allowed_principals}"

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${var.environment}-bastion-host"
    Team        = "${var.team}"
  }
}

resource "aws_security_group" "bastion_host_endpoint" {
  count       = "${length(var.bastion_host_vpc_endpoint_service_list) >= 1 ? 1 : 0}"
  name        = "${module.naming.aws_security_group}-bastion-host-vpc-endpoint"
  description = "SG for the VPC endpoints in the ${var.environment} VPC"
  vpc_id      = "${aws_vpc.main.id}"

  ingress {
    description = "allow access to bastion host port"
    from_port   = "${var.bastion_host_port}"
    to_port     = "${var.bastion_host_port}"
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  ingress {
    description = "allow port for http based services"
    from_port   = "80"
    to_port     = "80"
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  ingress {
    description = "allow port for https based services"
    from_port   = "443"
    to_port     = "443"
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment}"
    Name        = "${module.naming.aws_security_group}-bastion-host-vpc-endpoint"
    Team        = "${var.team}"
  }
}

# Creates the connection to the above endpoint service
# Once this is deployed, you will need to access the account and manually accept.
# Doing this for security reasons and a extra step to access the account over ssh
resource "aws_vpc_endpoint" "bastion_host" {
  count = "${length(var.bastion_host_vpc_endpoint_service_list)}"

  vpc_id            = "${aws_vpc.main.id}"
  service_name      = "${lookup(var.bastion_host_vpc_endpoint_service_list[count.index], "service")}"
  vpc_endpoint_type = "Interface"

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.bastion_host_endpoint.id}"]
}
