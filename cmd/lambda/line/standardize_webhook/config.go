package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/storage/amazons3"
)

type config struct {
	discordWebhookURL string
	snsTopicARN       string
	snsClient         livechat.PushNotificationClient
	dbClient          livechat.DBClient
	uploader          amazons3.Uploader
}
