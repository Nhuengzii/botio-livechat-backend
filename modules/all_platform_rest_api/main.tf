terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
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

resource "aws_api_gateway_resource" "platform" {
  rest_api_id = var.rest_api_id
  parent_id   = var.parent_id
  path_part   = "all"
}

resource "aws_api_gateway_resource" "conversations" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.platform.id
  path_part   = "conversations"
}

resource "aws_iam_role" "assume_role_lambda" {
  name               = format("all_platform_assume_role_lambda")
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

module "get_conversations_handler" {
  source                = "../lambda_handler"
  handler_name          = var.get_conversations_handler.handler_name
  handler_path          = var.get_conversations_handler.handler_path
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = var.get_conversations_handler.environment_variables
  dependencies          = var.get_conversations_handler.dependencies
}

resource "aws_api_gateway_method" "get_conversations" {
  rest_api_id   = var.rest_api_id
  resource_id   = aws_api_gateway_resource.conversations.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "get_conversations" {
  http_method             = aws_api_gateway_method.get_conversations.http_method
  resource_id             = aws_api_gateway_resource.conversations.id
  rest_api_id             = var.rest_api_id
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = module.get_conversations_handler.lambda.invoke_arn
}

resource "aws_lambda_permission" "endpoint_handler_permissions" {
  function_name = module.get_conversations_handler.lambda.function_name
  statement_id  = "AllowMethodExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = format("%s/*/GET/%s", var.rest_api_execution_arn, aws_api_gateway_resource.conversations.path)
}
