#!/bin/bash

# Get the environment
read -p "Enter environment: " wise_env
read -p "Enter KMS ID: " kms_id
read -p "Enter AWS Profile Name: " profile_name

param_array=(
  "/${wise_env}/bbva/app_env"
  "/${wise_env}/bbva/app_id"
  "/${wise_env}/bbva/app_name"
  "/${wise_env}/bbva/app_secret"
  "/${wise_env}/firebase/config"
  "/${wise_env}/money_request/wise/clearing/business_id"
  "/${wise_env}/money_request/wise/clearing/user_id"
  "/${wise_env}/plaid/client_id"
  "/${wise_env}/plaid/public_key"
  "/${wise_env}/plaid/secret"
  "/${wise_env}/rds/db_port"
  "/${wise_env}/rds/bank_db_name"
  "/${wise_env}/rds/bank_password"
  "/${wise_env}/rds/core_db_name"
  "/${wise_env}/rds/core_username"
  "/${wise_env}/rds/identity_username"
  "/${wise_env}/rds/master_endpoint"
  "/${wise_env}/rds/read_endpoint"
  "/${wise_env}/rds/txn_username"
  "/${wise_env}/rds/bank_username"
  "/${wise_env}/rds/core_password"
  "/${wise_env}/rds/identity_db_name"
  "/${wise_env}/rds/identity_password"
  "/${wise_env}/rds/txn_db_name"
  "/${wise_env}/rds/txn_password"
  "/${wise_env}/redis/endpoint"
  "/${wise_env}/redis/password"
  "/${wise_env}/redis/port"
  "/${wise_env}/segment/write_key"
  "/${wise_env}/sendgrid/api_key"
  "/${wise_env}/stripe/key"
  "/${wise_env}/stripe/publish_key"
  "/${wise_env}/stripe/webhook_secret"
  "/${wise_env}/wise/invoice_email_address"
  "/${wise_env}/wise/support_email_address"
  "/${wise_env}/wise/support_email_name"
  "/${wise_env}/wise/support_phone"
)

for param in "${param_array[@]}"; do
  read -p "Enter value for ${param}: " param_val

  if [[ param_val != "" ]]; then
    aws ssm put-parameter \
      --region us-west-2 \
      --profile $profile_name \
      --name "${param}" \
      --type "SecureString" \
      --key-id $kms_id \
      --value "${param_val}"
  fi

done
