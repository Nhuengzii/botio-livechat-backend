package main

import "github.com/Nhuengzii/botio-livechat-backend/livechat/snswrapper"

type config struct {
	DiscordWebhookURL       string
	SnsQueueURL             string
	SnsClient               snswrapper.Client
	FacebookPageAccessToken string
}
