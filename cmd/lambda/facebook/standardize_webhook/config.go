package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	DiscordWebhookURL       string
	SnsQueueURL             string
	SnsClient               livechat.PushNotificationClient
	FacebookPageAccessToken string
}
