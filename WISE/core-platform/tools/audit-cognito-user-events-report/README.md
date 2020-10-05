# audit-cognito-user-events-report

Audits users in Cognito. This does not edit anything in the cognito pool. It reads all the auth events for the specified users and cross checks them against `ONLY` the ones that are being checked. So if there are two numbers being checked, they will only cross check each other.

## What does it do?

Runs the following checks on auth events:

- If IP address was used with another cognito user
- If the event was outside the United States
- If the event was flagged as high risk by cognito

## What environment variables are needed?

- AWS_PROFILE
- AWS_REGION

## How do I run it?

## For only specific numbers:

```
AWS_PROFILE=<aws_profile> AWS_REGION=<aws_region> go run *.go -poolid <user_pool_id> +16263805101 +12063536237
```

## For all users in cognito

```
AWS_PROFILE=<aws_profile> AWS_REGION=<aws_region> go run *.go -poolid <user_pool_id>
```
