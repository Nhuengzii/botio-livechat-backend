resource "aws_apigatewayv2_api" "botio_livechat_websocket" {
  name                       = "botio_livechat_websocket"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "botio_livechat_websocket_test" {
  api_id      = aws_apigatewayv2_api.botio_livechat_websocket.id
  name        = "test"
  auto_deploy = true
}

resource "aws_apigatewayv2_deployment" "botio_livechat_websocket_test" {
  api_id      = aws_apigatewayv2_api.botio_livechat_websocket.id
  description = "test"
  depends_on  = [aws_apigatewayv2_route.botio_livechat_websocket_default]
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_route" "botio_livechat_websocket_default" {
  api_id    = aws_apigatewayv2_api.botio_livechat_websocket.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.botio_livechat_websocket_default.id}"
}

resource "aws_apigatewayv2_route" "botio_livechat_websocket_connect" {
  api_id    = aws_apigatewayv2_api.botio_livechat_websocket.id
  route_key = "$connect"
  target    = "integrations/${aws_apigatewayv2_integration.botio_livechat_websocket_connect.id}"
}

resource "aws_apigatewayv2_route" "botio_livechat_websocket_disconnect" {
  api_id    = aws_apigatewayv2_api.botio_livechat_websocket.id
  route_key = "$disconnect"
  target    = "integrations/${aws_apigatewayv2_integration.botio_livechat_websocket_disconnect.id}"
}

resource "aws_apigatewayv2_route" "botio_livechat_websocket_broadcast" {
  api_id    = aws_apigatewayv2_api.botio_livechat_websocket.id
  route_key = "broadcast"
  target    = "integrations/${aws_apigatewayv2_integration.botio_livechat_websocket_broadcast.id}"
}
resource "aws_apigatewayv2_route" "botio_livechat_websocket_typing_broadcast" {
  api_id    = aws_apigatewayv2_api.botio_livechat_websocket.id
  route_key = "typing_broadcast"
  target    = "integrations/${aws_apigatewayv2_integration.botio_livechat_websocket_typing_broadcast.id}"
}

resource "aws_apigatewayv2_integration" "botio_livechat_websocket_default" {
  api_id                    = aws_apigatewayv2_api.botio_livechat_websocket.id
  integration_type          = "AWS_PROXY"
  integration_uri           = aws_lambda_function.botio_livechat_websocket_default_handler.invoke_arn
  integration_method        = "POST"
  content_handling_strategy = "CONVERT_TO_TEXT"
  passthrough_behavior      = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_integration" "botio_livechat_websocket_connect" {
  api_id                    = aws_apigatewayv2_api.botio_livechat_websocket.id
  integration_type          = "AWS_PROXY"
  integration_uri           = aws_lambda_function.botio_livechat_websocket_connect_handler.invoke_arn
  integration_method        = "POST"
  content_handling_strategy = "CONVERT_TO_TEXT"
  passthrough_behavior      = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_integration" "botio_livechat_websocket_disconnect" {
  api_id                    = aws_apigatewayv2_api.botio_livechat_websocket.id
  integration_type          = "AWS_PROXY"
  integration_uri           = aws_lambda_function.botio_livechat_websocket_disconnect_handler.invoke_arn
  integration_method        = "POST"
  content_handling_strategy = "CONVERT_TO_TEXT"
  passthrough_behavior      = "WHEN_NO_MATCH"
}
resource "aws_apigatewayv2_integration" "botio_livechat_websocket_broadcast" {
  api_id                    = aws_apigatewayv2_api.botio_livechat_websocket.id
  integration_type          = "AWS_PROXY"
  integration_uri           = aws_lambda_function.botio_livechat_websocket_broadcast_handler.invoke_arn
  integration_method        = "POST"
  content_handling_strategy = "CONVERT_TO_TEXT"
  passthrough_behavior      = "WHEN_NO_MATCH"
}
resource "aws_apigatewayv2_integration" "botio_livechat_websocket_typing_broadcast" {
  api_id                    = aws_apigatewayv2_api.botio_livechat_websocket.id
  integration_type          = "AWS_PROXY"
  integration_uri           = aws_lambda_function.botio_livechat_websocket_typing_broadcast_handler.invoke_arn
  integration_method        = "POST"
  content_handling_strategy = "CONVERT_TO_TEXT"
  passthrough_behavior      = "WHEN_NO_MATCH"
}


