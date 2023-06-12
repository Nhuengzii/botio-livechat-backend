package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
)

type config struct {
	discordWebhookURL      string
	lineChannelSecret      string // TODO remove and get from some db
	lineChannelAccessToken string // TODO remove and get from some db
	dbClient               livechat.DBClient
}
