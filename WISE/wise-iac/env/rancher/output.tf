output "rancher-sg-id" {
  description = "The ID of the security group"
  value       = module.rancher-sg.this_security_group_id
}
output "rancher-private-ip" {
  value  = module.ec2_cluster.private_ip
}

output "rancher-instance-id" {
  description = "instance-id"
  value  = module.ec2_cluster.id
}

output "arn" {
  description = "The ARN of IAM Role"
  value       = module.ec2-iam-role.arn
}
output "profile_name" {
  description = "The Instance profile Name"
  value       = module.ec2-iam-role.profile_name
}



output "rancher-public-ip" {
  description = "Elastic IP of Rancher Instance"
  value       = aws_eip.rancher-elastic-ip.public_ip
}
