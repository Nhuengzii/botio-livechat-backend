
module "websocket_api" {
  source                   = "./modules/websocket_api"
  redis_password           = var.redis_password
  redis_addr               = var.redis_addr
  discord_webhook_url      = var.discord_webhook_url
  websocket_api_stage_name = var.websocket_api_stage_name
}

output "websocket_api" {
  value = {
    id = module.websocket_api.websocket_api_id
  }
}
