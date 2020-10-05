provider "rancher2" {
  api_url   = "https://rancher.sbx.wise.us"
  bootstrap = true
}


resource "rancher2_bootstrap" "admin" {
  password = var.adminpass
  telemetry = true
}


# Provider config for admin
provider "rancher2" {
  alias = "admin"
  api_url = "https://rancher.sbx.wise.us"
  token_key = rancher2_bootstrap.admin.token
  insecure = true
}

//output "secret_key" {
//  value = "${rancher2_bootstrap.admin.secret_key}"
//}
//output "access_key" {
//  value = "${rancher2_bootstrap.admin.access_key}"
//}

# Create a new rancher2 User
resource "rancher2_user" "dibyanshu" {
  provider = rancher2.admin
  name = "dibyanshu"
  username = "dibyanshu"
  password = var.userpass
  enabled = true
}
# Create a new rancher2 global_role_binding for User
resource "rancher2_global_role_binding" "dibyanshu" {
  provider = rancher2.admin
  name = "admin-role"
  global_role_id = "admin"
  user_id = "rancher2_user.dibyanshu.id"
}
# Create a new rancher2 User for cluster profile
resource "rancher2_user" "cluster-admin" {
#  depends_on = ["${module.ec2_cluster.private_ip}"]
  provider = rancher2.admin
  name = "cluster-admin"
  username = "cluster-admin"
  password = var.cadminpass
  enabled = true
}
# Create a new rancher2 global_role_binding for User cluster-admin
resource "rancher2_global_role_binding" "cluster-admin" {
  provider = rancher2.admin
  name = "cluster-admin-role"
  global_role_id = "admin"
  user_id = "rancher2_user.cluster-admin.id"
}
