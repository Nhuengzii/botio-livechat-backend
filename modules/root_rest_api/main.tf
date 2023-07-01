terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

resource "aws_api_gateway_rest_api" "rest_api" {
  name = var.rest_api_name
}
