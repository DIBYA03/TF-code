data "aws_iam_policy_document" "lambda_hsts" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type = "Service"

      identifiers = [
        "lambda.amazonaws.com",
        "edgelambda.amazonaws.com",
      ]
    }
  }

  provider = "aws.${var.provider_name}"
}

resource "aws_iam_role" "lambda_hsts" {
  name               = "${module.naming.aws_iam_role}-hsts"
  assume_role_policy = "${data.aws_iam_policy_document.lambda_hsts.json}"

  provider = "aws.${var.provider_name}"
}

resource "aws_iam_role_policy_attachment" "lambda_hsts_cw_logs" {
  role       = "${aws_iam_role.lambda_hsts.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

  provider = "aws.${var.provider_name}"
}
