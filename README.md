# Lambda Signer

Auto signs certificates for EC2 instances. Credit to [Steven Jack](https://github.com/stevenjack)
for suggesting it.

## Design
The application is shipped using AWS Lambda. It uses KMS to store the certificate files encrypted at rest. This probably isn't the best way to
store the certificates.

## Build
The project can be built with `go build -o main .`

## Configuration
The binary should be shipped with a `config.json`

| Key                | Description                                |
|:-------------------|:-------------------------------------------|
| `bucket`           | Bucket name to store certs in              |
| `environment_name` | Used as the organisation name in the cert  |
| `kms_key_id`       | ARN of the KMS key                         |

```json
{
  "bucket": "bucket_name",
  "environment_name": "name",
  "kms_key_id": "arn:aws:kms:eu-west-1:0123456789:key/1234-1234-1234"
}
```

## Deploy
A [terraform module](support/terraform/README.md) has been created that will help you manage deployments of this application.
