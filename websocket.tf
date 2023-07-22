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


module "websocket_api" {
  source                   = "./modules/websocket_api"
  redis_password           = var.redis_password
  redis_addr               = var.redis_addr
  discord_webhook_url      = var.discord_webhook_url
  websocket_api_stage_name = var.websocket_api_stage_name
}

output "websocket_api" {
  value = {
    id = aws_apigatewayv2_api.botio_livechat_websocket.id
  }
}
