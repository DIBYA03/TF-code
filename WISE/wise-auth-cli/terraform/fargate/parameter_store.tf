data "aws_ssm_parameter" "google_idp_url" {
  name = "/${var.environment}/google/idp_url"
}
