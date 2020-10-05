export
TF_WORKSPACE_KEY_PREFIX=core-platform-ecr

us-west-2-master:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/core-platform-ecr)
	$(eval export AWS_PROFILE=master-us-west-2-saml-roles-admin)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=master)
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
