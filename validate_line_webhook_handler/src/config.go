package main

import "os"

var discordWebhookURL, logToDiscordEnabled = os.LookupEnv("DISCORD_WEBHOOK_URL")
var lineChannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
var sqsQueueURL = os.Getenv("SQS_QUEUE_URL")
