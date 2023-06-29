terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

resource "aws_api_gateway_resource" "shops" {
  rest_api_id = var.rest_api_id
  parent_id   = var.parent_resource_id
  path_part   = "shops"
}

resource "aws_api_gateway_resource" "shop_id" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.shops.id
  path_part   = "{shop_id}"
}
