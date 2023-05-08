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

resource "aws_api_gateway_resource" "line_conversation" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.line_page_id.id
  path_part   = "conversations"
}

resource "aws_api_gateway_resource" "line_conversation_id" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.line_conversation.id
  path_part   = "{conversation_id}"
}

resource "aws_api_gateway_resource" "line_message" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.line_conversation_id.id
  path_part   = "messages"
}

resource "aws_api_gateway_method" "post_validate_line_webhook" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.line_webhook.id
  authorization = "NONE"
  http_method   = "POST"
}

resource "aws_api_gateway_method" "post_line_message" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.line_message.id
  authorization = "NONE"
  http_method   = "POST"
}

resource "aws_api_gateway_method" "get_line_conversations" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.line_conversation.id
  authorization = "NONE"
  http_method   = "GET"
}

resource "aws_api_gateway_method" "get_line_messages" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.line_message.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "get_line_messages" {
  http_method             = aws_api_gateway_method.get_line_messages.http_method
  resource_id             = aws_api_gateway_resource.line_message.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.get_line_messages_handler.invoke_arn
}

resource "aws_api_gateway_integration" "get_line_conversations" {
  http_method             = aws_api_gateway_method.get_line_conversations.http_method
  resource_id             = aws_api_gateway_resource.line_conversation.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.get_line_conversations_handler.invoke_arn
}

resource "aws_api_gateway_integration" "post_validate_line_webhook" {
  http_method             = aws_api_gateway_method.post_validate_line_webhook.http_method
  resource_id             = aws_api_gateway_resource.line_webhook.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.validate_line_webhook_handler.invoke_arn
}

resource "aws_api_gateway_integration" "post_line_message" {
  http_method             = aws_api_gateway_method.post_line_message.http_method
  resource_id             = aws_api_gateway_resource.line_message.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.post_line_message_handler.invoke_arn
}

resource "aws_lambda_permission" "get_line_messages_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_line_messages_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "validate_line_webhook_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.validate_line_webhook_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "post_line_message_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.post_line_message_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "get_line_conversations_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_line_conversations_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "aws_sns_topic" "line_receive_message" {
  name = "line_receive_message"
}

resource "aws_sqs_queue" "line_receive_message_to_frontend" {
  name = "line_receive_message_to_frontend"
}

resource "aws_sqs_queue" "line_receive_message_to_database" {
  name = "line_receive_message_to_database"
}

resource "aws_sqs_queue" "line_webhook_to_standardize_line_webhook_handler" {
  name = "line_webhook_to_standardize_facebook_webhook_handler"
}

resource "aws_lambda_event_source_mapping" "event_source_mapping_line_webhook_to_standardize_line_webhook_handler" {
  event_source_arn = aws_sqs_queue.line_webhook_to_standardize_line_webhook_handler.arn
  function_name    = aws_lambda_function.standardize_line_webhook_handler.function_name
  batch_size       = 1
}


resource "aws_lambda_event_source_mapping" "event_source_mapping_line_receive_message_to_frontend" {
  event_source_arn = aws_sqs_queue.line_receive_message_to_frontend.arn
  function_name    = aws_lambda_function.send_line_received_message_handler.function_name
  batch_size       = 1
}

resource "aws_lambda_event_source_mapping" "event_source_mapping_line_receive_message_to_database" {
  event_source_arn = aws_sqs_queue.line_receive_message_to_database.arn
  function_name    = aws_lambda_function.save_line_received_message_handler.function_name
  batch_size       = 1
}

resource "aws_sns_topic_subscription" "line_receive_message_to_frontend" {
  topic_arn = aws_sns_topic.line_receive_message.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.line_receive_message_to_frontend.arn
}

resource "aws_sns_topic_subscription" "line_receive_message_to_database" {
  topic_arn = aws_sns_topic.line_receive_message.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.line_receive_message_to_database.arn
}

