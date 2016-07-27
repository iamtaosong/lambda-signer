variable "aws_access_key_id" {
  description = "AWS Access Key ID"
  type        = "string"
}

variable "aws_secret_access_key" {
  description = "AWS Secret Access Key"
  type        = "string"
}

variable "aws_region" {
  description = "AWS Region"
  type        = "string"
}

variable "bucket_name" {
  description = "Name of the bucket to store keys"
  type        = "string"
}

variable "filename" {
  description = "Filename of lambda bundle"
  type        = "string"
}
