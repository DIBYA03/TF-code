provider "aws" {
  region                  = "us-west-2"
  shared_credentials_file = "~/.aws/credentials"
  profile                 = "sbx"
}

terraform {
  backend "s3" {
    bucket         = "terraform.sbx-k8s-wise-us.states"
    key            = "us-west-2/cloudops/vpc-k8s.tfstate"
    region         = "us-west-2"
    dynamodb_table = "sbx-k8s-terraform.state.lock"
    shared_credentials_file = "~/.aws/credentials"
    profile                 = "sbx"
  }
}

locals {
  env = terraform.workspace
  hostedzone = "sbx-k8s.us-west-2.internal.wise.us."
  tags = {
    Name        = "vpc"
    environment = terraform.workspace
    service     = "vpc"
    terraform   = "true"
    repo        = "https://github.com/wiseco/wise-iac"
  }
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${local.env}-vpc"
  cidr = "10.5.0.0/16"
  single_nat_gateway = true
  azs             = ["us-west-2a", "us-west-2b", "us-west-2c"]
  private_subnets = ["10.5.0.0/20", "10.5.16.0/20", "10.5.32.0/19"]
  public_subnets  = ["10.5.64.0/19", "10.5.96.0/19", "10.5.128.0/19"]
  database_subnets  = ["10.5.160.0/19", "10.5.192.0/19", "10.5.224.0/19"]

  enable_nat_gateway = true
  enable_vpn_gateway = true
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Terraform = "true"
    Environment = "sbx"
    "kubernetes.io/cluster/c-7h672" = "owned"
  }
}


resource "aws_route53_zone" "private" {
  name = "k8s-internal.wise.us"
  depends_on = [module.vpc]
  vpc {
    vpc_id = module.vpc.vpc_id
  }
}
