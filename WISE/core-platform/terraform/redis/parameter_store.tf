data "aws_elasticache_replication_group" "redis" {
  replication_group_id = "redis-cluster"
}

resource "aws_ssm_parameter" "redis_configuration_endpoint_address" {
  name      = "/${var.environment}/redis/configuration_endpoint_address"
  description = "The address of the replication group configuration endpoint when cluster mode is enabled"
  type      = "StringList"
  value     = "${aws_elasticache_replication_group.default.configuration_endpoint_address}"
  overwrite = true
}

