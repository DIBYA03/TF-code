
# Steps to make changes or execute the script  

Make sure you are connected to sbx VPN and the resolv.conf is pointing to nameserver 10.8.0.2

Install terraform version 0.12+ 

execute `terraform init` to initialise the modules

execute `terraform validate ` to check any syntax errors

Export the Rancher user credentials before planning the changes as below  

```
# export Admin user password

export TF_VAR_adminpass=<password>

# export the secondary user password

export TF_VAR_userpass=<password>

# export the cluster-admin profile password

export TF_VAR_cadminpass=<password>

```

Run `terraform plan` to get the changes details

Finally after confirming the changes, run `terraform apply` and approve the changes by entering `yes` when prompted to

Confirm the changes by loggin in to the rancher console at https://rancher.k8s-internal.wise.us


