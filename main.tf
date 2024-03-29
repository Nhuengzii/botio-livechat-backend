terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
  backend "s3" {
  }
}

provider "aws" {
  region     = "ap-southeast-1"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

module "rest_api" {
  source              = "./modules/root_rest_api"
  rest_api_name       = "botio_livechat_rest_api"
  s3_bucket_arn       = module.bucket.bucket_arn
  rest_api_stage_name = var.rest_api_stage_name
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
      dependencies = "{discord,db,apigateway,api,shops}/**/*.go"
    }
    get_shop_id = {
      handler_name = "get_shop_id"
      handler_path = format("%s/cmd/lambda/shops/get_shop_id", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,apigateway,api,shops}/**/*.go"
    }
    patch_shop_id = {
      handler_name = "patch_shop_id"
      handler_path = format("%s/cmd/lambda/shops/patch_shop_id", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,apigateway,api,shops}/**/*.go"
    }
    get_config = {
      handler_name = "get_config"
      handler_path = format("%s/cmd/lambda/shops/get_config", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,apigateway,api}/**/*.go"
    }
    delete_template = {
      handler_name = "delete_template"
      handler_path = format("%s/cmd/lambda/shops/delete_template", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,apigateway,api,templates}/**/*.go"
    }
    get_templates = {
      handler_name = "get_templates"
      handler_path = format("%s/cmd/lambda/shops/get_templates", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,apigateway,api,templates}/**/*.go"
    }
    post_templates = {
      handler_name = "post_templates"
      handler_path = format("%s/cmd/lambda/shops/post_templates", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,apigateway,api,templates}/**/*.go"
    }
  }
}

module "bucket" {
  source      = "./modules/bucket"
  bucket_name = var.media_storage_bucket_name
}

output "rest_api" {
  value = {
    id = module.rest_api.id
    base_url = format("https://%s.execute-api.ap-southeast-1.amazonaws.com/%s", module.rest_api.id, var.rest_api_stage_name)
  }
}
