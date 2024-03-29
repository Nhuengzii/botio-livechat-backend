# BotioLivechat Backend

The BotioLivechat Service is a powerful platform that brings together multiple chat platforms into a centralized and seamless communication hub. It allows businesses and organizations to streamline their customer support, sales, and engagement processes by integrating various chat channels into a single, user-friendly interface. This service is designed to enhance efficiency, boost customer satisfaction, and simplify chat management for businesses of all sizes.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Configurations](#configurations)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Technologies](#technologies)
- [Other Services!](#other-services)

## Features

- Centralized `Facebook, Instagram, Line` chat services into one communication platform!
- Interact with customers on different platforms with ease by just a change of conversation's tab.
- Customize and save your own template messages .

## Requirements

### MongoDB

Create mongoDB's database with following collections.

- conversations
- messages
- shop_config
- shops
- templates

### Redis

Create Redis's database.

### Discord (Optional)

Can be Used for logging. **Discord logging take some resources and time**, Recommended that it should only be used for error logging.

- Create discord server and text channel.
- add channel's webhook URL to `terraform.tfvars`.

### Facebook pages and Instagram official accounts

- Setup via [Facebook Developer Console](https://developers.facebook.com/).
- For facebook webhook services subscribe to `messages` and `messages_echoes` webhook.
- For instagram webhook services subscribe to `messages` webhook.

### Line bot
- [Getting started with the Messaging API](https://developers.line.biz/en/docs/messaging-api/getting-started/)
- [Building a bot](https://developers.line.biz/en/docs/messaging-api/building-bot/)

### Terraform
- Have `Terraform CLI` (1.2.0+) installed.

### AWS
- Have AWS account
- Create S3 bucket for storing Terraform's state files. This must be created manually.
- Have `AWS CLI` installed.





## Configurations

Editing `terraform.tfvars` allows you to edit your deployment options.

```
aws_access_key           = "aws_access_key" #AWS account's access key. Can be retrieved via AWS's console.
aws_secret_key           = "aws_secret_key" #AWS account's secret access key. Can be retrieved via AWS's console.
rest_api_stage_name      = "dev" #Deployment stage for REST api endpoints.
websocket_api_stage_name = "dev" #Deployment stage for websocket endpoints.
mongo_database           = "mongodb_database" #MongoDB's database name.
discord_webhook_url      = "discord_webhook_url" #Discord text channel's webhook URL for logging.
redis_addr               = "redis_addr" #Redis database address.
redis_password           = "redis_password" #Redis database password.
media_storage_bucket_name           = "botio-lifechat-bucket-name"  #S3 storage bucket name

facebook_app_secret                  = "facebook_app_secret" #Facebook app secret. Can be retrieved via facebook's developer console.
facebook_webhook_verification_string = "facebook_webhook_verification_string" #Set to the same value with facebook's developer console's webhook setup.

instagram_app_secret                  = "instagram_app_secret" #Instagram app secret. Can be retrieved via facebook's developer console.
instagram_webhook_verification_string = "instagram_webhook_verification_string"  #Set to the same value with facebook's developer console's webhook setup.


```

## Usage

All deployment changes will happen on your registered AWS account. **Do not forget to destroy the deployment before deleting your local repository otherwise you will need to clean the deployed services yourself!**
### Setup
- copy `terraform_backend.conf.example` to `terraform_backend.conf`
- edit value of each keys in `terraform_backend.conf`
- initialize Terraform by running command `make init`
- get `aws access key` and `aws secret access key` ready
-  copy `variables.example.tfvars` to `terraform.tfvars`
- edit value of each keys in `terraform.tfvars`
### Deploy

- make sure that values in `terraform.tfvars` are valid
- run command `make deploy`

### Apply

- run `make apply` to apply changes to the Lambdas.

### Destroy

- empty the S3 bucket that stores media files
- run `make destroy` to destroy the system

## API Documentation

View documentation in `openapi_specs/openapi_specs.yaml`.
Visualize the API with [SWAGGER](https://swagger.io/) tool.

## Technologies

### AWS services

- Serverless cloud computing service`Lambda` .
- Storage service `S3`
- Message queue service `SQS`.
- Publish subsciption service `SNS`.
- API proxy servicefor REST API and Websocket connection `API gateway`.

### Database

- `MongoDB` storing service's data.
- `Redis` for caching websocket connection.

### Logging

- AWS's `Cloudwatch`.
- Discord's webhook.

## Other Services!

- [botio](https://www.botio.services) : One platform, endless possiblities.
