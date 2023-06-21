package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoPSIDParam                      = errors.New("err query string parameter psid not given")
	errNoShopIDPath                     = errors.New("err path parameter shop_id not given")
	errNoPageIDPath                     = errors.New("err path parameter parameters page_id not given")
	errNoConversationIDPath             = errors.New("err path parameter conversation_id not given")
	errAttachmentTypeNotSupported       = errors.New("err attachment type given is not supported")
	errNoSrcFoundForBasicPayload        = errors.New("err this attachment type should not have an empty url")
	errNoPayloadFoundForTemplatePayload = errors.New("err this template attachment type should not have empty elements ")
	errSendingFacebookMessage           = errors.New("err sending facebook message check the body of the request")
)

const (
	templateButtonURLType  = "web_url"
	templateTypeGeneric    = "generic"
	attachmentTypeTemplate = "template"
)

func (c *config) handler(ctx context.Context, request events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			discord.Log(c.discordWebhookURL, fmt.Sprintln(err))
		}
	}()
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	var (
		mongodbURI        = os.Getenv("MONGODB_URI")
		mongodbDatabase   = os.Getenv("MONGODB_DATABASE")
		discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	)

	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     mongodbURI,
		Database:                mongodbDatabase,
		CollectionMessages:      "messages",
		CollectionConversations: "conversations",
		CollectionShops:         "shops",
	})
	c := config{
		discordWebhookURL: discordWebhookURL,
		dbClient:          dbClient,
	}
	if err != nil {
		discord.Log(c.discordWebhookURL, fmt.Sprintln(err))
		log.Fatalln(err)
	}
	defer func() {
		discord.Log(c.discordWebhookURL, "defer dbClient close")
		c.dbClient.Close(ctx)
	}()

	lambda.Start(c.handler)
}
