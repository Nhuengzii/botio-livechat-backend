terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

variable "platform" {
  type = string
}

variable "rest_api_id" {
  type = string
}

variable "rest_api_execution_arn" {
  type = string
}

variable "line_channel_secret" {
  type = string
}

variable "discord_webhook_url" {
  type = string
}

variable "line_channel_access_token" {
  type = string
}

variable "mongo_uri" {
  type = string
}

variable "mongo_database" {
  type = string
}

variable "mongo_collection_line_conversations" {
  type = string
}

variable "mongo_collection_line_messages" {
  type = string
}

variable "parent_id" {
  type = string
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

resource "aws_api_gateway_resource" "platform" {
  rest_api_id = var.rest_api_id
  parent_id   = var.parent_id
  path_part   = var.platform
}

resource "aws_api_gateway_resource" "page_id" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.platform.id
  path_part   = "{page_id}"
}

resource "aws_api_gateway_resource" "webhook" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.page_id.id
  path_part   = "webhook"
}

resource "aws_api_gateway_resource" "conversations" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.page_id.id
  path_part   = "conversations"
}

resource "aws_api_gateway_resource" "conversation_id" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.conversations.id
  path_part   = "{conversation_id}"
}

resource "aws_api_gateway_resource" "messages" {
  rest_api_id = var.rest_api_id
  parent_id   = aws_api_gateway_resource.conversation_id.id
  path_part   = "messages"
}

resource "aws_api_gateway_method" "get_post_webhook" {
  for_each      = toset(["GET", "POST"])
  http_method   = each.key
  rest_api_id   = var.rest_api_id
  resource_id   = aws_api_gateway_resource.webhook.id
  authorization = "NONE"
}

module "get_post_webhook_handler" {
  source                = "../lambda_handler/"
  handler_name          = format("%s_get_post_webhook_handler", var.platform)
  handler_path          = format("%s/cmd/lambda/line/validate_webhook", path.root)
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = {
    DISCORD_WEBHOOK_URL = var.discord_webhook_url
    SQS_QUEUE_URL       = aws_sqs_queue.webhook_standardizer.url
    SQS_QUEUE_ARN       = aws_sqs_queue.webhook_standardizer.arn
    LINE_CHANNEL_SECRET = var.line_channel_secret
  }
}

resource "aws_lambda_permission" "get_post_webhook" {
  statement_id  = format("AllowMethod_%s_ExecutionFromAPIGateway", each.key)
  action        = "lambda:InvokeFunction"
  function_name = module.get_post_webhook_handler.lambda.function_name
  principal     = "apigateway.amazonaws.com"
  for_each      = toset(["GET", "POST"])
  source_arn    = format("%s/*/%s%s", var.rest_api_execution_arn, each.key, aws_api_gateway_resource.webhook.path)
}

resource "aws_api_gateway_integration" "get_post_webhook" {
  for_each                = toset(["GET", "POST"])
  http_method             = aws_api_gateway_method.get_post_webhook[each.key].http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_resource.webhook.id
  rest_api_id             = var.rest_api_id
  type                    = "AWS_PROXY"
  uri                     = module.get_post_webhook_handler.lambda.invoke_arn
}

locals {
  endpoint_with_handlers = {
    get_message = {
      method                = "GET"
      resource_id           = aws_api_gateway_resource.messages.id
      resource_path         = aws_api_gateway_resource.messages.path
      handler_name          = format("%s_get_messages_handler", var.platform)
      handler_path          = format("%s/cmd/lambda/line/get_messages", path.root)
      role_arn              = aws_iam_role.assume_role_lambda.arn
      environment_variables = {
        DISCORD_WEBHOOK_URL              = var.discord_webhook_url
        MONGODB_URI                      = var.mongo_uri
        MONGODB_DATABASE                 = var.mongo_database
        MONGODB_COLLECTION_LINE_MESSAGES = var.mongo_collection_line_messages
      }
    }
    post_message = {
      method                = "POST"
      resource_id           = aws_api_gateway_resource.messages.id
      resource_path         = aws_api_gateway_resource.messages.path
      handler_name          = format("%s_post_message_handler", var.platform)
      handler_path          = format("%s/cmd/lambda/line/post_message", path.root)
      role_arn              = aws_iam_role.assume_role_lambda.arn
      environment_variables = {}
    }
    get_conversations = {
      method                = "GET"
      resource_id           = aws_api_gateway_resource.conversations.id
      resource_path         = aws_api_gateway_resource.conversations.path
      handler_name          = format("%s_get_conversations_handler", var.platform)
      handler_path          = format("%s/cmd/lambda/line/get_conversations", path.root)
      role_arn              = aws_iam_role.assume_role_lambda.arn
      environment_variables = {
        DISCORD_WEBHOOK_URL                   = var.discord_webhook_url
        MONGODB_URI                           = var.mongo_uri
        MONGODB_DATABASE                      = var.mongo_database
        MONGODB_COLLECTION_LINE_CONVERSATIONS = var.mongo_collection_line_conversations
      }
    }
  }
}

