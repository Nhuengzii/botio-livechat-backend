provider "aws" {
  region = "ap-southeast-1"
}


resource "aws_api_gateway_stage" "botio_rest_api_test_stage" {
  stage_name    = "test"
  deployment_id = aws_api_gateway_deployment.botio_rest_api_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
}
resource "aws_api_gateway_rest_api" "botio_rest_api" {
  name        = "botio_rest_api"
  description = "Api endpoint for interacting with botio"
}

resource "aws_api_gateway_resource" "shops" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_rest_api.botio_rest_api.root_resource_id
  path_part   = "shops"
}

resource "aws_api_gateway_resource" "shop_id" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.shops.id
  path_part   = "{shop_id}"
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
  name               = "assume_role_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_basic_sqsexecution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}
resource "aws_iam_role_policy_attachment" "lambda_basic_snsexecution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSNSFullAccess"
}

resource "aws_api_gateway_deployment" "botio_rest_api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [aws_api_gateway_method.get_validate_facebook_webhook, aws_api_gateway_integration.get_validate_facebook_webhook]
}

output "botio_invoke_url" {
  value = aws_api_gateway_deployment.botio_rest_api_deployment.invoke_url
}
