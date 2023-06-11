package main

import "github.com/Nhuengzii/botio-livechat-backend/livechat"

type config struct {
	discordWebhookURL                  string
	mongodbURI                         string
	mongodbDatabase                    string
	mongodbCollectionLineConversations string
	dbClient                           livechat.DBClient
}
