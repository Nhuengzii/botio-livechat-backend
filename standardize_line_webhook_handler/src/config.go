package main

import "os"

var discordWebhookURL, logToDiscordEnabled = os.LookupEnv("DISCORD_WEBHOOK_URL")
var snsTopicARN = os.Getenv("SNS_TOPIC_ARN")
