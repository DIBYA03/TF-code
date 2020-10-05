
terraform {
  backend "s3" {
    bucket         = "terraform.sbx-k8s-wise-us.states"
    key            = "rke-cluster/cluster.tfstate"
//    key            = "us-west-2/cloudops/ec2/rancher.tfstate"
    region         = "us-west-2"
    dynamodb_table = "sbx-k8s-terraform.state.lock"
    shared_credentials_file = "~/.aws/credentials"
    profile                 = "sbx"
//    workspaces {
//      prefix = "rke-cluster"
//    }
  }
}

locals {
  env = terraform.workspace
  tags = {
    Name	= ""
    environment = terraform.workspace
    app         = "rancher"
    service     = "cluster"
    terraform   = "true"
    repo        = "https://github.com/wiseco/wise-iac"
  }
}
