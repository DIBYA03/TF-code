locals {
  cliet_api_gateway_origin_id = "${var.environment}-api-gw-regional"
}

resource "aws_cloudfront_origin_access_identity" "client_api" {
  comment = "${var.environment_name} api gateway origin"
}

resource "aws_cloudfront_distribution" "client_api" {
  aliases     = ["${var.cloudfront_domain_name}"]
  enabled     = "true"
  price_class = "${var.cloudfront_price_class}"

  logging_config {
    include_cookies = false
    bucket          = "${aws_s3_bucket.cloudfront.bucket_domain_name }"
    prefix          = "${var.environment}/api-gateway"
  }

  origin {
    domain_name = "${data.aws_api_gateway_rest_api.client_api.id}.execute-api.${var.aws_region}.amazonaws.com"
    origin_id   = "${local.cliet_api_gateway_origin_id}"

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
    allowed_methods        = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods         = ["GET", "HEAD", "OPTIONS"]
    default_ttl            = 0
    viewer_protocol_policy = "redirect-to-https"
    target_origin_id       = "${local.cliet_api_gateway_origin_id}"

    lambda_function_association {
      event_type = "origin-response"
      lambda_arn = "${module.hsts_lambda.lambda_qualified_arn}"
    }

    forwarded_values {
      query_string = true

      headers = [
        "Access-Control-Allow-Origin",
        "Access-Control-Request-Headers",
        "Access-Control-Request-Method",
        "Authorization",
        "Origin",
        "X-Requested-With",
        "x-wise-business-id",
      ]

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
      restriction_type = "${var.clodfront_country_restriction_type}"
      locations        = "${var.clodfront_country_allowed_country_codes}"
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
