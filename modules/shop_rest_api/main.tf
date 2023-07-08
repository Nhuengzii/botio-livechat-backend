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


module "handlers" {
  source                = "../lambda_handler"
  for_each              = var.handlers
  handler_name          = each.value.handler_name
  handler_path          = each.value.handler_path
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = each.value.environment_variables
}

locals {
  method_mapping = {
    post_shops = {
      method        = "POST"
      resource_id   = aws_api_gateway_resource.shops.id
      resource_path = aws_api_gateway_resource.shops.path
    }
    get_shop_id = {
      method        = "GET"
      resource_id   = aws_api_gateway_resource.shop_id.id
      resource_path = aws_api_gateway_resource.shop_id.path
    }
    patch_shop_id = {
      method        = "PATCH"
      resource_id   = aws_api_gateway_resource.shop_id.id
      resource_path = aws_api_gateway_resource.shop_id.path
    }
  }
}

module "method_lambda_integration" {
  source                 = "../method_lambda_integration"
  for_each               = var.handlers
  method                 = local.method_mapping[each.key].method
  resource_id            = local.method_mapping[each.key].resource_id
  resource_path          = local.method_mapping[each.key].resource_path
  rest_api_id            = var.rest_api_id
  rest_api_execution_arn = var.rest_api_execution_arn
  lambda_invoke_arn      = module.handlers[each.key].lambda.invoke_arn
  lambda_function_name   = module.handlers[each.key].lambda.function_name
}

