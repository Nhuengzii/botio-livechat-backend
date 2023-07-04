module "facebook" {
  source                         = "./modules/platform_rest_api"
  platform                       = "facebook"
  rest_api_id                    = module.rest_api.id
  rest_api_execution_arn         = module.rest_api.execution_arn
  parent_id                      = module.shops.shop_id_resource_id
  relay_received_message_handler = module.websocket_api.relay_received_message_handler.function_name
  bucket_name                    = module.bucket.bucket_name
  bucket_arn                     = module.bucket.bucket_arn
  handlers = {
    get_page_id = {
      handler_name = "facebook_get_page_id"
      handler_path = format("%s/cmd/lambda/facebook/get_page_id", path.root)
      environment_variables = {
        MONGODB_DATABASE    = var.mongo_database
        MONGODB_URI         = var.mongo_uri
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
      }
      dependencies = "{discord,db,apigateway}/**/*.go"
    }
    validate_webhook = {
      handler_name = "facebook_validate_webhook"
      handler_path = format("%s/cmd/lambda/facebook/validate_webhook", path.root)
      environment_variables = {
        APP_SECRET                           = var.facebook_app_secret
        ACCESS_TOKEN                         = var.facebook_access_token
        FACEBOOK_WEBHOOK_VERIFICATION_STRING = var.facebook_webhook_verification_string
        DISCORD_WEBHOOK_URL                  = var.discord_webhook_url
      }
      dependencies = "{discord,db,sqswrapper,apigateway}/**/*.go"
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
      dependencies = "{discord,api,db,stdconversation,apigateway}/**/*.go"
    }
    patch_conversation = {
      handler_name = "facebook_patch_conversation"
      handler_path = format("%s/cmd/lambda/facebook/patch_conversation", path.root)
      environment_variables = {
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
      dependencies = "{discord,api,db,stdconversation,apigateway}/**/*.go"
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
      dependencies = "{discord,api,db,stdmessage,apigateway}/**/*.go"
    }
    post_message = {
      handler_name = "facebook_post_message"
      handler_path = format("%s/cmd/lambda/facebook/post_message", path.root)
      environment_variables = {
        MONGODB_DATABASE    = var.mongo_database
        MONGODB_URI         = var.mongo_uri
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
      }
      dependencies = "{discord,api,db,external_api,stdmessage,apigateway}/**/*.go"
    }
    standardize_webhook = {
      handler_name = "facebook_standardize_webhook"
      handler_path = format("%s/cmd/lambda/facebook/standardize_webhook", path.root)
      environment_variables = {
        DISCORD_WEBHOOK_URL = var.discord_webhook_url
        MONGODB_URI         = var.mongo_uri
        MONGODB_DATABASE    = var.mongo_database
      }
      dependencies = "{discord,db,snswrapper,stdmessage,external_api,storage}/**/*.go"
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
      dependencies = "{discord,external_api,stdconversation,stdmessage,db}/**/*.go"
    }
  }
  method_integrations = {
    get_page_id = {
      method  = "GET"
      handler = "get_page_id"
    }
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


