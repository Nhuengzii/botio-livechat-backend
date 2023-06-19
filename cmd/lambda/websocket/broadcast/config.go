package main

import "github.com/Nhuengzii/botio-livechat-backend/livechat"

type Config struct {
	cacheClient       livechat.CacheClient
	webSocketClient   livechat.WebsocketClient
	discordWebhookURL string
}
