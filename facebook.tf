resource "aws_api_gateway_resource" "facebook" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.shop_id.id
  path_part   = "facebook"
}

resource "aws_api_gateway_resource" "facebook_webhook" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.facebook.id
  path_part   = "webhook"
}

