package main

import "os"

var (
	discordWebhookURL             = os.Getenv("DISCORD_WEBHOOK_URL")
	mongodbURI                    = os.Getenv("MONGODB_URI")
	mongodbDatabase               = os.Getenv("MONGODB_DATABASE")
	mongodbCollectionLineMessages = os.Getenv("MONGODB_COLLECTION_LINE_MESSAGES")
)
