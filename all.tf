module "all_platform_rest_api" {
  source                 = "./modules/all_platform_rest_api"
  rest_api_id            = module.rest_api.id
  rest_api_execution_arn = module.rest_api.execution_arn
  parent_id              = module.shops.shop_id_resource_id
  get_conversations_handler = {
    handler_path = format("%s/cmd/lambda/all/get_conversations", path.root)
    handler_name = "all_platform_get_conversations"
    environment_variables = {
      DISCORD_WEBHOOK_URL = var.discord_webhook_url
      MONGODB_URI         = var.mongo_uri
      MONGODB_DATABASE    = var.mongo_database
    },
    dependencies = "{discord,db,stdconversation,stdmessage}/**/*.go"
  }
  get_all = {
    handler_path = format("%s/cmd/lambda/all/get_all", path.root)
    handler_name = "all_platform_get_all"
    environment_variables = {
      DISCORD_WEBHOOK_URL = var.discord_webhook_url
      MONGODB_URI         = var.mongo_uri
      MONGODB_DATABASE    = var.mongo_database
    },
    dependencies = "{discord,db,stdconversation,stdmessage}/**/*.go"
  }
}
