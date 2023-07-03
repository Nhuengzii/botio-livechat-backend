package main

import "github.com/Nhuengzii/botio-livechat-backend/livechat"

type config struct {
	discordWebhookURL string
	awsRegion         string
	storageClient     livechat.StorageClient
}