variable "redis_access" {
  type = object({
    addr     = string
    password = string
  })
}

resource "aws_lambda_function" "botio_livechat_websocket_default_handler" {
  function_name    = "botio_livechat_websocket_default_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  filename         = data.archive_file.botio_livechat_websocket_default_handler.output_path
  source_code_hash = data.archive_file.botio_livechat_websocket_default_handler.output_base64sha256
  depends_on       = [data.archive_file.botio_livechat_websocket_default_handler]
  environment {
    variables = {
      WEBSOCKET_API_ENDPOINT = "https://${aws_apigatewayv2_api.botio_livechat_websocket.id}.execute-api.ap-southeast-1.amazonaws.com/test"
      REDIS_ACCESS_ADDR      = var.redis_access.addr
      REDIS_ACCESS_PASSWORD  = var.redis_access.password
    }
  }
}
resource "aws_lambda_function" "botio_livechat_websocket_connect_handler" {
  function_name    = "botio_livechat_websocket_connect_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  filename         = data.archive_file.botio_livechat_websocket_connect_handler.output_path
  source_code_hash = data.archive_file.botio_livechat_websocket_connect_handler.output_base64sha256
  depends_on       = [data.archive_file.botio_livechat_websocket_connect_handler]
  environment {
    variables = {
      REDIS_ACCESS_ADDR     = var.redis_access.addr
      REDIS_ACCESS_PASSWORD = var.redis_access.password
    }
  }
}

resource "aws_lambda_function" "botio_livechat_websocket_disconnect_handler" {
  function_name    = "botio_livechat_websocket_disconnect_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  filename         = data.archive_file.botio_livechat_websocket_disconnect_handler.output_path
  source_code_hash = data.archive_file.botio_livechat_websocket_disconnect_handler.output_base64sha256
  depends_on       = [data.archive_file.botio_livechat_websocket_disconnect_handler]
  environment {
    variables = {
      REDIS_ACCESS_ADDR     = var.redis_access.addr
      REDIS_ACCESS_PASSWORD = var.redis_access.password
    }
  }
}

resource "aws_lambda_function" "botio_livechat_websocket_broadcast_handler" {
  function_name    = "botio_livechat_websocket_broadcast_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  filename         = data.archive_file.botio_livechat_websocket_broadcast_handler.output_path
  source_code_hash = data.archive_file.botio_livechat_websocket_broadcast_handler.output_base64sha256
  depends_on       = [data.archive_file.botio_livechat_websocket_broadcast_handler]
  environment {
    variables = {
      REDIS_ACCESS_ADDR      = var.redis_access.addr
      REDIS_ACCESS_PASSWORD  = var.redis_access.password
      WEBSOCKET_API_ENDPOINT = "https://${aws_apigatewayv2_api.botio_livechat_websocket.id}.execute-api.ap-southeast-1.amazonaws.com/test"
    }
  }
}
resource "aws_lambda_function" "botio_livechat_websocket_typing_broadcast_handler" {
  function_name    = "botio_livechat_websocket_typing_broadcast_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  filename         = data.archive_file.botio_livechat_websocket_typing_broadcast_handler.output_path
  source_code_hash = data.archive_file.botio_livechat_websocket_typing_broadcast_handler.output_base64sha256
  depends_on       = [data.archive_file.botio_livechat_websocket_typing_broadcast_handler]
  environment {
    variables = {
      REDIS_ACCESS_ADDR      = var.redis_access.addr
      REDIS_ACCESS_PASSWORD  = var.redis_access.password
      WEBSOCKET_API_ENDPOINT = "https://${aws_apigatewayv2_api.botio_livechat_websocket.id}.execute-api.ap-southeast-1.amazonaws.com/test"
    }
  }
}


resource "aws_lambda_permission" "botio_livechat_websocket_default_handler_allow_execution_form_api_gateway" {
  statement_id  = "AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.botio_livechat_websocket_default_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.botio_livechat_websocket.execution_arn}/*/$default"
}

resource "aws_lambda_permission" "botio_livechat_websocket_connect_handler_allow_execution_form_api_gateway" {
  statement_id  = "AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.botio_livechat_websocket_connect_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.botio_livechat_websocket.execution_arn}/*/$connect"
}

