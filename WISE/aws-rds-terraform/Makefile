export
TF_WORKSPACE_KEY_PREFIX=aws-rds-terraform

# run `make (env) plan` etc
us-west-2-csp:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/rds)
	$(eval export AWS_PROFILE=wiseus)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=csp)
	@true

us-west-2-csp-prod:
	$(eval export AWS_PROFILE=wise-prod)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=prod-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-west-2-prod-terraform-config)
	$(eval export ENV_NAME=csp-prod)
	@true

us-west-2-dev:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/rds)
	$(eval export AWS_PROFILE=wiseus)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=dev)
	@true

us-west-2-dev1:
	$(eval export AWS_PROFILE=wise-dev)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=dev-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-west-2-dev-terraform-config)
	$(eval export ENV_NAME=dev1)
	@true

us-west-2-staging:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/rds)
	$(eval export AWS_PROFILE=wiseus)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=staging)
	@true

us-west-2-prod:
	$(eval export AWS_PROFILE=wise-prod)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=prod-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-west-2-prod-terraform-config)
	$(eval export ENV_NAME=prod)
	@true

# terraform commands
terraform-clean:
	rm -rf \
		terraform.tfstate.d \
		.terraform/environment \
		.terraform/terraform.tfstate

# terraform
init: terraform-clean
	terraform init \
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
	terraform plan -var-file=./${AWS_REGION}/${ENV_NAME}/terraform.tfvars -var-file=./${AWS_REGION}/${ENV_NAME}/secrets.tfvars

apply: init
	terraform apply -var-file=./${AWS_REGION}/${ENV_NAME}/terraform.tfvars -var-file=./${AWS_REGION}/${ENV_NAME}/secrets.tfvars
