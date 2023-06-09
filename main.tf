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

resource "aws_api_gateway_rest_api" "rest_api" {
  name = "botio_rest_api"
}

resource "aws_api_gateway_resource" "shops" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_rest_api.rest_api.root_resource_id
  path_part   = "shops"
}

resource "aws_api_gateway_resource" "shop_id" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_resource.shops.id
  path_part   = "{shop_id}"
}

module "facebook_rest_api" {
  source                 = "./modules/facebook_rest_api"
  platform               = "facebook"
  rest_api_id            = aws_api_gateway_rest_api.rest_api.id
  rest_api_execution_arn = aws_api_gateway_rest_api.rest_api.execution_arn
  parent_id              = aws_api_gateway_resource.shop_id.id
}

output "rest_api" {
  value = {
    id = aws_api_gateway_rest_api.rest_api.id
  }
}
