locals {
  client_merchant_logo_origin_id = "${var.environment}-merchant-logo-alb"
}

resource "aws_cloudfront_origin_access_identity" "merchant_logo" {
  comment = "${var.environment_name} merchant logo origin"
}

resource "aws_cloudfront_distribution" "merchant_logo" {
  aliases     = ["${var.cloudfront_domain_name}"]
  enabled     = "true"
  price_class = "${var.cloudfront_price_class}"

  origin {
    domain_name = "${var.ecs_merchant_logo_domain}"
    origin_id   = "${local.client_merchant_logo_origin_id}"

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"

      origin_ssl_protocols = [
        "TLSv1.2",
      ]
    }
  }

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD", "OPTIONS"]
    default_ttl            = 1296000                                   # 15 days
    viewer_protocol_policy = "redirect-to-https"
    target_origin_id       = "${local.client_merchant_logo_origin_id}"

    forwarded_values {
      query_string            = true
      query_string_cache_keys = ["name"]

      headers = ["*"]

      cookies {
        forward = "none"
      }
    }
  }

  viewer_certificate {
    ssl_support_method             = "sni-only"
    acm_certificate_arn            = "${module.cert.arn}"
    minimum_protocol_version       = "TLSv1.1_2016"
    cloudfront_default_certificate = false
  }

  restrictions {
    geo_restriction {
      restriction_type = "${var.cloudfront_country_restriction_type}"
      locations        = "${var.cloudfront_country_allowed_country_codes}"
    }
  }

  tags {
    AddAPIGatewayCFWAF = "${var.cloudfront_add_waf}"
    Application        = "${var.application}"
    Environment        = "${var.environment_name}"
    Component          = "${var.component}"
    Name               = "${module.naming.aws_cloudfront_distribution}"
    Team               = "${var.team}"
  }
}
