# aws-vpc-terraform

This repo creates the VPCs for each environment. Each environment is provided from the workspace that it's in. For example, you can call the environment as:

```
${var.environment}
```

If you are in the workspace `dev`, then it will return `dev`

## Gotchas
If creating a new VPC, in a new account, and there's no bastion host... you need to manually create an autoscaling group and then remove it before deploying a new VPC.
This is because we use the default AWS generated autoscaling role and doesn't exist until there's an autoscaling group created.


## Preparing to make changes

This will initialize terraform:

```
make init
```

## Check for current setup

This will show the current configuration that is launched in AWS:

```
terraform show
```

## Planning

Before you apply any changes, you always want to check what is going to be changed:

```
make <region>-<environment> plan
```

### Applying changes

Be careful doing this, as it will affect the networking. In fact, if you have any doubts, don't make changes:

```
make <region>-<environment> apply
```
