data "archive_file" "lambda_hsts" {
  type        = "zip"
  output_path = "./modules/hsts_lambda/lambdas/hsts/lambda.zip"

  source {
    filename = "index.js"
    content  = "${file("./modules/hsts_lambda/lambdas/hsts/lambda.js")}"
  }
}

resource "aws_lambda_function" "hsts" {
  function_name    = "${module.naming.aws_lambda_function}-hsts"
  filename         = "${data.archive_file.lambda_hsts.output_path}"
  source_code_hash = "${data.archive_file.lambda_hsts.output_base64sha256}"
  role             = "${aws_iam_role.lambda_hsts.arn}"
  runtime          = "nodejs10.x"
  handler          = "index.handler"
  memory_size      = 128
  timeout          = 3
  publish          = true

  tags {
    Application = "${var.application}"
    Environment = "${terraform.workspace}"
    Component   = "${var.component}"
    Name        = "${module.naming.aws_lambda_function}-hsts"
    Team        = "${var.team}"
  }

  provider = "aws.${var.provider_name}"
}
