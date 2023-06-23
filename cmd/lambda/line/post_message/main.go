package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/snswrapper"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errConversationNotExist = errors.New("err conversation ID does not exist")
	errPageNotExist         = errors.New("err page ID does not exist")
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/line/post_message.main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	pathParameters := req.PathParameters
	shopID := pathParameters["shop_id"]
	pageID := pathParameters["page_id"]
	conversationID := pathParameters["conversation_id"]
	err = c.dbClient.CheckConversationExists(ctx, conversationID)
	if err != nil {
		return apigateway.NewProxyResponse(404, errConversationNotExist.Error(), "*"), nil
	}
	page, err := c.dbClient.QueryLineAuthentication(ctx, pageID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return apigateway.NewProxyResponse(404, errPageNotExist.Error(), "*"), nil
		}
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	lineChannelAccessToken := page.AccessToken
	lineChannelSecret := page.Secret
	bot, err := linebot.New(lineChannelSecret, lineChannelAccessToken)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	var postMessageRequestBody postmessage.Request
	err = json.Unmarshal([]byte(req.Body), &postMessageRequestBody)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	err = c.handlePostMessageRequest(ctx, shopID, pageID, conversationID, bot, postMessageRequestBody)
	if err != nil {
		if errors.Is(err, errUnsupportedAttachmentType) {
			return apigateway.NewProxyResponse(400, err.Error(), "*"), nil
		}
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	return apigateway.NewProxyResponse(200, "OK", "*"), nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()
	var (
		mongodbURI        = os.Getenv("MONGODB_URI")
		mongodbDatabase   = os.Getenv("MONGODB_DATABASE")
		discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
		snsTopicARN       = os.Getenv("SNS_TOPIC_ARN")
		awsRegion         = os.Getenv("AWS_REGION")
	)
	dbClient, err := mongodb.NewClient(ctx, mongodb.Target{
		URI:                     mongodbURI,
		Database:                mongodbDatabase,
		CollectionConversations: "conversations",
		CollectionMessages:      "messages",
		CollectionShops:         "shops",
	})
	if err != nil {
		logMessage := "cmd/lambda/line/post_message.main.main: " + err.Error()
		discord.Log(discordWebhookURL, logMessage)
		log.Fatalln(logMessage)
	}
	defer dbClient.Close(ctx)
	c := &config{
		discordWebhookURL: discordWebhookURL,
		snsTopicARN:       snsTopicARN,
		snsClient:         snswrapper.NewClient(awsRegion),
		dbClient:          dbClient,
	}
	lambda.Start(c.handler)
}
