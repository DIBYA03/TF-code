export
TF_WORKSPACE_KEY_PREFIX=core-platform-csp-fargate

# run `make (env) plan` etc
us-west-2-dev:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/csp-api-fargate)
	$(eval export AWS_PROFILE=master-us-west-2-saml-roles-deployment)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=dev)
	$(eval export OLD_WORKSPACE=us-west-2-dev)
	@true

us-west-2-qa:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/csp-api-fargate)
	$(eval export AWS_PROFILE=master-us-west-2-saml-roles-deployment)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=qa)
	$(eval export OLD_WORKSPACE=us-west-2-qa)
	@true

us-west-2-staging:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/csp-api-fargate)
	$(eval export AWS_PROFILE=master-us-west-2-saml-roles-deployment)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=staging)
	$(eval export OLD_WORKSPACE=us-west-2-staging)
	@true

us-west-2-prod:
	$(eval export AWS_PROFILE=prod-us-west-2-saml-roles-deployment)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=prod-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-west-2-prod-terraform-config)
	$(eval export ENV_NAME=prod)
	$(eval export OLD_WORKSPACE=prod)
	@true

us-west-2-sbx:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/csp-api-fargate)
	$(eval export AWS_PROFILE=wiseus)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=sbx-terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.sbx-wise-us.states)
	$(eval export ENV_NAME=sbx)
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
	- terraform workspace new ${OLD_WORKSPACE}
	terraform workspace select ${OLD_WORKSPACE}

show:
	terraform show

plan: init
	terraform plan -var-file=./${AWS_REGION}/${ENV_NAME}/terraform.tfvars -var "api_gw_stage=${STAGE}"

apply: init
	terraform apply -var-file=./${AWS_REGION}/${ENV_NAME}/terraform.tfvars -var "api_gw_stage=${STAGE}"

apply-auto-approve: init
	terraform apply -var-file=./${AWS_REGION}/${ENV_NAME}/terraform.tfvars -var "api_gw_stage=${STAGE}" -auto-approve
