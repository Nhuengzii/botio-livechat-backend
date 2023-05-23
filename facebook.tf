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

resource "aws_api_gateway_resource" "facebook_conversation" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.facebook_page_id.id
  path_part   = "conversations"
}

resource "aws_api_gateway_resource" "facebook_conversation_id" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.facebook_conversation.id
  path_part   = "{conversation_id}"
}

resource "aws_api_gateway_resource" "facebook_message" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.facebook_conversation_id.id
  path_part   = "messages"
}

module "aws_api_gateway_enable_cors" {
  source          = "squidfunk/api-gateway-enable-cors/aws"
  version         = "0.3.3"
  api_id          = aws_api_gateway_rest_api.botio_rest_api.id
  api_resource_id = aws_api_gateway_resource.facebook_message.id
}

resource "aws_api_gateway_method" "post_facebook_message" {
  http_method   = "POST"
  resource_id   = aws_api_gateway_resource.facebook_message.id
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  authorization = "NONE"
}

resource "aws_api_gateway_method" "get_facebook_conversation" {
  http_method   = "GET"
  resource_id   = aws_api_gateway_resource.facebook_conversation.id
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  authorization = "NONE"
}

resource "aws_api_gateway_method" "get_facebook_messages" {
  http_method   = "GET"
  resource_id   = aws_api_gateway_resource.facebook_message.id
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  authorization = "NONE"
}


resource "aws_api_gateway_method" "get_validate_facebook_webhook" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.facebook_webhook.id
  authorization = "NONE"
  http_method   = "GET"
}
resource "aws_api_gateway_method" "post_validate_facebook_webhook" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.facebook_webhook.id
  authorization = "NONE"
  http_method   = "POST"
}

resource "aws_sqs_queue" "facebook_webhook_to_standardize_facebook_webhook_handler" {
  name = "facebook_webhook_to_standardize_facebook_webhook_handler"
}

resource "aws_sns_topic" "facebook_receive_message" {
  name = "facebook_receive_message"
}

resource "aws_sqs_queue" "facebook_receive_message_to_database" {
  name = "facebook_receive_message_to_database"
}

resource "aws_sqs_queue" "facebook_receive_message_to_frontend" {
  name = "facebook_receive_message_to_frontend"
}

data "aws_iam_policy_document" "sqs_allow_send_message_from_facebook_receive_message_topic" {
  statement {
    sid = "AllowSendMessageFromFacebookReceiveMessageTopic"
    actions = [
      "sqs:SendMessage"
    ]
    effect = "Allow"
    resources = [
      aws_sqs_queue.facebook_receive_message_to_database.arn,
      aws_sqs_queue.facebook_receive_message_to_frontend.arn
    ]
    principals {
      type        = "Service"
      identifiers = ["sns.amazonaws.com"]
    }
    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_sns_topic.facebook_receive_message.arn]
    }
  }
}

resource "aws_sqs_queue_policy" "facebook_receive_message_to_database_allow_send_message_from_facebook_receive_message_topic" {
  queue_url = aws_sqs_queue.facebook_receive_message_to_database.id
  policy    = data.aws_iam_policy_document.sqs_allow_send_message_from_facebook_receive_message_topic.json
}

resource "aws_sqs_queue_policy" "facebook_recieve_message_to_frontend_allow_send_message_from_facebook_receive_message_topic" {
  queue_url = aws_sqs_queue.facebook_receive_message_to_frontend.id
  policy    = data.aws_iam_policy_document.sqs_allow_send_message_from_facebook_receive_message_topic.json
}


resource "aws_lambda_event_source_mapping" "event_source_mapping_facebook_webhook_to_standardize_facebook_webhook_handler" {
  event_source_arn = aws_sqs_queue.facebook_webhook_to_standardize_facebook_webhook_handler.arn
  function_name    = aws_lambda_function.standardize_facebook_webhook_handler.arn
  batch_size       = 10
}


resource "aws_api_gateway_integration" "get_facebook_messages" {
  http_method             = aws_api_gateway_method.get_facebook_messages.http_method
  resource_id             = aws_api_gateway_resource.facebook_message.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.get_facebook_messages_handler.invoke_arn
}

resource "aws_lambda_event_source_mapping" "event_source_mapping_facebook_recieve_message_to_save_facebook_received_message_handler" {
  event_source_arn = aws_sqs_queue.facebook_receive_message_to_database.arn
  function_name    = aws_lambda_function.save_facebook_received_message_handler.arn
  batch_size       = 1
}

resource "aws_lambda_event_source_mapping" "event_source_mapping_facebook_recieve_message_to_send_facebook_received_message_handler" {
  event_source_arn = aws_sqs_queue.facebook_receive_message_to_frontend.arn
  function_name    = aws_lambda_function.send_facebook_received_message_handler.arn
  batch_size       = 1
}


resource "aws_sns_topic_subscription" "facebook_recieve_message_to_database" {
  topic_arn = aws_sns_topic.facebook_receive_message.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.facebook_receive_message_to_database.arn
}

resource "aws_sns_topic_subscription" "facebook_recieve_message_to_frontend" {
  topic_arn = aws_sns_topic.facebook_receive_message.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.facebook_receive_message_to_frontend.arn
}


