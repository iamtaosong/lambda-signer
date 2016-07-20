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
  name          = "alias/signer-keys"
  target_key_id = "${aws_kms_key.key.key_id}"
}
