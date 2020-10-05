output "client_api_cloudfront_domain_name" {
  value = "${aws_cloudfront_distribution.client_api.domain_name}"
}

output "client_api_domain_name" {
  value = "${module.client_api.name}"
}
