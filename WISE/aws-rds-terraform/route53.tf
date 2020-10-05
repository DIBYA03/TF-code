resource "aws_route53_record" "master" {
  zone_id = "${var.route53_zone_id}"
  name    = "${var.route53_master_domain}"
  type    = "CNAME"
  ttl     = "300"
  records = ["${aws_db_instance.master.address}"]
}

resource "aws_route53_record" "read-replicas" {
  count   = "${var.rds_read_replica_count}"
  zone_id = "${var.route53_zone_id}"
  name    = "${count.index}.${var.route53_read_replica_domain}"
  type    = "CNAME"
  ttl     = "300"
  records = ["${element(aws_db_instance.read_replica.*.address, count.index)}"]
}
