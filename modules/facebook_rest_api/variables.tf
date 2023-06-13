variable "platform" {
  type = string
}

variable "facebook_webhook_verification_string" {
  type = string
}

variable "mongo_uri" {
  type = string
}
variable "mongo_database" {
  type = string
}
variable "discord_webhook_url" {
  type = string
}
variable "rest_api_id" {
  type = string
}

variable "rest_api_execution_arn" {
  type = string
}

variable "facebook_access_token" {
  type = string
}

variable "facebook_app_secret" {
  type = string
}

variable "parent_id" {
  type = string
}

variable "relay_received_message_queue" {
  type = object({
    arn = string
    id  = string
  })
}
