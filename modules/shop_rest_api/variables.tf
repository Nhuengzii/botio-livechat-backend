variable "rest_api_id" {
  type = string
}

variable "parent_resource_id" {
  type = string
}

variable "rest_api_execution_arn" {
  type = string
}

variable "handlers" {
  type = map(object({
    handler_name          = string
    handler_path          = string
    environment_variables = map(string)
  }))
}