resource "aws_api_gateway_integration" "get_validate_facebook_webhook" {
  http_method             = aws_api_gateway_method.get_validate_facebook_webhook.http_method
  resource_id             = aws_api_gateway_resource.facebook_webhook.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.validate_facebook_webhook_handler.invoke_arn
}
resource "aws_api_gateway_integration" "post_validate_facebook_webhook" {
  http_method             = aws_api_gateway_method.post_validate_facebook_webhook.http_method
  resource_id             = aws_api_gateway_resource.facebook_webhook.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.validate_facebook_webhook_handler.invoke_arn
}

resource "aws_api_gateway_integration" "post_facebook_message_handler" {
  http_method             = aws_api_gateway_method.post_facebook_message.http_method
  resource_id             = aws_api_gateway_resource.facebook_message.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.post_facebook_message_handler.invoke_arn
}

resource "aws_api_gateway_integration" "get_facebook_conversation" {
  http_method             = aws_api_gateway_method.get_facebook_conversation.http_method
  resource_id             = aws_api_gateway_resource.facebook_conversation.id
  rest_api_id             = aws_api_gateway_rest_api.botio_rest_api.id
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.get_facebook_conversation_handler.invoke_arn
}

variable "facebook_access_token" {
  type = string
}

variable "facebook_app_secret" {
  type = string
}

resource "aws_lambda_function" "get_facebook_conversation_handler" {
  filename         = "get_facebook_conversation_handler/get_facebook_conversation_handler.zip"
  function_name    = "get_facebook_conversation_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.get_facebook_conversation_handler.output_base64sha256
  depends_on       = [data.archive_file.get_facebook_conversation_handler]
  environment {
    variables = {
      ACCESS_TOKEN = var.facebook_access_token
    }
  }
}

resource "aws_lambda_function" "validate_facebook_webhook_handler" {
  filename         = "validate_facebook_webhook_handler/validate_facebook_webhook_handler.zip"
  function_name    = "validate_facebook_webhook_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.validate_facebook_webhook_handler.output_base64sha256
  depends_on       = [data.archive_file.validate_facebook_webhook_handler]
  environment {
    variables = {
      SQS_QUEUE_URL = aws_sqs_queue.facebook_webhook_to_standardize_facebook_webhook_handler.url
      SQS_QUEUE_ARN = aws_sqs_queue.facebook_webhook_to_standardize_facebook_webhook_handler.arn
      foo           = "bar"
      ACCESS_TOKEN  = var.facebook_access_token
      APP_SECRET    = var.facebook_app_secret
    }
  }
}

resource "aws_lambda_function" "standardize_facebook_webhook_handler" {
  filename         = "standardize_facebook_webhook_handler/standardize_facebook_webhook_handler.zip"
  function_name    = "standardize_facebook_webhook_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.standardize_facebook_webhook_handler.output_base64sha256
  depends_on       = [data.archive_file.standardize_facebook_webhook_handler]
  environment {
    variables = {
      SNS_TOPIC_ARN = aws_sns_topic.facebook_receive_message.arn
      SNS_TOPIC_URL = aws_sns_topic.facebook_receive_message.arn
      foo           = "bar"
      ACCESS_TOKEN  = var.facebook_access_token
    }
  }
}

resource "aws_lambda_function" "post_facebook_message_handler" {
  filename         = "post_facebook_message_handler/post_facebook_message_handler.zip"
  function_name    = "post_facebook_message_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.post_facebook_message_handler.output_base64sha256
  depends_on       = [data.archive_file.post_facebook_message_handler]
  timeout          = 10
  environment {
    variables = {
      ACCESS_TOKEN = var.facebook_access_token
      APP_SECRET   = var.facebook_app_secret
    }
  }
}

resource "aws_lambda_function" "get_facebook_messages_handler" {
  filename         = "get_facebook_messages_handler/get_facebook_messages_handler.zip"
  function_name    = "get_facebook_messages_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.get_facebook_messages_handler.output_base64sha256
  depends_on       = [data.archive_file.get_facebook_messages_handler]
  environment {
    variables = {
      ACCESS_TOKEN = var.facebook_access_token
    }
  }
}

resource "aws_lambda_function" "save_facebook_received_message_handler" {
  filename         = "save_facebook_received_message_handler/save_facebook_recieved_message_handler.zip"
  function_name    = "save_facebook_received_message_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.save_facebook_received_message_handler.output_base64sha256
  depends_on       = [data.archive_file.save_facebook_received_message_handler]
  environment {
    variables = {
      ACCESS_TOKEN = var.facebook_access_token
    }
  }
}

resource "aws_lambda_function" "send_facebook_received_message_handler" {
  filename         = "send_facebook_received_message_handler/send_facebook_recieved_message_handler.zip"
  function_name    = "send_facebook_received_message_handler"
  role             = aws_iam_role.assume_role_lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.send_facebook_received_message_handler.output_base64sha256
  depends_on       = [data.archive_file.send_facebook_received_message_handler]
  environment {
    variables = {
      WEBSOCKET_API_ENDPOINT = "https://${aws_apigatewayv2_api.botio_livechat_websocket.id}.execute-api.ap-southeast-1.amazonaws.com/test"
      ACCESS_TOKEN           = var.facebook_access_token
      REDIS_ACCESS_ADDR      = var.redis_access.addr
      REDIS_ACCESS_PASSWORD  = var.redis_access.password
    }
  }
}


