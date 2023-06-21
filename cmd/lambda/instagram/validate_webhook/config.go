package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	discordWebhookURL                  string
	sqsQueueURL                        string
	instagramAppSecret                 string
	instagramWebhookVerificationString string
	sqsClient                          livechat.MessageQueueClient
}
