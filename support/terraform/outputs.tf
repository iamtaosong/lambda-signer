// KMS Key ARN
output "kms_key_arn" {
  value = "${aws_kms_key.key.arn}"
}

// Lambda function ARN
output "lambda_func_arn" {
  value = "${aws_lambda_function.signer.arn}"
}