resource "aws_lambda_permission" "validate_facebook_webhook_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.validate_facebook_webhook_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "post_facebook_message_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.post_facebook_message_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "get_facebook_messages_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_facebook_messages_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "get_facebook_conversation_handler_allow_execution_from_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_facebook_conversation_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.botio_rest_api.execution_arn}/*/*/*"
}


resource "null_resource" "build_validate_facebook_webhook_handler" {
  triggers = {
    source_code_hash  = "${filebase64sha256("validate_facebook_webhook_handler/src/main.go")}"
    source_code_hash1 = "${filebase64sha256("validate_facebook_webhook_handler/src/sendQueueMessage.go")}"
    source_code_hash2 = "${filebase64sha256("validate_facebook_webhook_handler/src/verificationCheck.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./validate_facebook_webhook_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_standardize_facebook_webhook_handler" {
  triggers = {
    source_code_hash  = "${filebase64sha256("standardize_facebook_webhook_handler/src/main.go")}"
    source_code_hash1 = "${filebase64sha256("standardize_facebook_webhook_handler/src/recieveMessage.go")}"
    source_code_hash2 = "${filebase64sha256("standardize_facebook_webhook_handler/src/standardMessageStruct.go")}"
    source_code_hash3 = "${filebase64sha256("standardize_facebook_webhook_handler/src/standardizeMessage.go")}"
    source_code_hash4 = "${filebase64sha256("standardize_facebook_webhook_handler/src/sendSnsMessage.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./standardize_facebook_webhook_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_post_facebook_message_handler" {
  triggers = {
    source_code_hash = "${filebase64sha256("post_facebook_message_handler/src/main.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./post_facebook_message_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_get_facebook_conversation_handler" {
  triggers = {
    source_code_hash = "${filebase64sha256("get_facebook_conversation_handler/src/main.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./get_facebook_conversation_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_save_facebook_received_message_handler" {
  triggers = {
    source_code_hash = "${filebase64sha256("save_facebook_received_message_handler/src/main.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./save_facebook_received_message_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_send_facebook_received_message_handler" {
  triggers = {
    source_code_hash = "${filebase64sha256("send_facebook_received_message_handler/src/main.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./send_facebook_received_message_handler/src/ -o ../bin/main ."
  }
}

resource "null_resource" "build_get_facebook_messages_handler" {
  triggers = {
    source_code_hash = "${filebase64sha256("get_facebook_messages_handler/src/main.go")}"
  }
  provisioner "local-exec" {
    command = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./get_facebook_messages_handler/src/ -o ../bin/main ."
  }
}

data "archive_file" "get_facebook_conversation_handler" {
  type        = "zip"
  source_file = "get_facebook_conversation_handler/bin/main"
  output_path = "get_facebook_conversation_handler/get_facebook_conversation_handler.zip"
  depends_on  = [null_resource.build_get_facebook_conversation_handler]
}

data "archive_file" "validate_facebook_webhook_handler" {
  type        = "zip"
  source_file = "./validate_facebook_webhook_handler/bin/main"
  output_path = "./validate_facebook_webhook_handler/validate_facebook_webhook_handler.zip"
  depends_on  = [null_resource.build_validate_facebook_webhook_handler]
}

data "archive_file" "standardize_facebook_webhook_handler" {
  type        = "zip"
  source_file = "./standardize_facebook_webhook_handler/bin/main"
  output_path = "./standardize_facebook_webhook_handler/standardize_facebook_webhook_handler.zip"
  depends_on  = [null_resource.build_standardize_facebook_webhook_handler]
}

data "archive_file" "save_facebook_received_message_handler" {
  type        = "zip"
  source_file = "./save_facebook_received_message_handler/bin/main"
  output_path = "./save_facebook_received_message_handler/save_facebook_recieved_message_handler.zip"
  depends_on  = [null_resource.build_save_facebook_received_message_handler]
}


data "archive_file" "post_facebook_message_handler" {
  type        = "zip"
  source_file = "./post_facebook_message_handler/bin/main"
  output_path = "./post_facebook_message_handler/post_facebook_message_handler.zip"
  depends_on  = [null_resource.build_post_facebook_message_handler]
}

data "archive_file" "send_facebook_received_message_handler" {
  type        = "zip"
  source_file = "./send_facebook_received_message_handler/bin/main"
  output_path = "./send_facebook_received_message_handler/send_facebook_recieved_message_handler.zip"
  depends_on  = [null_resource.build_send_facebook_received_message_handler]
}



data "archive_file" "get_facebook_messages_handler" {
  type        = "zip"
  source_file = "./get_facebook_messages_handler/bin/main"
  output_path = "./get_facebook_messages_handler/get_facebook_messages_handler.zip"
  depends_on  = [null_resource.build_get_facebook_messages_handler]
}

