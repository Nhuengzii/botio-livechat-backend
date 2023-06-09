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
  handler_path = format("%s/validate_facebook_webhook_handler", path.root)
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
