# BotioLivechat Backend

### - Prerequisites

- The `Terraform CLI` (1.2.0+) installed.
- The `AWS CLI` installed.
- `aws access key` and `aws aws secret key`
- copy `variables.example.tfvars` to `terraform.tfvars`
- edit value of each keys in `terraform.tfvars`

### - Deploy

- initialize terraform by running command `make init`
- make sure that value in terraform.tfvars are valid
- run command `make deploy`

### - Destroy

- empty s3 bucket that use to store image (s3_bucket_name key in `terraform.tfvars`)
- run `make destroy` to destroy the system
