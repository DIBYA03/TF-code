locals {
  lambda_invocation_prefix = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.account.account_id}:function:${module.naming.aws_lambda_function}"
}

data "template_file" "statement_request" {
  template = "${file("${path.module}/specs/request_templates/statement.json")}"
}

# Template for OpenAPI file
data "template_file" "openapi" {
  template = "${file("${var.api_gw_openapi_dir_location}/${var.api_gw_stage}/${var.api_gw_openapi_file_source}")}"

  vars = {
    api_stage          = "${var.api_gw_stage}"
    api_version        = "${var.api_gw_stage}"                       # remove this when removed from openapi specs
    domain_name        = "${var.route53_domain_name}"
    api_name           = "${module.naming.aws_api_gateway_rest_api}"
    server_description = "${var.api_gw_server_description}"

    banking_business_account_lambda                     = "${local.lambda_invocation_prefix}-ban-bus-acc-${var.api_gw_stage}/invocations"
    banking_business_account_balance_lambda             = "${local.lambda_invocation_prefix}-ban-bus-acc-bal-${var.api_gw_stage}/invocations"
    banking_business_account_card_lambda                = "${local.lambda_invocation_prefix}-ban-bus-acc-crd-${var.api_gw_stage}/invocations"
    banking_business_account_statement_lambda           = "${local.lambda_invocation_prefix}-ban-bus-acc-smt-${var.api_gw_stage}/invocations"
    banking_business_account_statement_request_template = "${jsonencode(data.template_file.statement_request.rendered)}"
    banking_business_account_statement_list_lambda      = "${local.lambda_invocation_prefix}-ban-bus-acc-smt-lst-${var.api_gw_stage}/invocations"

    banking_business_card_activation_lambda = "${local.lambda_invocation_prefix}-ban-bus-crd-acv-${var.api_gw_stage}/invocations"
    banking_business_card_block_lambda      = "${local.lambda_invocation_prefix}-ban-bus-crd-blk-${var.api_gw_stage}/invocations"
    banking_business_card_unblock_lambda    = "${local.lambda_invocation_prefix}-ban-bus-crd-ubk-${var.api_gw_stage}/invocations"

    banking_business_connectaccount_lambda = "${local.lambda_invocation_prefix}-ban-bus-cta-${var.api_gw_stage}/invocations"

    banking_business_contact_linkedaccount_lambda = "${local.lambda_invocation_prefix}-ban-bus-con-lka-${var.api_gw_stage}/invocations"
    banking_business_contact_linkedcard_lambda    = "${local.lambda_invocation_prefix}-ban-bus-con-lkc-${var.api_gw_stage}/invocations"
    banking_business_contact_moneyrequest_lambda  = "${local.lambda_invocation_prefix}-ban-bus-con-myr-${var.api_gw_stage}/invocations"

    banking_business_linkedaccount_lambda = "${local.lambda_invocation_prefix}-ban-bus-lka-${var.api_gw_stage}/invocations"
    banking_business_linkedcard_lambda = "${local.lambda_invocation_prefix}-ban-bus-lkc-${var.api_gw_stage}/invocations"

    banking_business_transaction_lambda                   = "${local.lambda_invocation_prefix}-ban-bus-tns-${var.api_gw_stage}/invocations"
    banking_business_transaction_dispute_lambda           = "${local.lambda_invocation_prefix}-ban-bus-tns-dsp-${var.api_gw_stage}/invocations"
    banking_business_transaction_export_lambda            = "${local.lambda_invocation_prefix}-ban-bus-tns-exp-${var.api_gw_stage}/invocations"
    banking_business_transaction_receipt_lambda           = "${local.lambda_invocation_prefix}-ban-bus-tns-rct-${var.api_gw_stage}/invocations"
    banking_business_transaction_receipt_signedurl_lambda = "${local.lambda_invocation_prefix}-ban-bus-tns-rct-srl-${var.api_gw_stage}/invocations"

    banking_business_moneytransfer_lambda = "${local.lambda_invocation_prefix}-ban-bus-myt-${var.api_gw_stage}/invocations"

    banking_pending_transaction_lambda        = "${local.lambda_invocation_prefix}-ban-pnd-txn-${var.api_gw_stage}/invocations"
    banking_pending_transaction_export_lambda = "${local.lambda_invocation_prefix}-ban-pnd-txn-exp-${var.api_gw_stage}/invocations"

    business_lambda                       = "${local.lambda_invocation_prefix}-bus-${var.api_gw_stage}/invocations"
    business_activity_lambda              = "${local.lambda_invocation_prefix}-bus-act-${var.api_gw_stage}/invocations"
    business_contact_lambda               = "${local.lambda_invocation_prefix}-bus-con-${var.api_gw_stage}/invocations"
    business_contact_moneytransfer_lambda = "${local.lambda_invocation_prefix}-bus-con-myt-${var.api_gw_stage}/invocations"
    business_contact_address_lambda       = "${local.lambda_invocation_prefix}-bus-con-adr-${var.api_gw_stage}/invocations"
    business_document_lambda              = "${local.lambda_invocation_prefix}-bus-doc-${var.api_gw_stage}/invocations"
    business_document_signedurl_lambda    = "${local.lambda_invocation_prefix}-bus-doc-srl-${var.api_gw_stage}/invocations"
    business_member_lambda                = "${local.lambda_invocation_prefix}-bus-mem-${var.api_gw_stage}/invocations"
    business_member_submission_lambda     = "${local.lambda_invocation_prefix}-bus-mem-sub-${var.api_gw_stage}/invocations"
    business_signature_lambda             = "${local.lambda_invocation_prefix}-bus-sig-${var.api_gw_stage}/invocations"
    business_submission_lambda            = "${local.lambda_invocation_prefix}-bus-sub-${var.api_gw_stage}/invocations"
    business_partner_lambda               = "${local.lambda_invocation_prefix}-bus-ptr-${var.api_gw_stage}/invocations"
    business_account_closure_lambda       = "${local.lambda_invocation_prefix}-bus-acc-cls-${var.api_gw_stage}/invocations"
    business_subscription_lambda          = "${local.lambda_invocation_prefix}-bus-sbc-${var.api_gw_stage}/invocations"

    consumer_document_lambda     = "${local.lambda_invocation_prefix}-cmr-doc-${var.api_gw_stage}/invocations"
    consumer_document_url_lambda = "${local.lambda_invocation_prefix}-cmr-doc-url-${var.api_gw_stage}/invocations"

    delegate_lambda = "${local.lambda_invocation_prefix}-dlg-${var.api_gw_stage}/invocations"

    payment_capture_lambda          = "${local.lambda_invocation_prefix}-pmt-cpr-${var.api_gw_stage}/invocations"
    payment_card_reader_lambda      = "${local.lambda_invocation_prefix}-pmt-crd-rdr-${var.api_gw_stage}/invocations"
    payment_connection_token_lamba  = "${local.lambda_invocation_prefix}-pmt-ctn-tkn-${var.api_gw_stage}/invocations"
    payment_invoice_lambda          = "${local.lambda_invocation_prefix}-pmt-inv-${var.api_gw_stage}/invocations"
    payment_receipt_lambda          = "${local.lambda_invocation_prefix}-pmt-rct-${var.api_gw_stage}/invocations"
    payment_request_lambda          = "${local.lambda_invocation_prefix}-pmt-rqt-${var.api_gw_stage}/invocations"
    payment_request_resend_lambda   = "${local.lambda_invocation_prefix}-pmt-rqt-rsd-${var.api_gw_stage}/invocations"
    payment_transfer_request_lambda = "${local.lambda_invocation_prefix}-pmt-tnf-rqt-${var.api_gw_stage}/invocations"

    user_lambda                         = "${local.lambda_invocation_prefix}-usr-${var.api_gw_stage}/invocations"
    user_delete_lambda                  = "${local.lambda_invocation_prefix}-usr-del-${var.api_gw_stage}/invocations"
    user_device_logout_lambda           = "${local.lambda_invocation_prefix}-usr-dev-lgo-${var.api_gw_stage}/invocations"
    user_device_pushregistration_lambda = "${local.lambda_invocation_prefix}-usr-dev-psr-${var.api_gw_stage}/invocations"
    user_document_lambda                = "${local.lambda_invocation_prefix}-usr-doc-${var.api_gw_stage}/invocations"
    user_document_url_lambda            = "${local.lambda_invocation_prefix}-usr-doc-url-${var.api_gw_stage}/invocations"
    user_submission_lambda              = "${local.lambda_invocation_prefix}-usr-sub-${var.api_gw_stage}/invocations"
    user_activity_lambda                = "${local.lambda_invocation_prefix}-usr-act-${var.api_gw_stage}/invocations"
    user_notification_lambda            = "${local.lambda_invocation_prefix}-usr-ntf-${var.api_gw_stage}/invocations"
    user_partner_lambda                 = "${local.lambda_invocation_prefix}-usr-ptr-${var.api_gw_stage}/invocations"
    user_self_lambda                    = "${local.lambda_invocation_prefix}-usr-slf-${var.api_gw_stage}/invocations"

    cognito_pool_arn = "${aws_cognito_user_pool.default.arn}"
  }
}

