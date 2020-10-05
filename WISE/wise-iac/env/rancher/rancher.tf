provider "aws" {
  region                  = "us-west-2"
  shared_credentials_file = "~/.aws/credentials"
  profile                 = "sbx"
}



terraform {
  backend "s3" {
    bucket         = "terraform.sbx-k8s-wise-us.states"
    key            = "us-west-2/cloudops/ec2/rancher.tfstate"
    region         = "us-west-2"
    dynamodb_table = "sbx-k8s-terraform.state.lock"
    shared_credentials_file = "~/.aws/credentials"
    profile                 = "sbx"
  }
}

locals {
  env = terraform.workspace
  vpcname = var.vpcname
  subnetname = var.subnetname
  hostedzone = var.hostedzone
  tags = {
    Name	= "Rancher-server"
    environment = terraform.workspace
    app         = "rancher"
    service     = "Ec2"
    terraform   = "true"
    repo        = "https://github.com/wiseco/wise-iac"
  }
}


data "aws_vpc" "vpc" {
  tags= {
    Name = "${local.vpcname}"
  }
}

data "aws_subnet_ids" "subnet-ids" {
  vpc_id = (data.aws_vpc.vpc.id)
  tags = {
  Name = "${local.subnetname}"
  }
}
data "aws_route53_zone" "hosted-zone" {
    name = "${local.hostedzone}"
}

module "rancher-sg" {
  source = "terraform-aws-modules/security-group/aws"
  name        = "rancher-sg"
  description = "Security group for user-service with custom ports open within VPC, and http traffic open for vpc"
  vpc_id      = (data.aws_vpc.vpc.id)
  ingress_cidr_blocks      = ["0.0.0.0/0"]
  ingress_rules            = ["https-443-tcp","http-80-tcp"]
  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules = ["all-all"]
  tags = local.tags
}

resource "aws_security_group_rule" "ingress_rules" {
  cidr_blocks              = [
          data.aws_vpc.vpc.cidr_block, 
	  "10.8.0.0/16"
        ]
      description              = "SSH"
      from_port                = 22
      protocol                 = "tcp"
      security_group_id        = "sg-076a90209e0f1e9c2"
      to_port                  = 22
      type                     = "ingress"
    }

resource "aws_key_pair" "rancher-ssh" {
  key_name   = "rancher-ssh"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDuy2oxmfP9/jIbqYzC8hGRaYTrULNOZVtTIwrWOuguFta3QcWn20e96IAb6Gn5h3uvg6a9uZOy6VikDdr76aPpmjQdNRXHm9U379NTIo5JwDAOj0wdQ2o02s79SrWc39e5ZNGfh0Wv367CPl1dbqIcB9dJ2hy9SqpnZr/lcqmwv8QWju8Aqi2ITDWoyeS7hy8JBNiLuoHwnRzPpCXrcibhZWqmy7BgpvmBrZJ/qMbPpUgyxvmdt6rrYqZbDt5p5lOOCzMCD2ACdOWBQZ41Sh4TOqDQsiDclNhX0dgM/D33Q2iestXvrs1G1T514cGVIbBChvDNfzmpO/Y8h4Sli0vD niting@talentica-all.com@NitinG-ub"
}

module "ec2_cluster" {
  source                 = "../../modules/aws-ec2-instance/"
  name                   = "${local.env}-rancher"
  instance_count         = 1
  ami                    = "ami-003634241a8fcdec0"
  instance_type          = "t3.large"
  iam_instance_profile   = module.ec2-iam-role.profile_name
  private_ip		 = var.privateip
  key_name               = "rancher-ssh"
  monitoring             = true
  vpc_security_group_ids = [module.rancher-sg.this_security_group_id]
#  subnet_id              = "${data.aws_subnet_ids.subnet-ids.id}"
  subnet_id              = var.subnetid
  user_data = "${file("init.sh")}" 
  tags = local.tags

}

resource "aws_eip" "rancher-elastic-ip" {
  vpc = true
  instance  = "${element(module.ec2_cluster.id, 1)}"
  associate_with_private_ip = var.privateip 
}

resource "aws_route53_record" "rancher_internal" {
  zone_id = (data.aws_route53_zone.hosted-zone.zone_id)
  name    = "rancher"
  type    = "A"
  ttl     = "300"
  records = [aws_eip.rancher-elastic-ip.public_ip]
}

