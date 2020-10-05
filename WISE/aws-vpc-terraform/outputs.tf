output "vpc_id" {
  value = "${aws_vpc.main.id}"
}

output "nat_gateway_ips" {
  value = "${aws_eip.nat_eip.*.public_ip}"
}

output "public_subnets" {
  value = "${aws_subnet.public_subnets.*.id}"
}

output "app_subnets" {
  value = "${aws_subnet.app_subnets.*.id}"
}

output "db_subnets" {
  value = "${aws_subnet.db_subnets.*.id}"
}

output "kms_key_id" {
  value = "${aws_kms_key.default.id}"
}

output "route53_zone_id" {
  value = "${aws_route53_zone.private_default.zone_id}"
}

output "vpn_public_ip" {
  value = "${join("", aws_eip.vpn_eip.*.public_ip)}"
}

# Bastion Host Service
output "bastion_host_private_hostname" {
  value = "${join("", aws_vpc_endpoint_service.bastion_host.*.private_dns_name)}"
}

output "bastion_host_service_names" {
  value = "${join("", aws_vpc_endpoint_service.bastion_host.*.service_name)}"
}

# Bastion Host VPC Endpoints
output "bastion_host_vpc_endpoint_dns_names" {
  value = "${aws_vpc_endpoint.bastion_host.*.dns_entry}"
}
