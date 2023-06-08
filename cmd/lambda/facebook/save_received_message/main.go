package main

import (
	"context"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/internal/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Lambda struct {
	config
}

func (l Lambda) handler(ctx context.Context, event events.SQSEvent) {
}

func main() {
	dbClient, err := db.NewClient(context.TODO(), &db.Target{
		URI:                     os.Getenv("DATABASE_CONNECTION_URI"),
		Database:                "BotioLivechat",
		CollectionMessages:      "facebook_messages",
		CollectionConversations: "facebook_conversations",
	})
	if err != nil {
		return
	}
	l := Lambda{
		config: config{
			DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
			DbClient:          dbClient,
		},
	}
	lambda.Start(l.handler)
}
