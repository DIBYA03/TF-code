output "merchant_logo_cloudfront_domain_name" {
  value = "${aws_cloudfront_distribution.merchant_logo.domain_name}"
}

output "merchant_logo_domain_name" {
  value = "${module.merchant_logo.name}"
}
