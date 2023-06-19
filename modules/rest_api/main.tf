terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

# Define role
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

# Define Queue
resource "aws_sqs_queue" "webhook_standardizer" {
  name = format("%s_webhook_standardizer", var.platform)
}

resource "aws_sqs_queue" "save_and_relay_received_message" {
  for_each = toset(["save", "relay"])
  name     = format("%s_%s_received_message", var.platform, each.key)
}

data "aws_iam_policy_document" "sqs_allow_send_message_from_sns" {
  statement {
    sid = "AllowSendMessageFromFacebookReceiveMessageTopic"
    actions = [
      "sqs:SendMessage"
    ]
    effect = "Allow"
    resources = [
      aws_sqs_queue.save_and_relay_received_message["save"].arn,
      # aws_sqs_queue.relay_received_message.arn
    ]
    principals {
      type        = "Service"
      identifiers = ["sns.amazonaws.com"]
    }
    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_sns_topic.save_and_relay_received_message.arn]
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

resource "aws_lambda_event_source_mapping" "save_received_message" {
  event_source_arn = aws_sqs_queue.save_and_relay_received_message["save"].arn
  function_name    = module.handlers["save_received_message"].lambda.function_name
  batch_size       = 1
}

resource "aws_sns_topic_subscription" "save_received_message" {
  topic_arn = aws_sns_topic.save_and_relay_received_message.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.save_and_relay_received_message["save"].arn
}

# resource "aws_lambda_event_source_mapping" "relay_received_message" {
#   event_source_arn = aws_sqs_queue.relay_received_message.arn
#   function_name    = var.relay_received_message_handler.function_name
#   batch_size       = 1
# }

# Define Topic
resource "aws_sns_topic" "save_and_relay_received_message" {
  name = format("%s_save_and_relay_received_message", var.platform)
}

resource "aws_iam_role_policy_attachment" "sqs_full_access" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_lambda_event_source_mapping" "webhook_to_standardizer" {
  event_source_arn = aws_sqs_queue.webhook_standardizer.arn
  function_name    = module.handlers["standardize_webhook"].lambda.function_name
  batch_size       = 10
}



# Define Resource

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

locals {
  resource_id_mapping = {
    validate_webhook  = aws_api_gateway_resource.webhook.id
    get_conversations = aws_api_gateway_resource.conversations.id
    post_conversation = aws_api_gateway_resource.conversations.id
    get_conversation  = aws_api_gateway_resource.conversation_id.id
    get_messages      = aws_api_gateway_resource.messages.id
    post_message      = aws_api_gateway_resource.messages.id
  }
  resource_path_mapping = {
    validate_webhook  = aws_api_gateway_resource.webhook.path
    get_conversations = aws_api_gateway_resource.conversations.path
    post_conversation = aws_api_gateway_resource.conversations.path
    get_conversation  = aws_api_gateway_resource.conversation_id.path
    get_messages      = aws_api_gateway_resource.messages.path
    post_message      = aws_api_gateway_resource.messages.path
  }
  environment_variables_mapping = {
    validate_webhook = {
      SQS_QUEUE_URL = aws_sqs_queue.webhook_standardizer.id
      SQS_QUEUE_ARN = aws_sqs_queue.webhook_standardizer.arn
    }
    standardize_webhook = {
      SNS_TOPIC_ARN = aws_sns_topic.save_and_relay_received_message.arn
    }
    save_received_message = {}
    get_conversations     = {}
    post_conversation     = {}
    get_conversation      = {}
    get_messages          = {}
    post_message          = {}
  }
}

# Define Handler

module "handlers" {
  source   = "../lambda_handler"
  for_each = var.handlers

  handler_name          = each.value.handler_name
  handler_path          = each.value.handler_path
  role_arn              = aws_iam_role.assume_role_lambda.arn
  environment_variables = merge(each.value.environment_variables, local.environment_variables_mapping[each.key])
}

# Define Method
resource "aws_api_gateway_method" "methods" {
  for_each      = var.method_integrations
  http_method   = each.value.method
  rest_api_id   = var.rest_api_id
  resource_id   = local.resource_id_mapping[each.value.handler]
  authorization = "NONE"
}


# Define Integration
resource "aws_api_gateway_integration" "integrations" {
  for_each = var.method_integrations

  http_method             = aws_api_gateway_method.methods[each.key].http_method
  integration_http_method = "POST"
  resource_id             = local.resource_id_mapping[each.value.handler]
  rest_api_id             = var.rest_api_id
  type                    = "AWS_PROXY"
  uri                     = module.handlers[each.value.handler].lambda.invoke_arn
}

# Define Permission
resource "aws_lambda_permission" "endpoint_handler_permissions" {
  for_each      = var.method_integrations
  function_name = module.handlers[each.value.handler].lambda.function_name
  statement_id  = format("AllowMethod_%s_ExecutionFromAPIGateway", each.key)
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = format("%s/*/%s%s", var.rest_api_execution_arn, each.value.method, local.resource_path_mapping[each.value.handler])
}
