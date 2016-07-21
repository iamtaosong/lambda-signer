 Module usage:

     module "signer" {
       source "github.com/sthulb/signer-lambda/support/terraform"

       aws_access_key_id     = "${var.aws_access_key_id}"
       aws_secret_access_key = "${var.aws_secret_access_key}"
       aws_region            = "eu-west-1"

       filename = "../../build/archive.zip"
     }


## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| aws_access_key_id | AWS Access Key ID | - | yes |
| aws_secret_access_key | AWS Secret Access Key | - | yes |
| aws_region | AWS Region | - | yes |
| filename | Filename of lambda bundle | - | yes |

## Outputs

| Name | Description |
|------|-------------|
| kms_key_arn | KMS Key ARN |
| lambda_func_arn | Lambda function ARN |

