package main

import "github.com/Nhuengzii/botio-livechat-backend/internal/snswrapper"

type config struct {
	DiscordWebhookURL string
	SnsQueueURL       string
	SnsClient         snswrapper.Client
}