resource "null_resource" "build_get_line_messages_handler" {
  triggers = {
    source_code_hash = filebase64sha256("get_line_messages_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./get_line_messages_handler/src/ -o ../bin/main ."
  }
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

resource "null_resource" "build_post_line_message_handler" {
  triggers = {
    source_code_hash = filebase64sha256("post_line_message_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./post_line_message_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_get_line_conversations_handler" {
  triggers = {
    source_code_hash = filebase64sha256("get_line_conversations_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./get_line_conversations_handler/src/ -o ../bin/main ."
  }
}



resource "null_resource" "build_send_line_received_message_handler" {
  triggers = {
    source_code_hash = filebase64sha256("send_line_received_message_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./send_line_received_message_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_save_line_received_message_handler" {
  triggers = {
    source_code_hash = filebase64sha256("save_line_received_message_handler/src/main.go")
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./save_line_received_message_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "get_line_conversations_handler" {
  triggers = {
    source_code_hash = filebase64sha256("get_line_conversations_handler/src/main.go")
  }
  depends_on = [data.archive_file.get_line_conversations_handler]
}

data "archive_file" "get_line_messages_handler" {
  type        = "zip"
  source_file = "get_line_messages_handler/bin/main"
  output_path = "get_line_messages_handler/get_line_messages_handler.zip"
  depends_on  = [null_resource.build_get_line_messages_handler]
}

data "archive_file" "get_line_conversations_handler" {
  type        = "zip"
  source_file = "get_line_conversations_handler/bin/main"
  output_path = "get_line_conversations_handler/get_line_conversations_handler.zip"
  depends_on  = [null_resource.build_get_line_conversations_handler]
}

data "archive_file" "validate_line_webhook_handler" {
  type        = "zip"
  source_file = "validate_line_webhook_handler/bin/main"
  output_path = "validate_line_webhook_handler/validate_line_webhook_handler.zip"
  depends_on  = [null_resource.build_validate_line_webhook_handler]
}

data "archive_file" "standardize_line_webhook_handler" {
  type        = "zip"
  source_file = "standardize_line_webhook_handler/bin/main"
  output_path = "standardize_line_webhook_handler/standardize_line_webhook_handler.zip"
  depends_on  = [null_resource.build_standardize_line_webhook_handler]
}

data "archive_file" "post_line_message_handler" {
  type        = "zip"
  source_file = "post_line_message_handler/bin/main"
  output_path = "post_line_message_handler/post_line_message_handler.zip"
  depends_on  = [null_resource.build_post_line_message_handler]
}

data "archive_file" "save_line_received_message_handler" {
  type        = "zip"
  source_file = "save_line_received_message_handler/bin/main"
  output_path = "save_line_received_message_handler/save_line_received_message_handler.zip"
  depends_on  = [null_resource.build_save_line_received_message_handler]
}

data "archive_file" "send_line_received_message_handler" {
  type        = "zip"
  source_file = "send_line_received_message_handler/bin/main"
  output_path = "send_line_received_message_handler/send_line_received_message_handler.zip"
  depends_on  = [null_resource.build_send_line_received_message_handler]
}

resource "aws_lambda_function" "validate_line_webhook_handler" {
  filename         = data.archive_file.validate_line_webhook_handler.output_path
  function_name    = "validate_line_webhook_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.validate_line_webhook_handler.output_base64sha256
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
  filename         = data.archive_file.standardize_line_webhook_handler.output_path
  function_name    = "standardize_line_webhook_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.standardize_line_webhook_handler.output_base64sha256
  depends_on       = [data.archive_file.standardize_line_webhook_handler]
}

resource "aws_lambda_function" "post_line_message_handler" {
  filename         = data.archive_file.post_line_message_handler.output_path
  function_name    = "post_line_message_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.post_line_message_handler.output_base64sha256
  depends_on       = [data.archive_file.post_line_message_handler]
}

resource "aws_lambda_function" "get_line_conversations_handler" {
  filename         = data.archive_file.get_line_conversations_handler.output_path
  function_name    = "get_line_conversations_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.get_line_conversations_handler.output_base64sha256
  depends_on       = [data.archive_file.get_line_conversations_handler]
}

resource "aws_lambda_function" "get_line_messages_handler" {
  filename         = data.archive_file.get_line_messages_handler.output_path
  function_name    = "get_line_messages_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.get_line_messages_handler.output_base64sha256
  depends_on       = [data.archive_file.get_line_messages_handler]
}

resource "aws_lambda_function" "save_line_received_message_handler" {
  filename         = "save_line_received_message_handler/save_line_received_message_handler.zip"
  function_name    = "save_line_received_message_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  source_code_hash = data.archive_file.save_line_received_message_handler.output_base64sha256
  handler          = "main"
  runtime          = "go1.x"
  depends_on       = [data.archive_file.save_line_received_message_handler]
}

resource "aws_lambda_function" "send_line_received_message_handler" {
  filename         = "send_line_received_message_handler/send_line_received_message_handler.zip"
  function_name    = "send_line_received_message_handler"
  source_code_hash = data.archive_file.send_line_received_message_handler.output_base64sha256
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  depends_on       = [data.archive_file.send_line_received_message_handler]
}
