terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

provider "aws" {
  region = "ap-southeast-1"
}

module "rest_api" {
  source        = "./modules/root_rest_api"
  rest_api_name = "botio_livechat_rest_api"
  s3_bucket_arn = module.bucket.bucket_arn
  get_upload_url_handler = {
    handler_name = "get_upload_url"
    handler_path = format("%s/cmd/lambda/root/get_upload_url", path.root)
    dependencies = ""
    environment_variables = {
      DISCORD_WEBHOOK_URL = var.discord_webhook_url
      S3_BUCKET_NAME      = module.bucket.bucket_name
    }
  }
}

module "shops" {
  source                 = "./modules/shop_rest_api"
  rest_api_id            = module.rest_api.id
  rest_api_execution_arn = module.rest_api.execution_arn
  parent_resource_id     = module.rest_api.root_resource_id
  handlers = {
    post_shops = {
      handler_name = "post_shops"
      handler_path = format("%s/cmd/lambda/shops/post_shops", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
    }
    get_shop_id = {
      handler_name = "get_shop_id"
      handler_path = format("%s/cmd/lambda/shops/get_shop_id", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
    }
  }
}

module "bucket" {
  source      = "./modules/bucket"
  bucket_name = "botio-livechat-bucket"
}

resource "aws_api_gateway_deployment" "rest_api" {
  rest_api_id = module.rest_api.id
  lifecycle {
    create_before_destroy = true
  }
  triggers = {
    always_run = timestamp()
  }
  # depends_on = [module.websocket_api]
}

resource "aws_api_gateway_stage" "dev" {
  rest_api_id   = module.rest_api.id
  deployment_id = aws_api_gateway_deployment.rest_api.id
  stage_name    = "dev"
}

output "rest_api" {
  value = {
    id = module.rest_api.id
  }
}
