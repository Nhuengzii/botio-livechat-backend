package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/apirequest/sendmsgrequest"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/apiresponse/sendmsgresponse"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/db/mongodb"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external/fbrequest"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	errNoPSIDParam          = errors.New("err query string parameter psid not given")
	errNoPageIDPath         = errors.New("err path parameter parameters page_id not given")
	errNoConversationIDPath = errors.New("err path parameter conversation_id not given")
)

func (c *config) handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	psid, ok := request.QueryStringParameters["psid"]
	if !ok {
		discord.Log(c.DiscordWebhookURL, "err query string param psid was not given")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
		}, errNoPSIDParam
	}
	pageID, ok := request.PathParameters["page_id"]
	if !ok {
		discord.Log(c.DiscordWebhookURL, "err path param page_id was not given")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
		}, errNoPageIDPath
	}
	conversationID, ok := request.PathParameters["conversation_id"]
	if !ok {
		discord.Log(c.DiscordWebhookURL, "err path param conversation_id was not given")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
		}, errNoConversationIDPath
	}

	var requestMessage sendmsgrequest.Request
	err := json.Unmarshal([]byte(request.Body), &requestMessage)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	facebookRequest := fmtFbRequest(&requestMessage, pageID, psid)
	facebookResponse, err := fbrequest.RequestFacebookPostMessage(c.FacebookPageAccessToken, *facebookRequest, pageID, psid)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       "Bad Gateway",
		}, err
	}
	// map facebook response to api response
	response := sendmsgresponse.Response{
		RecipientID: facebookResponse.RecipientID,
		MessageID:   facebookResponse.MessageID,
		Timestamp:   facebookResponse.Timestamp,
	}

	jsonBodyByte, err := json.Marshal(response)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	err = c.updateDB(ctx, requestMessage, *facebookResponse, pageID, conversationID, psid)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       "Bad Gateway",
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBodyByte),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()
	dbClient, err := mongodb.NewClient(ctx, &mongodb.Target{
		URI:                     os.Getenv("DATABASE_CONNECTION_URI"),
		Database:                "BotioLivechat",
		CollectionMessages:      "facebook_messages",
		CollectionConversations: "facebook_conversations",
	})
	if err != nil {
		return
	}
	c := config{
		DiscordWebhookURL:       os.Getenv("DISCORD_WEBHOOK_URL"),
		FacebookPageAccessToken: os.Getenv("ACCESS_TOKEN"),
		DbClient:                dbClient,
	}
	defer func() {
		discord.Log(c.DiscordWebhookURL, "defer dbclient close")
		c.DbClient.Close(ctx)
	}()

	lambda.Start(c.handler)
}
