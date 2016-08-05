 Module usage:

   module "signer" {
     source "github.com/sthulb/signer-lambda/support/terraform"

     aws_access_key_id     = "${var.aws_access_key_id}"
     aws_secret_access_key = "${var.aws_secret_access_key}"
     aws_region            = "eu-west-1"

     ca_cert       = "ca.pem"
     function_name = "signer-lve"
     filename      = "archive.zip"
   }


## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| aws_access_key_id | AWS Access Key ID | - | yes |
| aws_secret_access_key | AWS Secret Access Key | - | yes |
| aws_region | AWS Region | - | yes |
| ca_cert | Filename of CA cert to sign things with | - | yes |
| filename | Filename of lambda bundle | - | yes |
| function_name | Name of lambda bundle | - | yes |
| vpc_id | VPC ID | - | yes |

## Outputs

| Name | Description |
|------|-------------|
| kms_key_arn | KMS Key ARN |
| lambda_func_arn | Lambda function ARN |

