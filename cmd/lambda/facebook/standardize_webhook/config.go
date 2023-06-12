package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	DiscordWebhookURL       string
	SnsTopicARN             string
	SnsClient               livechat.PushNotificationClient
	FacebookPageAccessToken string
}
