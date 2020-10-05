# Create a new rancher2 Node Template from Rancher 2.2.x

resource "rancher2_cloud_credential" "node-template-access" {
  provider = rancher2.admin
  name = "nodemaster"
  description = "access to launch node pool"
  amazonec2_credential_config {
    access_key = var.aws_access_key
    secret_key = var.aws_secret_key
  }
}
resource "rancher2_node_template" "master-base" {
  provider = rancher2.admin
  name = "master-base"
  description = "master node pool"
  cloud_credential_id = rancher2_cloud_credential.node-template-access.id
  engine_install_url = "https://releases.rancher.com/install-docker/18.09.sh"
  labels = {
    type = "general"
    compute = "na"
    memory = "na"
    role = "master"
    "cattle.io/creator" = "norman"
  }
  amazonec2_config {
    ami = var.ami_id
    region = "us-west-2"
    security_group = ["rancher-nodes"]
    subnet_id = "subnet-03cb8588897992896"
    vpc_id = "vpc-02f9af6bc9a992cc0"
    zone = "b"
    instance_type = "t2.medium"
    use_private_address = true
    private_address_only = true
    monitoring = true
    iam_instance_profile = "rancher-node-pool"
  }
}
resource "rancher2_node_template" "worker-base-c" {
  provider = rancher2.admin
  name = "worker-base-c"
  description = "master node pool"
  engine_install_url = "https://releases.rancher.com/install-docker/18.09.sh"
  cloud_credential_id = rancher2_cloud_credential.node-template-access.id
  labels = {
    type = "general" 
    compute = "na"
    memory = "na"
    role = "worker"
    "cattle.io/creator" = "norman"
  }
  amazonec2_config {
    ami = var.ami_id
    region = "us-west-2"
    security_group = ["rancher-nodes"]
    subnet_id = "subnet-0d788077e1af4bb58"
    vpc_id = "vpc-02f9af6bc9a992cc0"
    zone = "c"
    instance_type = "m5.large"
    private_address_only = true
    use_private_address = true
    monitoring = true
    iam_instance_profile = "rancher-node-pool"
  }
}

resource "rancher2_node_template" "worker-base-a" {
  provider = rancher2.admin
  name = "worker-base-a"
  description = "master node pool"
  engine_install_url = "https://releases.rancher.com/install-docker/18.09.sh"
  cloud_credential_id = rancher2_cloud_credential.node-template-access.id
  labels = {
    type = "general"
    compute = "na"
    memory = "na"
    role = "worker"
    "cattle.io/creator" = "norman"
  }
  amazonec2_config {
    ami = var.ami_id
    region = "us-west-2"
    security_group = ["rancher-nodes"]
    subnet_id = "subnet-0056a8a30873efb73"
    vpc_id = "vpc-02f9af6bc9a992cc0"
    zone = "a"
    instance_type = "m5.large"
    private_address_only = true
    use_private_address = true
    monitoring = true
    iam_instance_profile = "rancher-node-pool"
  }
}
