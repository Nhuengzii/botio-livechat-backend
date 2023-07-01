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

output "websocket_api" {
  value = {
    id = aws_apigatewayv2_api.botio_livechat_websocket.id
  }
}
