# This is a hack around getting the route53 entries in. We can't loop the vpc endpoints
# to get the dns_name and hosted_zone_id. So, to work around, we will run TF twice.
# One to create the bastion host endpoints and other to add the route53 here.
# Sucks, but will fix when we upgrade to TF .12

resource "aws_route53_record" "prod_bastion" {
  count = "${length(var.bastion_host_vpc_endpoint_service_list) >= 1 ? 1 : 0}"

  zone_id = "${var.internal_route53_zone_id}"
  name    = "prod-bastion.internal.wise.us"
  type    = "A"

  alias {
    name    = "vpce-0e690c4ceec89397f-zi296is7.vpce-svc-0c71a43062d4901f1.us-west-2.vpce.amazonaws.com"
    zone_id = "Z1YSA3EXCYUU9Z"

    evaluate_target_health = true
  }
}

resource "aws_route53_record" "csp_prod_bastion" {
  count = "${length(var.bastion_host_vpc_endpoint_service_list) >= 1 ? 1 : 0}"

  zone_id = "${var.internal_route53_zone_id}"
  name    = "csp-prod-bastion.internal.wise.us"
  type    = "A"

  alias {
    name    = "vpce-04e1d8e445ec4d80e-5hhidbvi.vpce-svc-0bafa0c126600ae70.us-west-2.vpce.amazonaws.com"
    zone_id = "Z1YSA3EXCYUU9Z"

    evaluate_target_health = true
  }
}

resource "aws_route53_record" "security_bastion" {
  count = "${length(var.bastion_host_vpc_endpoint_service_list) >= 1 ? 1 : 0}"

  zone_id = "${var.internal_route53_zone_id}"
  name    = "security-bastion.internal.wise.us"
  type    = "A"

  alias {
    name    = "vpce-031fbb543de144d13-knoa01zu.vpce-svc-09eb3760f1aeef91e.us-west-2.vpce.amazonaws.com"
    zone_id = "Z1YSA3EXCYUU9Z"

    evaluate_target_health = true
  }
}

resource "aws_route53_record" "private_ca" {
  count = "${length(var.bastion_host_vpc_endpoint_service_list) >= 1 ? 1 : 0}"

  zone_id = "${var.internal_route53_zone_id}"
  name    = "ca.internal.wise.us"
  type    = "A"

  alias {
    name    = "vpce-09963168148e6ecf6-kxpwbr0n.vpce-svc-0e2c686c72ccf1cfb.us-west-2.vpce.amazonaws.com"
    zone_id = "Z1YSA3EXCYUU9Z"

    evaluate_target_health = true
  }
}
