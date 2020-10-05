export
TF_WORKSPACE_KEY_PREFIX=wise-auth-cli

# run `make (env) plan` etc
us-west-2-security:
	$(eval export AWS_PROFILE=security-us-west-2-saml-roles-admin)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=security-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-west-2-security-terraform-config)
	$(eval export ENV_NAME=security)
	@true

# terraform
terraform-clean:
	rm -rf \
		terraform.tfstate.d \
		.terraform/environment \
		.terraform/terraform.tfstate

init: terraform-clean
	terraform init \
		-reconfigure \
		-backend-config="bucket=${AWS_S3_BUCKET}" \
		-backend-config="workspace_key_prefix=${TF_WORKSPACE_KEY_PREFIX}" \
		-backend-config="dynamodb_table=${AWS_DYNAMODB_TABLE}" \
		-backend-config="profile=${AWS_PROFILE}" \
		-backend-config="region=${AWS_REGION}"
	- terraform workspace new ${ENV_NAME}
	terraform workspace select ${ENV_NAME}

show:
	terraform show

plan: init
	terraform plan -var-file=./${AWS_REGION}/${ENV_NAME}/terraform.tfvars

apply: init
	terraform apply -var-file=./${AWS_REGION}/${ENV_NAME}/terraform.tfvars
