package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	discordWebhookURL string
	sqsQueueURL       string
	lineChannelSecret string // TODO to be removed and get from db with shopID and pageID
	sqsClient         livechat.MessageQueueClient
}
