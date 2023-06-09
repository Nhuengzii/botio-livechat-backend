package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/sqswrapper"
)

type config struct {
	DiscordWebhookURL                 string
	SqsQueueURL                       string
	FacebookAppSecret                 string
	FacebookWebhookVerificationString string
	SqsClient                         sqswrapper.Client
}
