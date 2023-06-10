package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	DiscordWebhookURL   string
	DbClient            livechat.DBClient
	FacebookAccessToken string
}