resource "local_file" "rendered_openapi" {
  filename          = "${var.api_gw_openapi_dir_location}/${var.api_gw_stage}/rendered-${var.api_gw_openapi_file_source}"
  sensitive_content = "${data.template_file.openapi.rendered}"

  file_permission      = "0644"
  directory_permission = "0755"
}

data "aws_iam_policy_document" "client" {
  statement {
    principals {
      type        = "*"
      identifiers = ["*"]
    }

    actions = [
      "execute-api:Invoke",
    ]

    resources = [
      "execute-api:/*",
    ]
  }
}

resource "aws_api_gateway_rest_api" "client" {
  name        = "${module.naming.aws_api_gateway_rest_api}"
  description = "${title(var.environment_name)} Wise Client API"

  policy = "${data.aws_iam_policy_document.client.json}"

  endpoint_configuration {
    types = ["${var.api_gw_endpoint_configuration}"]
  }

  body = "${data.template_file.openapi.rendered}"
}

# Bug with tf and deployment fix
resource "null_resource" "deploy_api" {
  provisioner "local-exec" {
    command    = "aws apigateway create-deployment --rest-api-id ${aws_api_gateway_rest_api.client.id} --stage-name ${var.api_gw_stage} --region ${var.aws_region}"
    on_failure = "fail"
  }

  triggers = {
    always_run = "${timestamp()}"
  }
}