resource "aws_lambda_permission" "botio_livechat_websocket_disconnect_handler_allow_execution_form_api_gateway" {
  statement_id  = "AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.botio_livechat_websocket_disconnect_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.botio_livechat_websocket.execution_arn}/*/$disconnect"
}

resource "aws_lambda_permission" "botio_livechat_websocket_broadcast_handler_allow_execution_form_api_gateway" {
  statement_id  = "AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.botio_livechat_websocket_broadcast_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.botio_livechat_websocket.execution_arn}/*/broadcast"
}
resource "aws_lambda_permission" "botio_livechat_websocket_typing_broadcast_handler_allow_execution_form_api_gateway" {
  statement_id  = "AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.botio_livechat_websocket_typing_broadcast_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.botio_livechat_websocket.execution_arn}/*/typing_broadcast"
}

resource "null_resource" "build_botio_livechat_websocket_default_handler" {
  triggers = {
    source_code_hash = filebase64sha256("./botio_livechat_websocket_default_handler/src/main.go")
  }

  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./botio_livechat_websocket_default_handler/src/ -o ../bin/main ."
  }
}
resource "null_resource" "build_botio_livechat_websocket_connect_handler" {
  triggers = {
    source_code_hash = filebase64sha256("./botio_livechat_websocket_connect_handler/src/main.go")
  }

  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./botio_livechat_websocket_connect_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_botio_livechat_websocket_disconnect_handler" {
  triggers = {
    source_code_hash = filebase64sha256("./botio_livechat_websocket_disconnect_handler/src/main.go")
  }

  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./botio_livechat_websocket_disconnect_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_botio_livechat_websocket_broadcast_handler" {
  triggers = {
    source_code_hash = filebase64sha256("./botio_livechat_websocket_broadcast_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./botio_livechat_websocket_broadcast_handler/src/ -o ../bin/main ."
  }
}
resource "null_resource" "build_botio_livechat_websocket_typing_broadcast_handler" {
  triggers = {
    source_code_hash = filebase64sha256("./botio_livechat_websocket_typing_broadcast_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./botio_livechat_websocket_typing_broadcast_handler/src/ -o ../bin/main ."
  }
}

data "archive_file" "botio_livechat_websocket_default_handler" {
  source_file = "botio_livechat_websocket_default_handler/bin/main"
  type        = "zip"
  output_path = "botio_livechat_websocket_default_handler/botio_livechat_websocket_default_handler.zip"
  depends_on  = [null_resource.build_botio_livechat_websocket_default_handler]
}
data "archive_file" "botio_livechat_websocket_connect_handler" {
  source_file = "botio_livechat_websocket_connect_handler/bin/main"
  type        = "zip"
  output_path = "botio_livechat_websocket_connect_handler/botio_livechat_websocket_connect_handler.zip"
  depends_on  = [null_resource.build_botio_livechat_websocket_connect_handler]
}

data "archive_file" "botio_livechat_websocket_disconnect_handler" {
  source_file = "botio_livechat_websocket_disconnect_handler/bin/main"
  type        = "zip"
  output_path = "botio_livechat_websocket_disconnect_handler/botio_livechat_websocket_disconnect_handler.zip"
  depends_on  = [null_resource.build_botio_livechat_websocket_disconnect_handler]
}

data "archive_file" "botio_livechat_websocket_broadcast_handler" {
  source_file = "botio_livechat_websocket_broadcast_handler/bin/main"
  type        = "zip"
  output_path = "botio_livechat_websocket_broadcast_handler/botio_livechat_websocket_broadcast_handler.zip"
  depends_on  = [null_resource.build_botio_livechat_websocket_broadcast_handler]
}
data "archive_file" "botio_livechat_websocket_typing_broadcast_handler" {
  source_file = "botio_livechat_websocket_typing_broadcast_handler/bin/main"
  type        = "zip"
  output_path = "botio_livechat_websocket_typing_broadcast_handler/botio_livechat_websocket_typing_broadcast_handler.zip"
  depends_on  = [null_resource.build_botio_livechat_websocket_typing_broadcast_handler]
}

output "botio_livechat_websocket_url" {
  value = aws_apigatewayv2_stage.botio_livechat_websocket_test.invoke_url
}

