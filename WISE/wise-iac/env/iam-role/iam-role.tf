provider "aws" {
  region                  = "us-west-2"
  shared_credentials_file = "~/.aws/credentials"
  profile                 = "wiseus"
}

terraform {
  backend "s3" {
    bucket         = "terraform.sbx-k8s-wise-us.states"
    key            = "us-west-2/cloudops/rancher-iam-role.tfstate"
    region         = "us-west-2"
    dynamodb_table = "sbx-k8s-terraform.state.lock"
    shared_credentials_file = "~/.aws/credentials"
    profile                 = "wiseus"
  }
}

resource "aws_iam_role" "EKS-sbx-k8s-role" {
  name = "EKS-sbx-k8s-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "eks.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags = {
    Name = "k8s-sbx EKS IAM role"
  }
}

resource "aws_iam_role_policy_attachment" "AmazonEKSClusterPolicy" {
  role       = "${aws_iam_role.EKS-sbx-k8s-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}
