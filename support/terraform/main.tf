/**
 * Module usage:
 *
 *   module "signer" {
 *     source "github.com/sthulb/signer-lambda/support/terraform"
 *
 *     aws_access_key_id     = "${var.aws_access_key_id}"
 *     aws_secret_access_key = "${var.aws_secret_access_key}"
 *     aws_region            = "eu-west-1"
 *
 *     function_name = "signer-lve"
 *     filename      = "archive.zip"
 *   }
 */

provider "aws" {
  access_key = "${var.aws_access_key_id}"
  secret_key = "${var.aws_secret_access_key}"
  region     = "${var.aws_region}"
}

resource "aws_kms_key" "key" {
  description             = "Signer key holder"
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "alias" {
  name          = "alias/${var.function_name}-keys"
  target_key_id = "${aws_kms_key.key.key_id}"
}

resource "aws_cloudwatch_event_target" "target" {
  rule = "${aws_cloudwatch_event_rule.rule.name}"
  arn  = "${aws_lambda_function.signer.arn}"
}

resource "aws_cloudwatch_event_rule" "rule" {
  name        = "${var.function_name}"
  description = "DNS Lambda rule"

  event_pattern = <<PATTERN
{
  "source": [
    "aws.autoscaling"
  ],
  "detail-type": [
    "EC2 Instance Launch Successful"
  ]
}
PATTERN
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.signer.arn}"
  principal     = "events.amazonaws.com"
  source_arn    = "${aws_cloudwatch_event_rule.rule.arn}"
}

resource "aws_iam_role" "role" {
  name = "${var.function_name}-${var.aws_region}"

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
}

resource "aws_iam_role_policy" "role" {
  name = "signer_${var.aws_region}"
  role = "${aws_iam_role.role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "*"
      ],
      "Effect": "Allow",
      "Resource": [
         "*"
      ]
    }
  ]
}
EOF
}

resource "aws_lambda_function" "signer" {
  filename         = "${var.filename}"
  function_name    = "${var.function_name}"
  role             = "${aws_iam_role.role.arn}"
  handler          = "index.handle"
  runtime          = "nodejs4.3"
  source_code_hash = "${base64sha256(file(var.filename))}"
}

resource "aws_s3_bucket" "bucket" {
  bucket = "${var.function_name}"
  acl    = "private"
}

resource "template_file" "config" {
  template = <<CONFIG
{
  "kms_key_id": "${kms_key_id}"
}
CONFIG

  vars {
    "kms_key_id" = "${aws_kms_key.key.arn}"
  }
}

resource "aws_s3_bucket_object" "config" {
  bucket  = "${var.function_name}"
  key     = "config.json"
  content = "${template_file.config.rendered}"
}
