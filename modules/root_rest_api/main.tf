terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

resource "aws_api_gateway_rest_api" "rest_api" {
  name = var.rest_api_name
}

resource "aws_api_gateway_resource" "upload_url" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_rest_api.rest_api.root_resource_id
  path_part   = "upload_url"
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

resource "aws_iam_role_policy_attachment" "lambda_basic_execution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

module "get_upload_url_handler" {
  source       = "../lambda_handler"
  handler_name = "get_upload_url"
  handler_path = format("%s/cmd/lambda/root/get_upload_url", path.root)
  role_arn     = aws_iam_role.assume_role_lambda.arn
}

module "get_upload_url" {
  source                 = "../method_lambda_integration"
  method                 = "GET"
  resource_id            = aws_api_gateway_resource.upload_url.id
  resource_path          = aws_api_gateway_resource.upload_url.path
  rest_api_id            = aws_api_gateway_rest_api.rest_api.id
  rest_api_execution_arn = aws_api_gateway_rest_api.rest_api.execution_arn
  lambda_invoke_arn      = module.get_upload_url_handler.lambda_invoke_arn
  lambda_function_name   = module.get_upload_url_handler.lambda_function_name
}
