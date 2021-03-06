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

variable "ca_cert" {
  description = "Filename of CA cert to sign things with"
  type        = "string"
}

variable "filename" {
  description = "Filename of lambda bundle"
  type        = "string"
}

variable "function_name" {
  description = "Name of lambda bundle"
  type        = "string"
}

variable "vpc_id" {
  description = "VPC ID"
  type        = "string"
}
