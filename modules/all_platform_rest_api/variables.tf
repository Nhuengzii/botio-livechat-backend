variable "rest_api_id" {
  type = string
}

variable "rest_api_execution_arn" {
  type = string
}

variable "parent_id" {
  type = string
}

variable "get_conversations_handler" {
  type = object({
    handler_name          = string
    handler_path          = string
    environment_variables = map(string)
    dependencies          = string
  })
}

variable "get_all" {
  type = object({
    handler_name          = string
    handler_path          = string
    environment_variables = map(string)
    dependencies          = string
  })
}
