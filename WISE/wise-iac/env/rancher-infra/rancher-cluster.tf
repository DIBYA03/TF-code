
# rancher2 RKE Cluster
resource "rancher2_cluster" "rancher-cluster" {
  provider = rancher2.admin
  name = "${local.env}-cluster"
  description = "rancher cluster managed by rancher2 terraform"
  rke_config {
    ignore_docker_version = false
    kubernetes_version    = "v1.16.9-rancher1-1"
    network {
      plugin = "canal"
    } 
    ingress {
      provider = "none"
    }
    services {
      etcd {
        creation = "6h"
        retention = "24h"
      }
    }
    cloud_provider {
      name = "aws"
      aws_cloud_provider {
        global {
          disable_security_group_ingress = false
          disable_strict_zone_check      = false
        }
      }
    }
// ## enable monitoring using prometheus and grafana
//  enable_cluster_monitoring = true
//  cluster_monitoring_input {
//    answers = {
//      "exporter-kubelets.https" = true
//      "exporter-node.enabled" = true
//      "exporter-node.ports.metrics.port" = 9796
//      "exporter-node.resources.limits.cpu" = "200m"
//      "exporter-node.resources.limits.memory" = "200Mi"
//      "grafana.persistence.enabled" = false
//      "grafana.persistence.size" = "10Gi"
//      "grafana.persistence.storageClass" = "default"
//      "operator.resources.limits.memory" = "500Mi"
//      "prometheus.persistence.enabled" = "false"
//      "prometheus.persistence.size" = "50Gi"
//      "prometheus.persistence.storageClass" = "default"
//      "prometheus.persistent.useReleaseName" = "true"
//      "prometheus.resources.core.limits.cpu" = "1000m",
//      "prometheus.resources.core.limits.memory" = "1500Mi"
//      "prometheus.resources.core.requests.cpu" = "750m"
//      "prometheus.resources.core.requests.memory" = "750Mi"
//      "prometheus.retention" = "12h"
//    }
//    version = "0.1.0"
//  }
}

}
