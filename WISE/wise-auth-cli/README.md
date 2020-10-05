# wise-auth-cli

Wise Authentication Command Line Interface is a tool that allows a Wise Guy to access AWS resources via command-line and web browser. It does this by managing your credentials through a process called SAML.

This will be a basic run through on how to setup this tool and use it locally for your access requirements.

## Before you begin

 - Reach out to the CloudOps or Security team to request access to AWS. Your supervisor's approval will be needed.
 - You need to have a VPN account setup and also be logged in
 - You will need a Github account and access to the `wise-auth-cli` repo
 - You will also need docker installed locally

## Installation

To install the application locally, you can follow these steps:

```
# get the repo
git clone git@github.com:wiseco/wise-auth-cli.git

# enter the directory for building the container locally
cd wise-auth-cli/cmd/docker/auth-cli-wise

# build the container locally
AWS_DEFAULT_PROFILE=wiseus \
  AWS_DEFAULT_REGION=us-west-2 \
  ECS_ENV=sec \
  BUILD_TAG=latest \
  ECR_IMAGE=379379777492.dkr.ecr.us-west-2.amazonaws.com/sec-auth-cli-wise \
  make build docker-build;
```

Once you have the above complete, you now have the required container built locally. Next, you need to setup a function to run this container.

*NOTE: THIS IS FOR OSX AND LINUX. IF RUNNING WINDOWS, GOOGLE THIS OR REACH OUT TO CLOUDOPS/ SECURITY*

Open your `~/.bash_profile` file for editing and add the following to your profile file:

```
aws-auth-cli-start () {
  docker run -it --rm -d \
    -v ~/.aws/credentials:/app/.aws/credentials \
    --env CONTAINER_LISTEN_PORT=4433 \
    -p 4433:4433 \
    --name aws-auth-cli \
    379379777492.dkr.ecr.us-west-2.amazonaws.com/sec-auth-cli-wise:latest
}

aws-auth-cli-stop () {
  docker stop aws-auth-cli
}
```

Once this is done, from the command-line, run the following:

```
source ~/.bash_profile
```

## Running the application

If you have successfully done everything previous to this, you can now run the following, which will start a local server to handle your auth:

```
aws-auth-cli-start
```

This will run the application in the background (BE AWARE OF THIS). Once you have this, please open a new tab in your browser and go to `https://localhost:4433`

You will see a SSL warning, this is fine and normal. We are currently using a locally generated certificate, so it won't validated.

## Stopping the application

When you would like to stop the local server from running, you can run the following:

```
aws-auth-cli-stop
```

## Containers

### auth-cli-wise

This is the container that is used to handle local credentials on CLI and in browser


### aws-vpn-auth

This container runs in fargate, inside the security account, and holds the source for only allowing login inside the VPN. You can think of this as the middle man
between a user's computer and AWS/Google for the getting the credentials
