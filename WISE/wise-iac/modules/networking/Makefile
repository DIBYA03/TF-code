export
TF_WORKSPACE_KEY_PREFIX=aws-vpc-terraform

# us-west-2 environments
us-west-2-sbx-k8s:
	$(eval export TF_WORKSPACE_KEY_PREFIX=us-west-2/cloud-ops/vpc)
	$(eval export AWS_PROFILE=wiseus)
	$(eval export AWS_REGION=us-west-2)
	$(eval export AWS_DYNAMODB_TABLE=sbx-k8s-terraform.state.lock)
	$(eval export AWS_S3_BUCKET=terraform.sbx-k8s-wise-us.states)
	$(eval export ENV_NAME=sbx-k8s)
	@true

terraform-clean:
	rm -rf \
		terraform.tfstate.d \
		.terraform/environment \
		.terraform/terraform.tfstate

# terraform
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
	terraform plan -var-file=./../../env/sbx-k8s/terraform.tfvars

apply: init
	terraform apply -var-file=./../../env/sbx-k8s/terraform.tfvars
