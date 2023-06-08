package main

import (
	"os"
)

var (
	discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	sqsQueueURL       = os.Getenv("SQS_QUEUE_URL")
	facebookAppSecret = os.Getenv("FACEBOOK_APP_SECRET") // TODO to be removed and get from some db instead
)
