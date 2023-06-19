variable "platform" {
  type = string
}

variable "rest_api_id" {
  type = string
}

variable "rest_api_execution_arn" {
  type = string
}

variable "parent_id" {
  type = string
}

variable "handlers" {
  type = map(object({
    handler_name          = string
    handler_path          = string
    environment_variables = map(string)
    dependencies          = string
  }))
}

variable "method_integrations" {
  type = map(object({
    method  = string
    handler = string
  }))
}

variable "relay_received_message_handler" {
  type = string
}



