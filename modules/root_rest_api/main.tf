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
  name               = "assume_role_lambda_for_root_handler"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "s3" {
  statement {
    actions = [
      "s3:*"
    ]
    effect = "Allow"
    resources = [
      "${var.s3_bucket_arn}/*"
    ]
  }
}

resource "aws_iam_policy" "s3" {
  name   = "s3_policy_for_root_handler"
  policy = data.aws_iam_policy_document.s3.json
}

resource "aws_iam_policy_attachment" "s3" {
  name = "s3_policy_attachment_for_root_handler"
  roles = [
    aws_iam_role.assume_role_lambda.name
  ]
  policy_arn = aws_iam_policy.s3.arn
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

module "get_upload_url_handler" {
  source                = "../lambda_handler"
  handler_name          = var.get_upload_url_handler.handler_name
  handler_path          = var.get_upload_url_handler.handler_path
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = merge(var.get_upload_url_handler.environment_variables, {})
  dependencies          = var.get_upload_url_handler.dependencies
}

module "get_upload_url" {
  source                 = "../method_lambda_integration"
  method                 = "GET"
  resource_id            = aws_api_gateway_resource.upload_url.id
  resource_path          = aws_api_gateway_resource.upload_url.path
  rest_api_id            = aws_api_gateway_rest_api.rest_api.id
  rest_api_execution_arn = aws_api_gateway_rest_api.rest_api.execution_arn
  lambda_invoke_arn      = module.get_upload_url_handler.lambda.invoke_arn
  lambda_function_name   = module.get_upload_url_handler.lambda.function_name
}

resource "aws_api_gateway_deployment" "rest_api" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  lifecycle {
    create_before_destroy = true
  }
  triggers = {
    always_run = timestamp()
  }
  depends_on = [module.get_upload_url.method, module.get_upload_url.integration]
}

resource "aws_api_gateway_stage" "dev" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  deployment_id = aws_api_gateway_deployment.rest_api.id
  stage_name    = "dev"
}
