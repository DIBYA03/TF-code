variable "userpass" {
  description = "secondary user password ask dibyanshu"
}

variable "adminpass" {
  description = "master admin password ask niting"
}

variable "cadminpass" {
  description = "cluster-admin password to manage cluster"
}


variable "aws_access_key" {}
variable "aws_secret_key" {}


variable "ami_id" {
  default = "ami-78a22900"
}


variable "clustername" {
  default = "sbx-k8s-cluster"
}
