terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

provider "aws" {
  region = "ap-southeast-1"
}

resource "aws_api_gateway_rest_api" "rest_api" {
  name = "botio_rest_api"
}

resource "aws_api_gateway_resource" "shops" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_rest_api.rest_api.root_resource_id
  path_part   = "shops"
}

resource "aws_api_gateway_resource" "shop_id" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_resource.shops.id
  path_part   = "{shop_id}"
}

resource "aws_apigatewayv2_api" "botio_livechat_websocket" {
  name                       = "botio_livechat_websocket"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "botio_livechat_websocket_dev" {
  api_id      = aws_apigatewayv2_api.botio_livechat_websocket.id
  name        = "dev"
  auto_deploy = true
}

resource "aws_apigatewayv2_deployment" "botio_livechat_websocket_dev" {
  api_id      = aws_apigatewayv2_api.botio_livechat_websocket.id
  description = "dev"
  triggers = {
    always_run = timestamp()
  }
  lifecycle {
    create_before_destroy = true
  }
}

module "websocket_api" {
  source                      = "./modules/websocket_api"
  websocket_api_id            = aws_apigatewayv2_api.botio_livechat_websocket.id
  websocket_api_execution_arn = aws_apigatewayv2_api.botio_livechat_websocket.execution_arn
  redis_password              = var.redis_password
  redis_addr                  = var.redis_addr
  discord_webhook_url         = var.discord_webhook_url
}

resource "aws_api_gateway_deployment" "rest_api" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  lifecycle {
    create_before_destroy = true
  }
  triggers = {
    always_run = timestamp()
  }
  # depends_on = [module.websocket_api]
}

resource "aws_api_gateway_stage" "dev" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  deployment_id = aws_api_gateway_deployment.rest_api.id
  stage_name    = "dev"
}

output "rest_api" {
  value = {
    id = aws_api_gateway_rest_api.rest_api.id
  }
}
module "facebook" {
  source                         = "./modules/rest_api"
  platform                       = "facebook"
  rest_api_id                    = aws_api_gateway_rest_api.rest_api.id
  rest_api_execution_arn         = aws_api_gateway_rest_api.rest_api.execution_arn
  parent_id                      = aws_api_gateway_resource.shop_id.id
  relay_received_message_handler = module.websocket_api.relay_received_message_handler.function_name
  handlers = {
    validate_webhook = {
      handler_name = "facebook_validate_webhook"
      handler_path = format("%s/cmd/lambda/facebook/validate_webhook", path.root)
      environment_variables = {
        APP_SECRET                           = var.facebook_app_secret
        ACCESS_TOKEN                         = var.facebook_access_token
        FACEBOOK_WEBHOOK_VERIFICATION_STRING = var.facebook_webhook_verification_string
        DISCORD_WEBHOOK_URL                  = var.discord_webhook_url
      }
      dependencies = "{discord,db,sqswrapper}/**/*.go"
    }
    get_conversations = {
      handler_name = "facebook_get_conversations"
      handler_path = format("%s/cmd/lambda/facebook/get_conversations", path.root)
      environment_variables = {
        ACCESS_TOKEN        = var.facebook_access_token
        MONGODB_DATABASE    = var.mongo_database
        MONGODB_URI         = var.mongo_uri
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
      }
      dependencies = ""
    }
    get_conversation = {
      handler_name = "facebook_get_conversation"
      handler_path = format("%s/cmd/lambda/facebook/get_conversation", path.root)
      environment_variables = {
        ACCESS_TOKEN        = var.facebook_access_token
        MONGODB_DATABASE    = var.mongo_database
        MONGODB_URI         = var.mongo_uri
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
      }
      dependencies = ""
    }
    get_messages = {
      handler_name = "facebook_get_messages"
      handler_path = format("%s/cmd/lambda/facebook/get_messages", path.root)
      environment_variables = {
        ACCESS_TOKEN        = var.facebook_access_token
        MONGODB_DATABASE    = var.mongo_database
        MONGODB_URI         = var.mongo_uri
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
      }
      dependencies = ""
    }
    post_message = {
      handler_name = "facebook_post_message"
      handler_path = format("%s/cmd/lambda/facebook/post_message", path.root)
      environment_variables = {
        MONGODB_DATABASE    = var.mongo_database
        MONGODB_URI         = var.mongo_uri
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
      }
      dependencies = ""
    }
    standardize_webhook = {
      handler_name = "facebook_standardize_webhook"
      handler_path = format("%s/cmd/lambda/facebook/standardize_webhook", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
    save_received_message = {
      handler_name = "facebook_save_received_message"
      handler_path = format("%s/cmd/lambda/facebook/save_received_message", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        ACCESS_TOKEN        = var.facebook_access_token
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
  }
  method_integrations = {
    get_validate_webhook = {
      method  = "GET"
      handler = "validate_webhook"
    }
    post_validate_webhook = {
      method  = "POST"
      handler = "validate_webhook"
    }
    get_conversations = {
      method  = "GET"
      handler = "get_conversations"
    }
    get_conversation = {
      method  = "GET"
      handler = "get_conversation"
    }
    get_messages = {
      method  = "GET"
      handler = "get_messages"
    }
    post_message = {
      method  = "POST"
      handler = "post_message"
    }
  }
}

module "line" {
  source                         = "./modules/rest_api"
  platform                       = "line"
  rest_api_id                    = aws_api_gateway_rest_api.rest_api.id
  rest_api_execution_arn         = aws_api_gateway_rest_api.rest_api.execution_arn
  parent_id                      = aws_api_gateway_resource.shop_id.id
  relay_received_message_handler = module.websocket_api.relay_received_message_handler.function_name
  handlers = {
    validate_webhook = {
      handler_name = "line_validate_webhook"
      handler_path = format("%s/cmd/lambda/line/validate_webhook", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
    get_conversations = {
      handler_name = "line_get_conversations"
      handler_path = format("%s/cmd/lambda/line/get_conversations", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
    get_conversation = {
      handler_name = "line_get_conversation"
      handler_path = format("%s/cmd/lambda/line/get_conversation", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
    get_messages = {
      handler_name = "line_get_messages"
      handler_path = format("%s/cmd/lambda/line/get_messages", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
    post_message = {
      handler_name = "line_post_message"
      handler_path = format("%s/cmd/lambda/line/post_message", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
    standardize_webhook = {
      handler_name = "line_standardize_webhook"
      handler_path = format("%s/cmd/lambda/line/standardize_webhook", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
    save_received_message = {
      handler_name = "line_save_received_message"
      handler_path = format("%s/cmd/lambda/line/save_received_message", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = ""
    }
  }
  method_integrations = {
    post_validate_webhook = {
      method  = "POST"
      handler = "validate_webhook"
    }
    get_conversations = {
      method  = "GET"
      handler = "get_conversations"
    }
    get_conversation = {
      method  = "GET"
      handler = "get_conversation"
    }
    get_messages = {
      method  = "GET"
      handler = "get_messages"
    }
    post_message = {
      method  = "POST"
      handler = "post_message"
    }
  }
}

