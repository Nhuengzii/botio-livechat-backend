package main

import "github.com/Nhuengzii/botio-livechat-backend/livechat"

type config struct {
	discordWebhookURL string
	dbClient          livechat.DBClient
}
