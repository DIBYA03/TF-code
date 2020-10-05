locals {
  api_env {
    API_ENV = "${var.environment_name}"
  }

  aws_s3_bucket_document {
    AWS_S3_BUCKET_DOCUMENT = "${data.aws_s3_bucket.documents.id}"
  }

  bbva_app_credentials = {
    BBVA_APP_ENV    = "${data.aws_ssm_parameter.bbva_app_env.value}"
    BBVA_APP_ID     = "${data.aws_ssm_parameter.bbva_app_id.value}"
    BBVA_APP_NAME   = "${data.aws_ssm_parameter.bbva_app_name.value}"
    BBVA_APP_SECRET = "${data.aws_ssm_parameter.bbva_app_secret.value}"
  }

  core_db_credentials = {
    CORE_DB_WRITE_URL  = "${data.aws_ssm_parameter.rds_master_endpoint.value}"
    CORE_DB_READ_URL   = "${data.aws_ssm_parameter.rds_read_endpoint.value}"
    CORE_DB_WRITE_PORT = "${data.aws_ssm_parameter.rds_port.value}"
    CORE_DB_READ_PORT  = "${data.aws_ssm_parameter.rds_port.value}"
    CORE_DB_NAME       = "${data.aws_ssm_parameter.core_rds_db_name.value}"
    CORE_DB_USER       = "${data.aws_ssm_parameter.core_rds_user_name.value}"
    CORE_DB_PASSWD     = "${data.aws_ssm_parameter.core_rds_password.value}"
  }

  csp_db_credentials = {
    CSP_DB_WRITE_URL  = "${data.aws_ssm_parameter.csp_rds_master_endpoint.value}"
    CSP_DB_READ_URL   = "${data.aws_ssm_parameter.csp_rds_read_endpoint.value}"
    CSP_DB_WRITE_PORT = "${data.aws_ssm_parameter.csp_rds_port.value}"
    CSP_DB_READ_PORT  = "${data.aws_ssm_parameter.csp_rds_port.value}"
    CSP_DB_NAME       = "${data.aws_ssm_parameter.csp_rds_db_name.value}"
    CSP_DB_USER       = "${data.aws_ssm_parameter.csp_rds_username.value}"
    CSP_DB_PASSWD     = "${data.aws_ssm_parameter.csp_rds_password.value}"
  }

  bank_db_credentials = {
    BANK_DB_WRITE_URL  = "${data.aws_ssm_parameter.rds_master_endpoint.value}"
    BANK_DB_READ_URL   = "${data.aws_ssm_parameter.rds_read_endpoint.value}"
    BANK_DB_WRITE_PORT = "${data.aws_ssm_parameter.rds_port.value}"
    BANK_DB_READ_PORT  = "${data.aws_ssm_parameter.rds_port.value}"
    BANK_DB_NAME       = "${data.aws_ssm_parameter.bank_rds_db_name.value}"
    BANK_DB_USER       = "${data.aws_ssm_parameter.bank_rds_user_name.value}"
    BANK_DB_PASSWD     = "${data.aws_ssm_parameter.bank_rds_password.value}"
  }

  identity_db_credentials = {
    IDENTITY_DB_WRITE_URL  = "${data.aws_ssm_parameter.rds_master_endpoint.value}"
    IDENTITY_DB_READ_URL   = "${data.aws_ssm_parameter.rds_read_endpoint.value}"
    IDENTITY_DB_WRITE_PORT = "${data.aws_ssm_parameter.rds_port.value}"
    IDENTITY_DB_READ_PORT  = "${data.aws_ssm_parameter.rds_port.value}"
    IDENTITY_DB_NAME       = "${data.aws_ssm_parameter.identity_rds_db_name.value}"
    IDENTITY_DB_USER       = "${data.aws_ssm_parameter.identity_rds_user_name.value}"
    IDENTITY_DB_PASSWD     = "${data.aws_ssm_parameter.identity_rds_password.value}"
  }

  txn_db_credentials = {
    TXN_DB_WRITE_URL  = "${data.aws_ssm_parameter.rds_master_endpoint.value}"
    TXN_DB_READ_URL   = "${data.aws_ssm_parameter.rds_read_endpoint.value}"
    TXN_DB_WRITE_PORT = "${data.aws_ssm_parameter.rds_port.value}"
    TXN_DB_READ_PORT  = "${data.aws_ssm_parameter.rds_port.value}"
    TXN_DB_NAME       = "${data.aws_ssm_parameter.txn_rds_db_name.value}"
    TXN_DB_USER       = "${data.aws_ssm_parameter.txn_rds_user_name.value}"
    TXN_DB_PASSWD     = "${data.aws_ssm_parameter.txn_rds_password.value}"

    KINESIS_TRX_REGION = "${var.txn_kinesis_region}"
    KINESIS_TRX_NAME   = "${var.txn_kinesis_name}"
  }

  grpc_container_ports = {
    GRPC_SERVICE_PORT = "${var.grpc_port}"
  }

  use_transaction_service = {
    USE_TRANSACTION_SERVICE = "${var.use_transaction_service}"
  }

  use_banking_service = {
    USE_BANKING_SERVICE = "${var.use_banking_service}"
  }

  use_invoice_service = {
    USE_INVOICE_SERVICE = "${var.use_invoice_service}"
  }

  business_document_sqs = {
    CSP_SQS_REGION = "${var.aws_region}"
    CSP_SQS_URL    = "${data.aws_sqs_queue.business_document_upload.id}"
  }

  review_sqs = {
    CSP_REVIEW_SQS_REGION = "${var.aws_region}"
    CSP_REVIEW_SQS_URL    = "${data.aws_sqs_queue.review.id}"
  }

  segment_sqs = {
    SQS_REGION      = "${var.aws_region}"
    SEGMENT_SQS_URL = "${data.aws_sqs_queue.segment_analytics.id}"
  }

  sendgrid_api = {
    SENDGRID_API_KEY = "${data.aws_ssm_parameter.sendgrid_api_key.value}"
  }

  twilio_vars = {
    TWILIO_ACCOUNT_SID  = "${data.aws_ssm_parameter.twilio_account_sid.value}"
    TWILIO_API_SID      = "${data.aws_ssm_parameter.twilio_api_sid.value}"
    TWILIO_API_SECRET   = "${data.aws_ssm_parameter.twilio_api_secret.value}"
    TWILIO_SENDER_PHONE = "${data.aws_ssm_parameter.twilio_sender_phone.value}"
  }

  wise_clearing_ids = {
    WISE_CLEARING_ACCOUNT_ID  = "${data.aws_ssm_parameter.wise_clearing_account_id.value}"
    WISE_CLEARING_BUSINESS_ID = "${data.aws_ssm_parameter.wise_clearing_business_id.value}"
    WISE_CLEARING_USER_ID     = "${data.aws_ssm_parameter.wise_clearing_user_id.value}"
  }

  wise_promo_clearing_ids = {
    WISE_PROMO_CLEARING_ACCOUNT_ID        = "${data.aws_ssm_parameter.wise_promo_clearing_account_id.value}"
    WISE_PROMO_CLEARING_LINKED_ACCOUNT_ID = "${data.aws_ssm_parameter.wise_promo_linked_clearing_account_id.value}"
  }

  wise_support_email = {
    WISE_SUPPORT_EMAIL = "${data.aws_ssm_parameter.wise_support_email_address.value}"
    WISE_SUPPORT_NAME  = "${data.aws_ssm_parameter.wise_support_email_name.value}"
  }

  intercom_access_token = {
    INTERCOM_ACCESS_TOKEN = "${data.aws_ssm_parameter.intercom_access_token.value}"
  }

  cognito_user_pool_id = {
    COGNITO_USER_POOL_ID = "${data.aws_ssm_parameter.cognito_user_pool_id.value}"
  }

  vgs = {
    VGS_CERT            = "${data.aws_ssm_parameter.vgs_cert.value}"
    VGS_HTTPS_PROXY_URL = "${data.aws_ssm_parameter.vgs_https_proxy_url.value}"
  }

  batch_default_timezone = {
    BATCH_TZ            = "${var.batch_default_timezone}"
  }
}
