resource "aws_api_gateway_resource" "line" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.shop_id.id
  path_part   = "line"
}

resource "aws_api_gateway_resource" "line_page_id" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.line.id
  path_part   = "{page_id}"
}

resource "aws_api_gateway_resource" "line_webhook" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_resource.line_page_id.id
  path_part   = "webhook"
}

resource "aws_api_gateway_method" "get_validate_line_webhook" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.line_webhook.id
  authorization = "NONE"
  http_method   = "GET"
}
resource "aws_api_gateway_method" "post_validate_line_webhook" {
  rest_api_id   = aws_api_gateway_rest_api.botio_rest_api.id
  resource_id   = aws_api_gateway_resource.line_webhook.id
  authorization = "NONE"
  http_method   = "POST"
}
