
resource "aws_security_group" "allow_redis" {
  name        = "${var.environment}-allow_redis"
  description = "Allow Redis inbound traffic"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr_block}"]
  }

  tags = {
    Name = "${var.environment}_allow_redis"
  }
}

resource "aws_subnet" "redis-subnet-us-west-2a" {
  vpc_id  = "${var.vpc_id}"
  cidr_block = "${var.cidr_block_us_west_2a}"
  availability_zone = "us-west-2a"

  tags = {
    Name = "wise-us.${var.environment}.us-west-2a-redis"
  }
}

resource "aws_elasticache_subnet_group" "default" {
  name       = "r${var.environment}-redis-cache-subnet"
  subnet_ids = ["${aws_subnet.redis-subnet-us-west-2a.id}"]
}


resource "aws_elasticache_replication_group" "default" {
  replication_group_id          = "${var.environment}-redis-cluster"
  replication_group_description = "Redis cluster for Hashicorp ElastiCache example"

  node_type            = "${var.total_node}"
  port                 = 6379
  parameter_group_name = "default.redis5.0.cluster.on"

  snapshot_retention_limit = 2
  snapshot_window          = "00:00-05:00"

  subnet_group_name          = "${aws_elasticache_subnet_group.default.name}"
  automatic_failover_enabled = true
  engine_version = "${var.redis_engine_version}"

  cluster_mode {
    replicas_per_node_group = 1
    num_node_groups         = "${var.total_node}"
  }
}


