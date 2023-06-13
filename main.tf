terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

variable "line_channel_secret" {
  type = string
}

variable "line_channel_access_token" {
  type = string
}

variable "discord_webhook_url" {
  type = string
}

variable "mongo_uri" {
  type = string
}

variable "mongo_database" {
  type = string
}

variable "mongo_collection_line_conversations" {
  type = string
}

variable "mongo_collection_line_messages" {
  type = string
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
  source                               = "./modules/facebook_rest_api"
  platform                             = "facebook"
  rest_api_id                          = aws_api_gateway_rest_api.rest_api.id
  rest_api_execution_arn               = aws_api_gateway_rest_api.rest_api.execution_arn
  parent_id                            = aws_api_gateway_resource.shop_id.id
  facebook_access_token                = var.facebook_access_token
  facebook_app_secret                  = var.facebook_app_secret
  facebook_webhook_verification_string = var.facebook_webhook_verification_string
  mongo_uri                            = var.mongo_uri
  mongo_database                       = var.mongo_database
  discord_webhook_url                  = var.discord_webhook_url
}

module "line_rest_api" {
  source                              = "./modules/line_rest_api"
  platform                            = "line"
  rest_api_id                         = aws_api_gateway_rest_api.rest_api.id
  rest_api_execution_arn              = aws_api_gateway_rest_api.rest_api.execution_arn
  parent_id                           = aws_api_gateway_resource.shop_id.id
  line_channel_access_token           = var.line_channel_access_token
  line_channel_secret                 = var.line_channel_secret
  mongo_uri                           = var.mongo_uri
  mongo_database                      = var.mongo_database
  mongo_collection_line_messages      = var.mongo_collection_line_messages
  mongo_collection_line_conversations = var.mongo_collection_line_conversations
  discord_webhook_url                 = var.discord_webhook_url
}

resource "aws_apigatewayv2_api" "botio_livechat_websocket" {
  name                       = "botio_livechat_websocket"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "botio_livechat_websocket_dev" {
  api_id      = aws_apigatewayv2_api.botio_livechat_websocket.id
  name        = "dev"
  auto_deploy = true
}

resource "aws_apigatewayv2_deployment" "botio_livechat_websocket_dev" {
  api_id      = aws_apigatewayv2_api.botio_livechat_websocket.id
  description = "dev"
  triggers = {
    always_run = timestamp()
  }
  lifecycle {
    create_before_destroy = true
  }
}

module "websocket_api" {
  source                      = "./modules/websocket_api"
  websocket_api_id            = aws_apigatewayv2_api.botio_livechat_websocket.id
  websocket_api_execution_arn = aws_apigatewayv2_api.botio_livechat_websocket.execution_arn
  redis_password              = var.redis_password
  redis_addr                  = var.redis_addr
}

resource "aws_api_gateway_deployment" "rest_api" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  lifecycle {
    create_before_destroy = true
  }
  triggers = {
    always_run = timestamp()
  }
  depends_on = [module.websocket_api]
}

resource "aws_api_gateway_stage" "dev" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  deployment_id = aws_api_gateway_deployment.rest_api.id
  stage_name    = "dev"
}

output "rest_api" {
  value = {
    id = aws_api_gateway_rest_api.rest_api.id
  }
}
