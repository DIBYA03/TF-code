terraform {
  backend "s3" {
    encrypt = true
    key     = "terraform.tfstate"
  }
}
