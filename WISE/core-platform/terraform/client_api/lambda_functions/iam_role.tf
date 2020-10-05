data "aws_iam_role" "lambda_default" {
  name = "${module.naming.aws_iam_role}-lambda"
}

# IAM role for delete lambda function
resource "aws_iam_role" "clientapi_delete_lambda" {
  name = "${module.naming.aws_iam_role}-del-lambda-${var.api_gw_stage}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_iam_role}-del-lambda-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

data "template_file" "clientapi_delete_lambda" {
  template = "${file("./policies/clientapi-delete-lambda.json")}"
}

resource "aws_iam_role_policy" "clientapi_delete_lambda" {
  name   = "${module.naming.aws_iam_policy}-del-lambda-${var.api_gw_stage}"
  role   = "${aws_iam_role.clientapi_delete_lambda.id}"
  policy = "${data.template_file.clientapi_delete_lambda.rendered}"
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "clientapi_delete_lambda_cw" {
  role       = "${aws_iam_role.clientapi_delete_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "clientapi_delete_lambda_vpc" {
  role       = "${aws_iam_role.clientapi_delete_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "clientapi_delete_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.clientapi_delete_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToInternalBankingSQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey",
        "sqs:GetQueue*",
        "sqs:SendMessage*"
      ],
      "Resource": [
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_sqs_queue.segment_analytics.arn}"
      ]
    }
  ]
}
EOF
}

# IAM role for user device logout lambda function
resource "aws_iam_role" "clientapi_user_device_logout_lambda" {
  name = "${module.naming.aws_iam_role}-usr-dev-lgo-lambda-${var.api_gw_stage}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_iam_role}-usr-dev-lgo-lambda-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

data "template_file" "clientapi_user_device_logout_lambda" {
  template = "${file("./policies/clientapi-user-device-logout-lambda.json")}"
}

resource "aws_iam_role_policy" "clientapi_user_device_logout_lambda" {
  name   = "${module.naming.aws_iam_role}-c-u-dev-lo-lambda-${var.api_gw_stage}"
  role   = "${aws_iam_role.clientapi_user_device_logout_lambda.id}"
  policy = "${data.template_file.clientapi_user_device_logout_lambda.rendered}"
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "clientapi_user_device_logout_lambda_cw" {
  role       = "${aws_iam_role.clientapi_user_device_logout_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "clientapi_user_device_logout_lambda_vpc" {
  role       = "${aws_iam_role.clientapi_user_device_logout_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "clientapi_user_device_logout_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.clientapi_user_device_logout_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToInternalBankingSQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey",
        "sqs:GetQueue*",
        "sqs:SendMessage*"
      ],
      "Resource": [
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_sqs_queue.segment_analytics.arn}"
      ]
    }
  ]
}
EOF
}

# business document iam role
resource "aws_iam_role" "clientapi_business_document_lambda" {
  name = "${module.naming.aws_iam_role}-bus-doc-lambda-${var.api_gw_stage}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_iam_role}-bus-doc-lambda-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

data "template_file" "clientapi_business_document_lambda" {
  template = "${file("./policies/clientapi-business-document-lambda.json")}"

  vars {
    bucket_name = "${data.aws_s3_bucket.documents.id}"
    kms_arn     = "${data.aws_kms_alias.documents_bucket.target_key_arn}"
  }
}

resource "aws_iam_role_policy" "clientapi_business_document_lambda" {
  name   = "${module.naming.aws_iam_policy}-bus-doc-lambda-${var.api_gw_stage}"
  role   = "${aws_iam_role.clientapi_business_document_lambda.id}"
  policy = "${data.template_file.clientapi_business_document_lambda.rendered}"
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "clientapi_business_document_lambda_cw" {
  role       = "${aws_iam_role.clientapi_business_document_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "clientapi_business_document_lambda_vpc" {
  role       = "${aws_iam_role.clientapi_business_document_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "clientapi_business_document_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.clientapi_business_document_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToInternalBankingSQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey",
        "sqs:GetQueue*",
        "sqs:SendMessage*"
      ],
      "Resource": [
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_sqs_queue.segment_analytics.arn}"
      ]
    }
  ]
}
EOF
}

# business document signurl iam role
resource "aws_iam_role" "clientapi_business_document_sign_url_lambda" {
  name = "${module.naming.aws_iam_role}-bus-doc-srl-lambda-${var.api_gw_stage}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags {
    Application = "${var.application}"
    Component   = "${var.component}"
    Environment = "${var.environment_name}"
    Name        = "${module.naming.aws_iam_role}-bus-doc-srl-lambda-${var.api_gw_stage}"
    Team        = "${var.team}"
  }
}

data "template_file" "clientapi_business_document_sign_url_lambda" {
  template = "${file("./policies/clientapi-business-document-signurl-lambda.json")}"

  vars {
    bucket_name = "${data.aws_s3_bucket.documents.id}"
    kms_arn     = "${data.aws_kms_alias.documents_bucket.target_key_arn}"
  }
}

resource "aws_iam_role_policy" "clientapi_business_document_sign_url_lambda" {
  name   = "${module.naming.aws_iam_policy}-bus-doc-surl-lambda-${var.api_gw_stage}"
  role   = "${aws_iam_role.clientapi_business_document_sign_url_lambda.id}"
  policy = "${data.template_file.clientapi_business_document_sign_url_lambda.rendered}"
}

# Policy to access CW logs
resource "aws_iam_role_policy_attachment" "clientapi_business_document_sign_url_lambda_cw" {
  role       = "${aws_iam_role.clientapi_business_document_sign_url_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Policy to access VPC resources (i.e. ENIs)
resource "aws_iam_role_policy_attachment" "clientapi_business_document_sign_url_lambda_vpc" {
  role       = "${aws_iam_role.clientapi_business_document_sign_url_lambda.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# inline policy to allow access to segment analytics sqs
resource "aws_iam_role_policy" "clientapi_business_document_sgined_url_segment_analytics_sqs" {
  name = "${module.naming.aws_iam_role_policy}-segment-analytics"
  role = "${aws_iam_role.clientapi_business_document_sign_url_lambda.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAccessToInternalBankingSQSQueue",
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt",
        "kms:GenerateDataKey",
        "sqs:GetQueue*",
        "sqs:SendMessage*"
      ],
      "Resource": [
        "${data.aws_kms_alias.internal_sqs.target_key_arn}",
        "${data.aws_sqs_queue.segment_analytics.arn}"
      ]
    }
  ]
}
EOF
}
