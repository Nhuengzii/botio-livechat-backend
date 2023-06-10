package main

import (
	"context"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
)

func main() {
	dbClient, err := mongodb.NewClient(context.TODO(), &mongodb.Target{
		URI:                     os.Getenv("DATABASE_CONNECTION_URI"),
		Database:                "BotioLivechat",
		CollectionMessages:      "facebook_messages",
		CollectionConversations: "facebook_conversations",
	})
	if err != nil {
		return
	}
	c := config{
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		DbClient:          dbClient,
	}
}
