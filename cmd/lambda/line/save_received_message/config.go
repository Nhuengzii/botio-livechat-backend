package main

import "github.com/Nhuengzii/botio-livechat-backend/livechat"

type config struct {
	discordWebhookURL                  string
	lineChannelAccessToken             string // TODO remove and get from db with shopID and pageID
	mongodbURI                         string
	mongodbDatabase                    string
	mongodbCollectionLineConversations string
	mongodbCollectionLineMessages      string
	dbClient                           livechat.DBClient
}
