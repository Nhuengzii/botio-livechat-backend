resource "aws_api_gateway_resource" "facebook" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.shop_id.id
  path_part   = "facebook"
}

resource "aws_api_gateway_resource" "facebook_page_id" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.facebook.id
  path_part   = "{page_id}"
}

resource "aws_api_gateway_resource" "facebook_webhook" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.facebook_page_id.id
  path_part   = "webhook"
}

resource "aws_api_gateway_method" "validate_facebook_webhook" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.facebook_webhook.id
  authorization = "NONE"
  http_method   = "GET"
}

resource "aws_api_gateway_integration" "get_validate_facebook_webhook" {
  http_method             = aws_api_gateway_method.validate_facebook_webhook.http_method
  resource_id             = aws_api_gateway_resource.facebook_webhook.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.validate_facebook_webhook_handler.invoke_arn
}

resource "aws_lambda_function" "validate_facebook_webhook_handler" {
  filename         = "validate_facebook_webhook_handler/validate_facebook_webhook_handler.zip"
  function_name    = "validate_facebook_webhook_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = filebase64sha256("validate_facebook_webhook_handler/src/main.go")
  depends_on       = [data.archive_file.validate_facebook_webhook_handler]
}

resource "aws_lambda_function" "standardize_facebook_webhook_handler" {
  filename         = "standardize_facebook_webhook_handler/standardize_facebook_webhook_handler.zip"
  function_name    = "standardize_facebook_webhook_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = filebase64sha256("standardize_facebook_webhook_handler/src/main.go")
  depends_on       = [data.archive_file.standardize_facebook_webhook_handler]
}

resource "aws_lambda_permission" "validate_facebook_webhook_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.validate_facebook_webhook_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "null_resource" "build_validate_facebook_webhook_handler" {
  triggers = {
    source_code_hash = "${filebase64sha256("validate_facebook_webhook_handler/src/main.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./validate_facebook_webhook_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_standardize_facebook_webhook_handler" {
  triggers = {
    source_code_hash = "${filebase64sha256("standardize_facebook_webhook_handler/src/main.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./standardize_facebook_webhook_handler/src/ -o ../bin/main ."
  }
}

data "archive_file" "validate_facebook_webhook_handler" {
  type        = "zip"
  source_file = "./validate_facebook_webhook_handler/bin/main"
  output_path = "./validate_facebook_webhook_handler/validate_facebook_webhook_handler.zip"
  depends_on  = [null_resource.build_validate_facebook_webhook_handler]
}

data "archive_file" "standardize_facebook_webhook_handler" {
  type        = "zip"
  source_file = "./standardize_facebook_webhook_handler/src/main.go"
  output_path = "./standardize_facebook_webhook_handler/standardize_facebook_webhook_handler.zip"
  depends_on  = [null_resource.build_standardize_facebook_webhook_handler]
}



resource "null_resource" "watch_validate_facebook_webhook_handler" {
  triggers = {
    source_code_hash = filebase64sha256("validate_facebook_webhook_handler/src/main.go")
  }
  depends_on = [null_resource.build_validate_facebook_webhook_handler]
}

resource "null_resource" "watch_standardize_facebook_webhook_handler" {
  triggers = {
    source_code_hash = filebase64sha256("standardize_facebook_webhook_handler/src/main.go")
  }
  depends_on = [null_resource.build_standardize_facebook_webhook_handler]
}
