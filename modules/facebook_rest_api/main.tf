terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

variable "platform" {
  type = string
}

variable "rest_api_id" {
  type = string
}

variable "rest_api_execution_arn" {
  type = string
}

variable "facebook_access_token" {
  type = string
}

variable "facebook_app_secret" {
  type = string
}

variable "parent_id" {
  type = string
}
data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "assume_role_lambda" {
  name               = format("%s_assume_role_lambda", var.platform)
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_api_gateway_resource" "platform" {
  rest_api_id = var.rest_api_id
  parent_id   = var.parent_id
  path_part   = var.platform
}

resource "aws_api_gateway_resource" "page_id" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.platform.id
  path_part   = "{page_id}"
}

resource "aws_api_gateway_resource" "webhook" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.page_id.id
  path_part   = "webhook"
}

resource "aws_api_gateway_resource" "conversations" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.page_id.id
  path_part   = "conversations"
}

resource "aws_api_gateway_resource" "conversation_id" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.conversations.id
  path_part   = "{conversation_id}"
}

resource "aws_api_gateway_resource" "messages" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.conversation_id.id
  path_part   = "messages"
}

resource "aws_api_gateway_method" "get_post_webhook" {
  for_each      = toset(["GET", "POST"])
  http_method   = each.key
  rest_api_id   = var.rest_api_id
  resource_id   = aws_api_gateway_resource.webhook.id
  authorization = "NONE"
}

module "get_post_webhook_handler" {
  source       = "../lambda_handler/"
  handler_name = format("%s_get_post_webhook_handler", var.platform)
  handler_path = format("%s/validate_%s_webhook_handler", path.root, var.platform)
  role_arn     = aws_iam_role.assume_role_lambda.arn
}

resource "aws_lambda_permission" "get_post_webhook" {
  statement_id  = format("AllowMethod_%s_ExecutionFromAPIGateway", each.key)
  action        = "lambda:InvokeFunction"
  function_name = module.get_post_webhook_handler.lambda.function_name
  principal     = "apigateway.amazonaws.com"
  for_each      = toset(["GET", "POST"])
  source_arn    = format("%s/*/%s%s", var.rest_api_execution_arn, each.key, aws_api_gateway_resource.webhook.path)
}

resource "aws_api_gateway_integration" "get_post_webhook" {
  for_each                = toset(["GET", "POST"])
  http_method             = aws_api_gateway_method.get_post_webhook[each.key].http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_resource.webhook.id
  rest_api_id             = var.rest_api_id
  type                    = "AWS_PROXY"
  uri                     = module.get_post_webhook_handler.lambda.invoke_arn
}

locals {
  endpoint_with_handlers = {
    get_message = {
      method        = "GET"
      resource_id   = aws_api_gateway_resource.messages.id
      resource_path = aws_api_gateway_resource.messages.path
      handler_name  = format("%s_get_messages_handler", var.platform)
      handler_path  = format("%s/get_%s_messages_handler", path.root, var.platform)
      role_arn      = aws_iam_role.assume_role_lambda.arn
      environment_variables = {
        ACCESS_TOKEN = var.facebook_access_token
      }
    }
    post_message = {
      method                = "POST"
      resource_id           = aws_api_gateway_resource.messages.id
      resource_path         = aws_api_gateway_resource.messages.path
      handler_name          = format("%s_post_message_handler", var.platform)
      handler_path          = format("%s/post_%s_message_handler", path.root, var.platform)
      role_arn              = aws_iam_role.assume_role_lambda.arn
      environment_variables = {}
    }
    get_conversations = {
      method        = "GET"
      resource_id   = aws_api_gateway_resource.conversations.id
      resource_path = aws_api_gateway_resource.conversations.path
      handler_name  = format("%s_get_conversations_handler", var.platform)
      handler_path  = format("%s/get_%s_conversation_handler", path.root, var.platform)
      role_arn      = aws_iam_role.assume_role_lambda.arn
      environment_variables = {
        ACCESS_TOKEN = var.facebook_access_token
      }
    }
  }
}

resource "aws_api_gateway_method" "endpoint_method" {
  for_each      = local.endpoint_with_handlers
  http_method   = each.value.method
  rest_api_id   = var.rest_api_id
  resource_id   = each.value.resource_id
  authorization = "NONE"
}

module "endpoint_handlers" {
  for_each              = local.endpoint_with_handlers
  source                = "../lambda_handler/"
  handler_name          = each.value.handler_name
  handler_path          = each.value.handler_path
  role_arn              = each.value.role_arn
  environment_variables = each.value.environment_variables
}

resource "aws_api_gateway_integration" "endpoint_handler_integrations" {
  for_each                = local.endpoint_with_handlers
  http_method             = aws_api_gateway_method.endpoint_method[each.key].http_method
  integration_http_method = "POST"
  resource_id             = each.value.resource_id
  rest_api_id             = var.rest_api_id
  type                    = "AWS_PROXY"
  uri                     = module.endpoint_handlers[each.key].lambda.invoke_arn
}

resource "aws_lambda_permission" "endpoint_handler_permissions" {
  for_each      = local.endpoint_with_handlers
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.endpoint_handlers[each.key].lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = format("%s/*/%s%s", var.rest_api_execution_arn, each.value.method, each.value.resource_path)
}
