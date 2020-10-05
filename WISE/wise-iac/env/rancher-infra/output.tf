
output "token" {
  value = "${rancher2_bootstrap.admin.token}"
}

output "admin-url" {
  value = "${rancher2_bootstrap.admin.url}"
}
output "cluster_id" {
  value = rancher2_cluster.rancher-cluster.id
}

output "kube_config" {
  value = rancher2_cluster.rancher-cluster.kube_config
}

