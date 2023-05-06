provider "aws" {
  region = "ap-southeast-1"
}

resource "aws_api_gateway_rest_api" "botio_rest_api" {
  name        = "botio_rest_api"
  description = "Api endpoint for interacting with botio"
}

resource "aws_api_gateway_resource" "shop_id" {
  rest_api_id = aws_api_gateway_rest_api.botio_rest_api.id
  parent_id   = aws_api_gateway_rest_api.botio_rest_api.root_resource_id
  path_part   = "{shop_id}"
}
