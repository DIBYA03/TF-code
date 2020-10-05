locals {
  ecs_autoscaling_role = "arn:aws:iam::${data.aws_caller_identity.account.account_id}:role/aws-service-role/ecs.application-autoscaling.amazonaws.com/AWSServiceRoleForApplicationAutoScaling_ECSService"
  csp_review_sqs_url   = "https://sqs.${var.aws_region}.amazonaws.com/${data.aws_caller_identity.account.account_id}/${var.environment}-csp-api-review"
}
