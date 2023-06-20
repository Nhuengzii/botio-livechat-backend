package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	discordWebhookURL string
	snsTopicARN       string
	snsClient         livechat.PushNotificationClient
	dbClient          livechat.DBClient
}
