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
}

resource "aws_api_gateway_resource" "shops" {
  rest_api_id = module.rest_api.id
  parent_id   = module.rest_api.root_resource_id
  path_part   = "shops"
}

resource "aws_api_gateway_resource" "shop_id" {
  rest_api_id = module.rest_api.id
  parent_id   = aws_api_gateway_resource.shops.id
  path_part   = "{shop_id}"
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
