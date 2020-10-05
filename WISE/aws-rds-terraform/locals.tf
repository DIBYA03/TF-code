locals {
  # max connections
  # 1 GiB = 1073741824 bytes
  rds_max_connections = "${floor(((var.rds_instance_class_mem * 1073741824) / var.rds_instance_max_connection_divider))}"

  rds_non_critical_conn_count_limit = "${floor(local.rds_max_connections *
    (__builtin_StringToFloat(var.rds_cw_conn_count_non_critical) / __builtin_StringToFloat(100)))}"

  rds_critical_conn_count_limit = "${floor(local.rds_max_connections *
    (__builtin_StringToFloat(var.rds_cw_conn_count_critical) / __builtin_StringToFloat(100)))}"

  # disk space
  # 1 GiB = 1073741824 bytes
  rds_disk_space_mb = "${floor(var.rds_storage_size * 1073741824)}"

  rds_non_critical_disk_limit = "${floor(local.rds_disk_space_mb *
    (__builtin_StringToFloat(var.rds_cw_free_disk_non_critical) / __builtin_StringToFloat(100)))}"

  #  40MB
  #  40
  rds_critical_disk_limit = "${floor(local.rds_disk_space_mb *
    (__builtin_StringToFloat(var.rds_cw_free_disk_critical) / __builtin_StringToFloat(100)))}"

  # memory
  # 1 GiB = 1073741824 bytes
  rds_mem_mb = "${floor(var.rds_instance_class_mem * 1073741824)}"

  rds_non_critical_mem_limit = "${floor(local.rds_mem_mb *
    (__builtin_StringToFloat(var.rds_cw_free_mem_non_critical) / __builtin_StringToFloat(100)))}"

  rds_critical_mem_limit = "${floor(local.rds_mem_mb *
    (__builtin_StringToFloat(var.rds_cw_free_mem_critical) / __builtin_StringToFloat(100)))}"
}
