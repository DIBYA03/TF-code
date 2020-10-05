locals {
  lambda_invocation_prefix = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.account.account_id}:function:${module.naming.aws_lambda_function}"
}

# Template for OpenAPI file
data "template_file" "openapi" {
  template = "${file("${var.api_gw_openapi_dir_location}/${var.api_gw_stage}/${var.api_gw_openapi_file_source}")}"

  vars = {
    api_name           = "${module.naming.aws_api_gateway_rest_api}"
    server_description = "${var.api_gw_server_description}"

    api_stage   = "${var.api_gw_stage}"
    api_version = "${var.api_gw_stage}"
    domain_name = "${var.route53_domain_name}"

    cognito_pool_arn = "${aws_cognito_user_pool.default.arn}"

    analytics_lambda = "${local.lambda_invocation_prefix}-anc-${var.api_gw_stage}/invocations"

    business_core_lambda                      = "${local.lambda_invocation_prefix}-bus-core-${var.api_gw_stage}/invocations"
    business_account_lambda                   = "${local.lambda_invocation_prefix}-bus-acc-${var.api_gw_stage}/invocations"
    business_external_account_lambda          = "${local.lambda_invocation_prefix}-bus-ext-acc-${var.api_gw_stage}/invocations"
    business_approve_lambda                   = "${local.lambda_invocation_prefix}-bus-apr-${var.api_gw_stage}/invocations"
    business_card_lambda                      = "${local.lambda_invocation_prefix}-bus-crd-${var.api_gw_stage}/invocations"
    business_card_reissue_lambda              = "${local.lambda_invocation_prefix}-bus-crd-ris-${var.api_gw_stage}/invocations"
    business_card_reader_lambda               = "${local.lambda_invocation_prefix}-bus-crd-rdr-${var.api_gw_stage}/invocations"
    business_decline_lambda                   = "${local.lambda_invocation_prefix}-bus-dcn-${var.api_gw_stage}/invocations"
    business_document_lambda                  = "${local.lambda_invocation_prefix}-bus-doc-${var.api_gw_stage}/invocations"
    business_document_formation_lambda        = "${local.lambda_invocation_prefix}-bus-doc-fmn-${var.api_gw_stage}/invocations"
    business_document_status_lambda           = "${local.lambda_invocation_prefix}-bus-doc-sts-${var.api_gw_stage}/invocations"
    business_document_url_lambda              = "${local.lambda_invocation_prefix}-bus-doc-url-${var.api_gw_stage}/invocations"
    business_item_lambda                      = "${local.lambda_invocation_prefix}-bus-itm-${var.api_gw_stage}/invocations"
    business_member_lambda                    = "${local.lambda_invocation_prefix}-bus-mem-${var.api_gw_stage}/invocations"
    business_member_verification_lambda       = "${local.lambda_invocation_prefix}-bus-mem-ver-${var.api_gw_stage}/invocations"
    business_member_phone_verification_lambda = "${local.lambda_invocation_prefix}-bus-mem-ph-ver-${var.api_gw_stage}/invocations"
    business_member_email_verification_lambda = "${local.lambda_invocation_prefix}-bus-mem-em-ver-${var.api_gw_stage}/invocations"
    business_member_alloy_verification_lambda = "${local.lambda_invocation_prefix}-bus-mem-alo-ver-${var.api_gw_stage}/invocations"
    business_member_clear_verification_lambda = "${local.lambda_invocation_prefix}-bus-mem-clr-ver-${var.api_gw_stage}/invocations"
    business_promofunds_lambda                = "${local.lambda_invocation_prefix}-bus-pmf-${var.api_gw_stage}/invocations"
    business_notes_lambda                     = "${local.lambda_invocation_prefix}-bus-nte-${var.api_gw_stage}/invocations"
    business_status_lambda                    = "${local.lambda_invocation_prefix}-bus-sts-${var.api_gw_stage}/invocations"
    business_state_lambda                     = "${local.lambda_invocation_prefix}-bus-ste-${var.api_gw_stage}/invocations"
    business_submit_document_lambda           = "${local.lambda_invocation_prefix}-bus-sdc-${var.api_gw_stage}/invocations"
    business_document_reupload_lambda         = "${local.lambda_invocation_prefix}-bus-doc-rup-${var.api_gw_stage}/invocations"
    business_middesk_lambda                   = "${local.lambda_invocation_prefix}-bus-mdk-${var.api_gw_stage}/invocations"
    business_clear_lambda                     = "${local.lambda_invocation_prefix}-bus-clr-${var.api_gw_stage}/invocations"

    business_csp_lambda                   = "${local.lambda_invocation_prefix}-bus-csp-${var.api_gw_stage}/invocations"
    business_account_closure_lambda       = "${local.lambda_invocation_prefix}-bus-acc-cls-${var.api_gw_stage}/invocations"
    business_account_closure_state_lambda = "${local.lambda_invocation_prefix}-bus-acc-cls-ste-${var.api_gw_stage}/invocations"

    csp_consumer_lambda                   = "${local.lambda_invocation_prefix}-cmr-${var.api_gw_stage}/invocations"
    csp_consumer_document_lambda          = "${local.lambda_invocation_prefix}-cmr-doc-${var.api_gw_stage}/invocations"
    csp_consumer_document_status_lambda   = "${local.lambda_invocation_prefix}-cmr-doc-sts-${var.api_gw_stage}/invocations"
    csp_consumer_document_submit_lambda   = "${local.lambda_invocation_prefix}-cmr-doc-sbt-${var.api_gw_stage}/invocations"
    csp_consumer_document_reupload_lambda = "${local.lambda_invocation_prefix}-cmr-doc-rup-${var.api_gw_stage}/invocations"
    csp_consumer_document_url_lambda      = "${local.lambda_invocation_prefix}-cmr-doc-url-${var.api_gw_stage}/invocations"
    csp_consumer_item_lambda              = "${local.lambda_invocation_prefix}-cmr-itm-${var.api_gw_stage}/invocations"
    csp_consumer_status_lambda            = "${local.lambda_invocation_prefix}-cmr-sts-${var.api_gw_stage}/invocations"
    csp_consumer_state_lambda             = "${local.lambda_invocation_prefix}-cmr-ste-${var.api_gw_stage}/invocations"
    csp_consumer_verification_lambda      = "${local.lambda_invocation_prefix}-cmr-ver-${var.api_gw_stage}/invocations"
    csp_user_change_phone_lambda          = "${local.lambda_invocation_prefix}-usr-chg-ph-${var.api_gw_stage}/invocations"
    csp_business_subscription_lambda      = "${local.lambda_invocation_prefix}-bus-sbc-${var.api_gw_stage}/invocations"

    csp_intercom_lambda          = "${local.lambda_invocation_prefix}-int-${var.api_gw_stage}/invocations"
    csp_intercom_tag_lambda      = "${local.lambda_invocation_prefix}-int-tag-${var.api_gw_stage}/invocations"
    csp_kyc_lambda               = "${local.lambda_invocation_prefix}-kyc-${var.api_gw_stage}/invocations"
    csp_kyb_lambda               = "${local.lambda_invocation_prefix}-kyb-${var.api_gw_stage}/invocations"
    csp_report_lambda            = "${local.lambda_invocation_prefix}-rpt-${var.api_gw_stage}/invocations"
    csp_report_statistics_lambda = "${local.lambda_invocation_prefix}-rpt-sts-${var.api_gw_stage}/invocations"

    csp_posted_transaction_lambda          = "${local.lambda_invocation_prefix}-pst-txn-${var.api_gw_stage}/invocations"
    csp_posted_transaction_export_lambda   = "${local.lambda_invocation_prefix}-pst-txn-exp-${var.api_gw_stage}/invocations"
    csp_pending_transaction_lambda         = "${local.lambda_invocation_prefix}-pnd-txn-${var.api_gw_stage}/invocations"
    csp_pending_transaction_export_lambda  = "${local.lambda_invocation_prefix}-pnd-txn-exp-${var.api_gw_stage}/invocations"
    csp_declined_transaction_lambda        = "${local.lambda_invocation_prefix}-dcl-txn-${var.api_gw_stage}/invocations"
    csp_declined_transaction_export_lambda = "${local.lambda_invocation_prefix}-dcl-txn-exp-${var.api_gw_stage}/invocations"
    csp_transaction_approve_lambda  = "${local.lambda_invocation_prefix}-txn-apr-${var.api_gw_stage}/invocations"
    csp_transaction_decline_lambda  = "${local.lambda_invocation_prefix}-txn-dcn-${var.api_gw_stage}/invocations"
    csp_transaction_transfer_lambda = "${local.lambda_invocation_prefix}-txn-tnf-${var.api_gw_stage}/invocations"

    csp_support_lambda                     = "${local.lambda_invocation_prefix}-sup-${var.api_gw_stage}/invocations"
    csp_support_account_lambda             = "${local.lambda_invocation_prefix}-sup-acc-${var.api_gw_stage}/invocations"
    csp_support_account_block_lambda       = "${local.lambda_invocation_prefix}-sup-acc-blk-${var.api_gw_stage}/invocations"
    csp_support_account_unblock_lambda     = "${local.lambda_invocation_prefix}-sup-acc-ubk-${var.api_gw_stage}/invocations"
    csp_support_phone_lambda               = "${local.lambda_invocation_prefix}-sup-phn-${var.api_gw_stage}/invocations"
    csp_support_card_block_lambda          = "${local.lambda_invocation_prefix}-sup-crd-blk-${var.api_gw_stage}/invocations"
    csp_support_card_block_status_lambda   = "${local.lambda_invocation_prefix}-sup-crd-blk-sts-${var.api_gw_stage}/invocations"
    csp_support_card_unblock_lambda        = "${local.lambda_invocation_prefix}-sup-crd-ubk-${var.api_gw_stage}/invocations"

    review_business_lambda        = "${local.lambda_invocation_prefix}-rvw-bus-${var.api_gw_stage}/invocations"
    review_business_item_lambda   = "${local.lambda_invocation_prefix}-rvw-bus-itm-${var.api_gw_stage}/invocations"
    review_business_status_lambda = "${local.lambda_invocation_prefix}-rvw-bus-sts-${var.api_gw_stage}/invocations"
    review_consumer_lambda        = "${local.lambda_invocation_prefix}-rvw-cmr-${var.api_gw_stage}/invocations"
    review_consumer_item_lambda   = "${local.lambda_invocation_prefix}-rvw-cmr-itm-${var.api_gw_stage}/invocations"

    cognito_pool_arn = "${aws_cognito_user_pool.default.arn}"
  }
}

