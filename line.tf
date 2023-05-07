resource "aws_api_gateway_resource" "line" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.shop_id.id
  path_part   = "line"
}

resource "aws_api_gateway_resource" "line_page_id" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.line.id
  path_part   = "{page_id}"
}

resource "aws_api_gateway_resource" "line_webhook" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.line_page_id.id
  path_part   = "webhook"
}

resource "aws_api_gateway_method" "post_validate_line_webhook" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.line_webhook.id
  authorization = "NONE"
  http_method   = "POST"
}

resource "aws_api_gateway_integration" "post_validate_line_webhook" {
  http_method             = aws_api_gateway_method.post_validate_line_webhook.http_method
  resource_id             = aws_api_gateway_resource.line_webhook.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.validate_line_webhook_handler.invoke_arn
}
resource "aws_lambda_permission" "validate_line_webhook_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.validate_line_webhook_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "aws_sqs_queue" "line_webhook_to_standardize_line_webhook_handler" {
  name = "line_webhook_to_standardize_facebook_webhook_handler"
}

resource "aws_lambda_event_source_mapping" "event_source_mapping_line_webhook_to_standardize_line_webhook_handler" {
  event_source_arn = aws_sqs_queue.line_webhook_to_standardize_line_webhook_handler.arn
  function_name    = aws_lambda_function.standardize_line_webhook_handler.function_name
  batch_size       = 1
}

resource "null_resource" "build_validate_line_webhook_handler" {
  triggers = {
    source_code_hash = filebase64sha256("validate_line_webhook_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./validate_line_webhook_handler/src/ -o ../bin/main ."
  }
}
resource "null_resource" "build_standardize_line_webhook_handler" {
  triggers = {
    source_code_hash = filebase64sha256("standardize_line_webhook_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./standardize_line_webhook_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "watch_validate_line_webhook_handler" {
  triggers = {
    validate_line_webhook_handler = aws_lambda_function.validate_line_webhook_handler.qualified_arn
  }
  depends_on = [null_resource.build_validate_line_webhook_handler]
}
resource "null_resource" "watch_standardize_line_webhook_handler" {
  triggers = {
    standardize_line_webhook_handler = aws_lambda_function.standardize_line_webhook_handler.qualified_arn
  }
  depends_on = [null_resource.build_validate_line_webhook_handler]
}

data "archive_file" "validate_line_webhook_handler" {
  type        = "zip"
  source_dir  = "validate_line_webhook_handler"
  output_path = "validate_line_webhook_handler/validate_line_webhook_handler.zip"
}

data "archive_file" "standardize_line_webhook_handler" {
  type        = "zip"
  source_dir  = "standardize_line_webhook_handler"
  output_path = "standardize_line_webhook_handler/standardize_line_webhook_handler.zip"
}

resource "aws_lambda_function" "validate_line_webhook_handler" {
  filename         = data.archive_file.validate_line_webhook_handler.output_path
  function_name    = "validate_line_webhook_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = filebase64sha256("validate_line_webhook_handler/src/main.go")
  depends_on       = [data.archive_file.validate_line_webhook_handler]
  environment {
    variables = {
      SQS_QUEUE_URL = aws_sqs_queue.line_webhook_to_standardize_line_webhook_handler.id
      SQS_QUEUE_ARN = aws_sqs_queue.line_webhook_to_standardize_line_webhook_handler.arn
      foo           = "bar"
    }
  }
}

resource "aws_lambda_function" "standardize_line_webhook_handler" {
  filename      = data.archive_file.standardize_line_webhook_handler.output_path
  function_name = "standardize_line_webhook_handler"
  role          = aws_iam_role.assume_role_lambda.arn
  handler       = "main"
  runtime       = "go1.x"
  environment {
    variables = {
      SQS_QUEUE_URL = aws_sqs_queue.line_webhook_to_standardize_line_webhook_handler.id
      SQS_QUEUE_ARN = aws_sqs_queue.line_webhook_to_standardize_line_webhook_handler.arn
      foo           = "bar"
    }
  }
  depends_on = [data.archive_file.standardize_line_webhook_handler]
}
