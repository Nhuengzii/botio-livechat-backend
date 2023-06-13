terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

locals {
  redis_addr     = var.redis_addr
  redis_password = var.redis_password
  routes_with_handler = {
    connect = {
      route_key    = "$connect"
      handler_name = "connect"
      handler_path = format("%s/cmd/lambda/websocket/connect", path.root)
      environment_variables = {
      }
    }
    disconnect = {
      route_key    = "$disconnect"
      handler_name = "disconnect"
      handler_path = format("%s/cmd/lambda/websocket/disconnect", path.root)
      environment_variables = {
      }
    }
    default = {
      route_key    = "$default"
      handler_name = "default"
      handler_path = format("%s/cmd/lambda/websocket/default", path.root)
      environment_variables = {
      }
    }
    broadcast = {
      route_key    = "broadcast"
      handler_name = "broadcast"
      handler_path = format("%s/cmd/lambda/websocket/broadcast", path.root)
      environment_variables = {
      }
    }
  }
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
  name               = "assume_role_for_route_handlers"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "basic_execution" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_permission" "allow_execution_from_api_websocket" {
  for_each      = local.routes_with_handler
  statement_id  = "AllowExecutionFrom_${each.key}_Route"
  action        = "lambda:InvokeFunction"
  function_name = module.routes_handler[each.key].lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = format("%s/*/%s", var.websocket_api_execution_arn, each.value.route_key)
}


module "routes_handler" {
  source = "../lambda_handler"

  for_each     = local.routes_with_handler
  handler_name = each.value.handler_name
  handler_path = each.value.handler_path
  environment_variables = {
    REDIS_ADDR     = local.redis_addr
    REDIS_PASSWORD = local.redis_password
  }
  role_arn = aws_iam_role.assume_role_lambda.arn

}

resource "aws_apigatewayv2_route" "routes_with_handler" {
  for_each  = local.routes_with_handler
  route_key = each.value.route_key
  target    = format("integrations/%s", aws_apigatewayv2_integration.route_handlers[each.key].id)
  api_id    = var.websocket_api_id
}

resource "aws_apigatewayv2_integration" "route_handlers" {
  for_each                  = local.routes_with_handler
  api_id                    = var.websocket_api_id
  integration_type          = "AWS_PROXY"
  integration_uri           = module.routes_handler[each.key].lambda.invoke_arn
  integration_method        = "POST"
  content_handling_strategy = "CONVERT_TO_TEXT"
  passthrough_behavior      = "WHEN_NO_MATCH"
}

resource "aws_sqs_queue" "relay_received_message" {
  name = "relay_received_message"
}

resource "aws_iam_role_policy_attachment" "lambda_basic_sqsexecution_to_assume_role_lambda" {
  role       = aws_iam_role.assume_role_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_lambda_event_source_mapping" "relay_received_message" {
  event_source_arn = aws_sqs_queue.relay_received_message.arn
  function_name    = module.relay_received_message.lambda.function_name
  batch_size       = 1
}

module "relay_received_message" {
  source = "../lambda_handler"

  handler_name = "relay_received_message"
  handler_path = format("%s/cmd/lambda/websocket/relay", path.root)
  role_arn     = aws_iam_role.assume_role_lambda.arn
  environment_variables = {
    REDIS_ADDR     = local.redis_addr
    REDIS_PASSWORD = local.redis_password
  }
}

output "relay_received_message_queue" {
  value = aws_sqs_queue.relay_received_message
}
