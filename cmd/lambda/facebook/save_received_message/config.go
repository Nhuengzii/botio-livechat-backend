package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	discordWebhookUrl   string
	dbClient            livechat.DBClient
	facebookAccessToken string
}
