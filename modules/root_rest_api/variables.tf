variable "rest_api_name" {
  type = string
}

variable "get_upload_url_handler" {
  type = object({
    handler_name          = string
    handler_path          = string
    environment_variables = map(string)
    dependencies          = string
  })
}

variable "s3_bucket_arn" {
  type = string
}
