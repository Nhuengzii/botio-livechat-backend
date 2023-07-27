aws_access_key            = "aws_access_key"      #AWS account's access key. Can be retrieved via AWS's console.
aws_secret_key            = "aws_secret_key"      #AWS account's secret access key. Can be retrieved via AWS's console.
rest_api_stage_name       = "dev"                 #Deployment stage for REST api endpoints.
websocket_api_stage_name  = "dev"                 #Deployment stage for websocket endpoints.
mongo_uri                 = "mongodb_uri"         #MongoDB's connection URI.
mongo_database            = "mongodb_database"    #MongoDB's database name.
discord_webhook_url       = "discord_webhook_url" #Discord text channel's webhook URL for logging.
redis_addr                = "redis_addr"          #Redis database address.
redis_password            = "redis_password"      #Redis database password.
media_storage_bucket_name = "name"                #S3 bucket name for media storage.

facebook_app_secret                  = "facebook_app_secret"                  #Facebook app secret. Can be retrieved via facebook's developer console.
facebook_webhook_verification_string = "facebook_webhook_verification_string" #Set to the same value with facebook's developer console's webhook setup.

instagram_app_secret                  = "instagram_app_secret"                  #Instagram app secret. Can be retrieved via facebook's developer console.
instagram_webhook_verification_string = "instagram_webhook_verification_string" #Set to the same value with facebook's developer console's webhook setup.
