resource "aws_vpc_endpoint" "dynamo_db" {
  vpc_id            = "${aws_vpc.main.id}"
  service_name      = "com.amazonaws.${var.aws_region}.dynamodb"
  vpc_endpoint_type = "Gateway"

  policy = <<POLICY
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Action": "*",
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "*"
    }
  ]
}
POLICY

  route_table_ids = [
    "${aws_route_table.public.id}",
    "${aws_route_table.app.*.id}",
    "${aws_route_table.db.*.id}",
  ]
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id            = "${aws_vpc.main.id}"
  service_name      = "com.amazonaws.${var.aws_region}.s3"
  vpc_endpoint_type = "Gateway"

  policy = <<POLICY
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Action": "*",
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "*"
    }
  ]
}
POLICY

  route_table_ids = [
    "${aws_route_table.public.id}",
    "${aws_route_table.app.*.id}",
    "${aws_route_table.db.*.id}",
  ]
}

resource "aws_vpc_endpoint" "cloudtrail" {
  vpc_id              = "${aws_vpc.main.id}"
  service_name        = "com.amazonaws.${var.aws_region}.cloudtrail"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.vpc_endpoints.id}"]
}

resource "aws_vpc_endpoint" "ec2" {
  vpc_id              = "${aws_vpc.main.id}"
  service_name        = "com.amazonaws.${var.aws_region}.ec2"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.vpc_endpoints.id}"]
}

resource "aws_vpc_endpoint" "ecr_api" {
  vpc_id              = "${aws_vpc.main.id}"
  service_name        = "com.amazonaws.${var.aws_region}.ecr.api"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.vpc_endpoints.id}"]
}

resource "aws_vpc_endpoint" "ecr_dkr" {
  vpc_id              = "${aws_vpc.main.id}"
  service_name        = "com.amazonaws.${var.aws_region}.ecr.dkr"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.vpc_endpoints.id}"]
}

resource "aws_vpc_endpoint" "execute-api" {
  vpc_id              = "${aws_vpc.main.id}"
  service_name        = "com.amazonaws.${var.aws_region}.execute-api"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.vpc_endpoints.id}"]
}

resource "aws_vpc_endpoint" "kms" {
  vpc_id              = "${aws_vpc.main.id}"
  service_name        = "com.amazonaws.${var.aws_region}.kms"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.vpc_endpoints.id}"]
}

resource "aws_vpc_endpoint" "logs" {
  vpc_id              = "${aws_vpc.main.id}"
  service_name        = "com.amazonaws.${var.aws_region}.logs"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.vpc_endpoints.id}"]
}

resource "aws_vpc_endpoint" "sns" {
  vpc_id              = "${aws_vpc.main.id}"
  service_name        = "com.amazonaws.${var.aws_region}.sns"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true

  subnet_ids         = ["${aws_subnet.app_subnets.*.id}"]
  security_group_ids = ["${aws_security_group.vpc_endpoints.id}"]
}
