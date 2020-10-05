resource "rancher2_node_pool" "master" {
  provider = rancher2.admin
  cluster_id = rancher2_cluster.rancher-cluster.id
  name = "sbx-k8s-master"
  hostname_prefix = "sbx-k8s-master"
  node_template_id = rancher2_node_template.master-base.id
  quantity = 3
  control_plane = true
  etcd = true
  worker = false
}



resource "rancher2_node_pool" "worker-a" {
  provider = rancher2.admin
  cluster_id = rancher2_cluster.rancher-cluster.id
  name = "sbx-k8s-worker-a"
  hostname_prefix = "sbx-k8s-worker"
  node_template_id = rancher2_node_template.worker-base-a.id
  quantity = 1  
  control_plane = false
  etcd = false
  worker = true
}

resource "rancher2_node_pool" "worker-c" {
  provider = rancher2.admin
  cluster_id = rancher2_cluster.rancher-cluster.id
  name = "sbx-k8s-worker-c"
  hostname_prefix = "sbx-k8s-worker"
  node_template_id = rancher2_node_template.worker-base-c.id
  quantity = 1  
  control_plane = false
  etcd = false
  worker = true
}

