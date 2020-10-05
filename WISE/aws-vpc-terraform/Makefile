export
TF_WORKSPACE_KEY_PREFIX=aws-vpc-terraform

# us-west-2 environments
us-west-2-csp:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/vpc)
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
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/vpc)
	$(eval export AWS_PROFILE=master-us-west-2-saml-roles-admin)
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
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/vpc)
	$(eval export AWS_PROFILE=wiseus)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=staging)
	@true

us-west-2-shared:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/vpc)
	$(eval export AWS_PROFILE=wiseus)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.wise-us.states)
	$(eval export ENV_NAME=shared)
	@true

us-west-2-prod:
	$(eval export AWS_PROFILE=wise-prod)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=prod-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-west-2-prod-terraform-config)
	$(eval export ENV_NAME=prod)
	@true

us-west-2-security:
	$(eval export AWS_PROFILE=wise-security)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=security-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-west-2-security-terraform-config)
	$(eval export ENV_NAME=security)
	@true

# us-east-1 environments
us-east-1-csp-prod:
	$(eval export AWS_PROFILE=wise-prod)
	$(eval export AWS_REGION=us-east-1)
	$(eval export AWS_DYNAMODB_TABLE=prod-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-east-1-prod-terraform-config)
	$(eval export ENV_NAME=csp-prod)
	@true

us-east-1-prod:
	$(eval export AWS_PROFILE=wise-prod)
	$(eval export AWS_REGION=us-east-1)
	$(eval export AWS_DYNAMODB_TABLE=prod-terraform-config.lock)
	$(eval export AWS_S3_BUCKET=wiseus-us-east-1-prod-terraform-config)
	$(eval export ENV_NAME=prod)
	@true

saltstack-clean-bastion:
	rm -f saltstack/bastion/saltstack.zip

saltstack-bastion: saltstack-clean-bastion
	cd saltstack/bastion && zip -r saltstack.zip ./*

terraform-clean:
	rm -rf \
		terraform.tfstate.d \
		.terraform/environment \
		.terraform/terraform.tfstate

# terraform
init: saltstack-bastion terraform-clean
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
	terraform apply -var-file=./${AWS_REGION}/${ENV_NAME}/terraform.tfvars -var-file=./${AWS_REGION}/${ENV_NAME}/secrets.tfvars
