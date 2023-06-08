package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/internal/sqswrapper"
)

type config struct {
	DiscordWebhookURL                 string
	SqsQueueURL                       string
	FacebookAppSecret                 string
	FacebookWebhookVerificationString string
	SqsClient                         sqswrapper.Client
}
