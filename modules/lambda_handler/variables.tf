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

