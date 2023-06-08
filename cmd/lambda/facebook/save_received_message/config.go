package main

import "github.com/Nhuengzii/botio-livechat-backend/internal/db"

type config struct {
	DiscordWebhookURL string
	DbClient          *db.Client
}
