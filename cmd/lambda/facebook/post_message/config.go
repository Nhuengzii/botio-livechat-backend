package main

import "github.com/Nhuengzii/botio-livechat-backend/livechat"

type config struct {
	DiscordWebhookURL       string
	FacebookPageAccessToken string
	DbClient                livechat.DBClient
}
