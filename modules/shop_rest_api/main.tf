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

module "shop_id_enable_cors" {
  source          = "squidfunk/api-gateway-enable-cors/aws"
  version         = "0.3.3"
  api_id          = var.rest_api_id
  api_resource_id = aws_api_gateway_resource.shop_id.id
}

data "aws_iam_policy_document" "assume_role_lambda" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "assume_role_lambda" {
  name               = "assume_role_lambda_for_shops_handlers"
  assume_role_policy = data.aws_iam_policy_document.assume_role_lambda.json
}

module "post_shop_id" {
  source                = "../lambda_handler"
  handler_name          = var.handlers["post_shop_id"].handler_name
  handler_path          = var.handlers["post_shop_id"].handler_path
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = var.handlers["post_shop_id"].environment_variables
}

module "get_shop_id" {
  source                = "../lambda_handler"
  handler_name          = var.handlers["get_shop_id"].handler_name
  handler_path          = var.handlers["get_shop_id"].handler_path
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = var.handlers["get_shop_id"].environment_variables
}

resource "aws_api_gateway_method" "post_shop_id" {
  http_method   = "POST"
  rest_api_id   = var.rest_api_id
  resource_id   = aws_api_gateway_resource.shop_id.id
  authorization = "NONE"
}

resource "aws_api_gateway_method" "get_shop_id" {
  http_method   = "GET"
  rest_api_id   = var.rest_api_id
  resource_id   = aws_api_gateway_resource.shop_id.id
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "post_shop_id" {
  http_method             = aws_api_gateway_method.post_shop_id.http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_resource.shop_id.id
  rest_api_id             = var.rest_api_id
  type                    = "AWS_PROXY"
  uri                     = module.post_shop_id.lambda.invoke_arn
}

resource "aws_api_gateway_integration" "get_shop_id" {
  http_method             = aws_api_gateway_method.get_shop_id.http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_resource.shop_id.id
  rest_api_id             = var.rest_api_id
  type                    = "AWS_PROXY"
  uri                     = module.get_shop_id.lambda.invoke_arn
}

resource "aws_lambda_permission" "allow_api_gateway_to_invoke_post_shop_id" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.post_shop_id.lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = format("%s/*/%s%s", var.rest_api_execution_arn, "POST", aws_api_gateway_resource.shop_id.path)
}

resource "aws_lambda_permission" "allow_api_gateway_to_invoke_get_shop_id" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.get_shop_id.lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = format("%s/*/%s%s", var.rest_api_execution_arn, "POST", aws_api_gateway_resource.shop_id.path)
}
