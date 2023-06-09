package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
)

type config struct {
	DiscordWebhookURL   string
	DbClient            *mongodb.Client
	FacebookAccessToken string
}
