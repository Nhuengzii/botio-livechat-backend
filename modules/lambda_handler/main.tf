terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.1.0"
    }
  }
}

variable "handler_name" {
  type = string
}

variable "handler_path" {
  type = string
}

variable "role_arn" {
  type = string
}

variable "environment_variables" {
  type = map(string)
  default = {
    foo = "bar"
  }
}


resource "null_resource" "build_handler" {
  triggers = {
    entire_source_code_hash = sha1(join("", [for f in fileset(format("%s/%s", path.root, var.handler_path), "*.go") : filesha1(format("%s/%s/%s", path.root, var.handler_path, f))]))
  }
  provisioner "local-exec" {
    command = format("CGO_ENABLED=0 GOOS=linux go build -C %s -o main .", var.handler_path)
  }
}

data "archive_file" "handler_zip" {
  type        = "zip"
  source_file = format("%s/%s/main", path.root, var.handler_path)
  output_path = format("%s/%s/handler.zip", path.root, var.handler_path)
  depends_on  = [null_resource.build_handler]
}

resource "aws_lambda_function" "handler" {
  filename         = data.archive_file.handler_zip.output_path
  function_name    = var.handler_name
  handler          = "main"
  runtime          = "go1.x"
  role             = var.role_arn
  source_code_hash = data.archive_file.handler_zip.output_base64sha256
  environment {
    variables = var.environment_variables
  }
  depends_on = [data.archive_file.handler_zip]
}

output "lambda" {
  value = aws_lambda_function.handler
}
