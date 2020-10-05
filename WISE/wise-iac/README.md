# wise-iac
Terraform modules for K8s infrastructure


To make any changes, follow below instructions: 

1. Configure AWS cli with proper privileges to access the AWS resources related to R53, EC2, VPC etc. 

2. Configure the aws credentials profile with name "sbx" as it is coded in the terraform scripts. 
 
3. Pull the latest code from master branch. 

4. For vpc changes, use env/sbx-k8s/ path. for rancher infra changes infra changes, use env/rancher/ path.

5. run `terraform init` to initialise the modules. 

6. Switch to the terraform workspace using `terraform workspace select sbx-k8s`. (currently for the sandbox env)

7. Make the changes in scripts accroding to terraform v0.12 syntax and apply the changes.

8. To apply the changes, run `terraform plan` and validate the changes , followed by `terraform apply` and type `yes` when prompted finally to apply the changes. 

9. To add new modules, add a directory in ./modules/ and refer the code in your scripts. 



  


