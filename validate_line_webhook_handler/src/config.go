package main

import "os"

var (
	discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	lineChannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	sqsQueueURL       = os.Getenv("SQS_QUEUE_URL")
)