resource "aws_api_gateway_method" "endpoint_method" {
  for_each      = local.endpoint_with_handlers
  http_method   = each.value.method
  rest_api_id   = var.rest_api_id
  resource_id   = each.value.resource_id
  authorization = "NONE"
}

module "endpoint_handlers" {
  for_each              = local.endpoint_with_handlers
  source                = "../lambda_handler/"
  handler_name          = each.value.handler_name
  handler_path          = each.value.handler_path
  role_arn              = each.value.role_arn
  environment_variables = each.value.environment_variables
}

resource "aws_api_gateway_integration" "endpoint_handler_integrations" {
  for_each                = local.endpoint_with_handlers
  http_method             = aws_api_gateway_method.endpoint_method[each.key].http_method
  integration_http_method = "POST"
  resource_id             = each.value.resource_id
  rest_api_id             = var.rest_api_id
  type                    = "AWS_PROXY"
  uri                     = module.endpoint_handlers[each.key].lambda.invoke_arn
}

resource "aws_lambda_permission" "endpoint_handler_permissions" {
  for_each      = local.endpoint_with_handlers
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.endpoint_handlers[each.key].lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = format("%s/*/%s%s", var.rest_api_execution_arn, each.value.method, each.value.resource_path)
}

resource "aws_sqs_queue" "webhook_standardizer" {
  name = format("%s_webhook_standardizer", var.platform)
}

resource "aws_iam_role_policy_attachment" "sqs_full_access" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_lambda_event_source_mapping" "webhook_to_standardizer" {
  event_source_arn = aws_sqs_queue.webhook_standardizer.arn
  function_name    = module.standardizer.lambda.function_name
  batch_size       = 10
}

module "standardizer" {
  source                = "../lambda_handler/"
  handler_name          = format("%s_standardizer", var.platform)
  handler_path          = format("%s/cmd/lambda/line/standardize_webhook", path.root)
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = {
    DISCORD_WEBHOOK_URL = var.discord_webhook_url
    SNS_TOPIC_ARN       = aws_sns_topic.save_and_send_received_message.arn
    SNS_TOPIC_NAME      = aws_sns_topic.save_and_send_received_message.name
  }
}
data "aws_iam_policy_document" "sqs_allow_send_message_from_sns" {
  statement {
    sid     = "AllowSendMessageFromLineReceiveMessageTopic"
    actions = [
      "sqs:SendMessage"
    ]
    effect    = "Allow"
    resources = [
      aws_sqs_queue.save_and_send_received_message["save"].arn,
    ]
    principals {
      type        = "Service"
      identifiers = ["sns.amazonaws.com"]
    }
    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_sns_topic.save_and_send_received_message.arn]
    }
  }
}

resource "aws_iam_role_policy_attachment" "lambda_basic_sqsexecution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_iam_role_policy_attachment" "lambda_basic_snsexecution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSNSFullAccess"
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_sqs_queue_policy" "sqs_allow_send_message_from_sns" {
  for_each  = toset(["save"])
  queue_url = aws_sqs_queue.save_and_send_received_message[each.key].id
  policy    = data.aws_iam_policy_document.sqs_allow_send_message_from_sns.json
}

resource "aws_iam_role_policy_attachment" "lambda_apigateway_invoke_full_access_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess"
}

resource "aws_sns_topic" "save_and_send_received_message" {
  name = format("%s_save_and_send_receive_message", var.platform)
}

resource "aws_sqs_queue" "save_and_send_received_message" {
  for_each = toset(["save", "send"])
  name     = format("%s_%s_received_message", var.platform, each.key)
}

resource "aws_lambda_event_source_mapping" "save_received_message" {
  event_source_arn = aws_sqs_queue.save_and_send_received_message["save"].arn
  function_name    = module.save_received_message.lambda.function_name
  batch_size       = 1
}

resource "aws_sns_topic_subscription" "save_received_message" {
  topic_arn = aws_sns_topic.save_and_send_received_message.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.save_and_send_received_message["save"].arn
}

module "save_received_message" {
  source                = "../lambda_handler"
  handler_name          = format("%s_save_received_message", var.platform)
  handler_path          = format("%s/cmd/lambda/line/save_received_message", path.root)
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = {
    DISCORD_WEBHOOK_URL                   = var.discord_webhook_url
    LINE_CHANNEL_ACCESS_TOKEN             = var.line_channel_access_token
    MONGODB_URI                           = var.mongo_uri
    MONGODB_COLLECTION_LINE_MESSAGES      = var.mongo_collection_line_messages
    MONGODB_COLLECTION_LINE_CONVERSATIONS = var.mongo_collection_line_conversations
    MONGODB_DATABASE                      = var.mongo_database
  }
}