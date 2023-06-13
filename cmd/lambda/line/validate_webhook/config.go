package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	discordWebhookURL string
	sqsQueueURL       string
	sqsClient         livechat.MessageQueueClient
	dbClient          livechat.DBClient
}
