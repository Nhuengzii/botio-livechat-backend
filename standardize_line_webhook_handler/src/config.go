package main

import "os"

var (
	discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	snsTopicARN       = os.Getenv("SNS_TOPIC_ARN")
)
