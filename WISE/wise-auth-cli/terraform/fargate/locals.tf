locals {
  ecs_autoscaling_role = "arn:aws:iam::${data.aws_caller_identity.account.account_id}:role/aws-service-role/ecs.application-autoscaling.amazonaws.com/AWSServiceRoleForApplicationAutoScaling_ECSService"
}
