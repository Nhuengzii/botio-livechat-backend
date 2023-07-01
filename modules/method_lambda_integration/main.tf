terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

resource "aws_api_gateway_method" "method" {
  http_method   = var.method
  resource_id   = var.resource_id
  rest_api_id   = var.rest_api_id
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "integration" {
  http_method             = aws_api_gateway_method.method.http_method
  integration_http_method = "POST"
  rest_api_id             = aws_api_gateway_method.method.rest_api_id
  resource_id             = aws_api_gateway_method.method.resource_id
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

resource "aws_lambda_permission" "endpoint_handler_permission" {
  function_name = var.lambda_function_name
  statement_id  = format("AllowMethod_%s_ExecutionFromAPIGateway", var.method)
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = format("%s/*/%s%s", var.rest_api_execution_arn, var.method, var.resource_path)
}

