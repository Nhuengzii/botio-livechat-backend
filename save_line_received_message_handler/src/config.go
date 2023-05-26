package main

import "os"

var (
	discordWebhookURL                  = os.Getenv("DISCORD_WEBHOOK_URL")
	lineChannelAccessToken             = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	mongodbURI                         = os.Getenv("MONGODB_URI")
	mongodbDatabase                    = os.Getenv("MONGODB_DATABASE")
	mongodbCollectionLineConversations = os.Getenv("MONGODB_COLLECTION_LINE_CONVERSATIONS")
	mongodbCollectionLineMessages      = os.Getenv("MONGODB_COLLECTION_LINE_MESSAGES")
)
