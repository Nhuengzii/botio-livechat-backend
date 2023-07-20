package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/instagram/reqigsendmessage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoPSIDParam                      = errors.New("err query string parameter psid not given")
	errNoShopIDPath                     = errors.New("err path parameter shop_id not given")
	errNoPageIDPath                     = errors.New("err path parameter parameters page_id not given")
	errNoConversationIDPath             = errors.New("err path parameter conversation_id not given")
	errConversationNotExist             = errors.New("err conversation ID does not exist")
	errPageNotExist                     = errors.New("err page ID does not exist")
	errAttachmentTypeNotSupported       = errors.New("err attachment type given is not supported")
	errNoSrcFoundForBasicPayload        = errors.New("err this attachment type should not have an empty url")
	errNoPayloadFoundForTemplatePayload = errors.New("err this template attachment type should not have empty elements ")
	errSendingInstagramMessage          = errors.New("err sending instagram message check the body of the request")
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

	//**check parameters**//
	psid, ok := request.QueryStringParameters["psid"]
	if !ok {
		return apigateway.NewProxyResponse(400, errNoPSIDParam.Error(), "*"), nil
	}
	pageID, ok := request.PathParameters["page_id"]
	if !ok {
		return apigateway.NewProxyResponse(400, errNoPageIDPath.Error(), "*"), nil
	}
	conversationID := request.PathParameters["conversation_id"]
	if !ok {
		return apigateway.NewProxyResponse(400, errNoConversationIDPath.Error(), "*"), nil
	}
	err = c.dbClient.CheckConversationExists(ctx, conversationID)
	if err != nil {
		return apigateway.NewProxyResponse(404, errConversationNotExist.Error(), "*"), nil
	}
	//**finish checking parameters**//

	var requestMessage postmessage.Request
	err = json.Unmarshal([]byte(request.Body), &requestMessage)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	igCredentials, err := c.dbClient.QueryInstagramAuthentication(ctx, pageID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return apigateway.NewProxyResponse(404, errPageNotExist.Error(), "*"), nil
		}
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	igRequest, err := fmtIgRequest(&requestMessage, psid)
	if !ok {
		return apigateway.NewProxyResponse(400, err.Error(), "*"), nil
	}

	shop, err := c.dbClient.QueryShop(ctx, pageID)
	if err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return apigateway.NewProxyResponse(404, errPageNotExist.Error(), "*"), nil
		}
		return apigateway.NewProxyResponse(503, "Service Unavailable", "*"), nil
	}
	igResponse, err := reqigsendmessage.SendMessage(igCredentials.AccessToken, *igRequest, shop.FacebookPageID)
	if err != nil {
		return apigateway.NewProxyResponse(503, "Service Unavailable", "*"), nil
	}

	// map instagram response to api response
	resp := postmessage.Response{
		RecipientID: igResponse.RecipientID,
		MessageID:   igResponse.MessageID,
		Timestamp:   igResponse.Timestamp,
	}
	if resp.MessageID == "" || resp.RecipientID == "" {
		return apigateway.NewProxyResponse(500, errSendingInstagramMessage.Error(), "*"), nil
	}
	jsonBodyByte, err := json.Marshal(resp)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}

	return apigateway.NewProxyResponse(200, string(jsonBodyByte), "*"), nil
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
		CollectionTemplates:     "templates",
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