resource "local_file" "rendered_openapi" {
  filename          = "${var.api_gw_openapi_dir_location}/${var.api_gw_stage}/rendered-${var.api_gw_openapi_file_source}"
  sensitive_content = "${data.template_file.openapi.rendered}"

  file_permission      = "0644"
  directory_permission = "0755"
}

data "aws_iam_policy_document" "csp" {
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

    condition {
      test     = "StringEquals"
      variable = "aws:SourceVpc"

      values = [
        "${var.vpc_id}",
        "${var.shared_vpc_id}",
        "${var.dev_vpc_id}",
      ]
    }
  }
}

resource "aws_api_gateway_rest_api" "csp" {
  name        = "${module.naming.aws_api_gateway_rest_api}"
  description = "${title(var.environment_name)} Wise CSP API"

  policy = "${data.aws_iam_policy_document.csp.json}"

  endpoint_configuration {
    types = ["${var.api_gw_endpoint_configuration}"]
  }

  body = "${data.template_file.openapi.rendered}"
}

# Bug with tf and deployment fix
resource "null_resource" "deploy_api" {
  provisioner "local-exec" {
    command    = "aws apigateway create-deployment --rest-api-id ${aws_api_gateway_rest_api.csp.id} --stage-name ${var.api_gw_stage} --region ${var.aws_region}"
    on_failure = "fail"
  }

  triggers = {
    always_run = "${timestamp()}"
  }
}
