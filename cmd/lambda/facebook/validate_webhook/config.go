package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	discordWebhookURL                 string
	sqsQueueURL                       string
	facebookAppSecret                 string
	facebookWebhookVerificationString string
	sqsClient                         livechat.MessageQueueClient
}
