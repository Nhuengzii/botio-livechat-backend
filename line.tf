module "line" {
  source                         = "./modules/platform_rest_api"
  platform                       = "line"
  rest_api_id                    = module.rest_api.id
  rest_api_execution_arn         = module.rest_api.execution_arn
  parent_id                      = aws_api_gateway_resource.shop_id.id
  relay_received_message_handler = module.websocket_api.relay_received_message_handler.function_name
  bucket_name                    = module.bucket.bucket_name
  bucket_arn                     = module.bucket.bucket_arn
  handlers = {
    get_page_id = {
      handler_name = "line_get_page_id"
      handler_path = format("%s/cmd/lambda/line/get_page_id", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,apigateway}/**/*.go"
    }
    validate_webhook = {
      handler_name = "line_validate_webhook"
      handler_path = format("%s/cmd/lambda/line/validate_webhook", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,sqswrapper,apigateway}/**/*.go"
    }
    get_conversations = {
      handler_name = "line_get_conversations"
      handler_path = format("%s/cmd/lambda/line/get_conversations", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,api,stdconversation,apigateway}/**/*.go"
    }
    get_conversation = {
      handler_name = "line_get_conversation"
      handler_path = format("%s/cmd/lambda/line/get_conversation", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,api,apigateway}/**/*.go"
    }
    patch_conversation = {
      handler_name = "line_patch_conversation"
      handler_path = format("%s/cmd/lambda/line/patch_conversation", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,api}/**/*.go"
    }
    get_messages = {
      handler_name = "line_get_messages"
      handler_path = format("%s/cmd/lambda/line/get_messages", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,api,stdmessage,apigateway}/**/*.go"
    }
    post_message = {
      handler_name = "line_post_message"
      handler_path = format("%s/cmd/lambda/line/post_message", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,api,stdmessage,apigateway}/**/*.go"
    }
    standardize_webhook = {
      handler_name = "line_standardize_webhook"
      handler_path = format("%s/cmd/lambda/line/standardize_webhook", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,snswrapper,storage,stdmessage}/**/*.go"
    }
    save_received_message = {
      handler_name = "line_save_received_message"
      handler_path = format("%s/cmd/lambda/line/save_received_message", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,snswrapper,stdmessage,stdconversation,external_api}/**/*.go"
    }
  }
  method_integrations = {
    get_page_id = {
      method  = "GET"
      handler = "get_page_id"
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
    patch_conversation = {
      method  = "PATCH"
      handler = "patch_conversation"
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
