variable "facebook_access_token" {
  type = string
}

variable "facebook_app_secret" {
  type = string
}

variable "facebook_webhook_verification_string" {
  type = string
}

variable "redis_addr" {
  type = string
}

variable "redis_password" {
  type = string
}

variable "line_channel_secret" {
  type = string
}

variable "line_channel_access_token" {
  type = string
}

variable "discord_webhook_url" {
  type = string
}

variable "mongo_uri" {
  type = string
}

variable "mongo_database" {
  type = string
}

variable "mongo_collection_line_conversations" {
  type = string
}

variable "mongo_collection_line_messages" {
  type = string
}

variable "instagram_access_token" {
  type = string
}

variable "instagram_app_secret" {
  type = string
}

variable "instagram_webhook_verification_string" {
  type = string
}

variable "s3_bucket_name" {
  type = string
}
variable "aws_access_key" {
  type = string
}

variable "aws_secret_key" {
  type = string
}

variable "rest_api_stage_name" {
  type    = string
  default = "dev"
}

variable "websocket_api_stage_name" {
  type    = string
  default = "dev"
}
