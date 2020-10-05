module "ec2-iam-role" {
  source  = "Smartbrood/ec2-iam-role/aws"
  name    = "${local.env}-ec2-iam-role"
  version = "0.3.0"

  policy_arn = ["arn:aws:iam::aws:policy/AmazonS3FullAccess",]
}
