package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/request/sendmsgrequest"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
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

	var requestMessage sendmsgrequest.RequestMessage
	err := json.Unmarshal([]byte(request.Body), &requestMessage)
	if err != nil {
		discord.Log(c.DiscordWebhookURL, fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	// request external facebook post
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	c := config{
		DiscordWebhookURL:       os.Getenv("DISCORD_WEBHOOK_URL"),
		FacebookPageAccessToken: os.Getenv("ACCESS_TOKEN"),
	}

	lambda.Start(c.handler)
}
