package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/cache/redis"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/websocketwrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	websocket_api_id := os.Getenv("WEBSOCKET_API_ID")
	cacheClient := redis.NewClient(addr, password)
	websocketClient := websocketwrapper.NewClient(ctx, fmt.Sprintf("https://%s.execute-api.ap-southeast-1.amazonaws.com/dev", websocket_api_id))
	if websocketClient == nil {
		fmt.Println("websocketClient is nil")
	}
	c := Config{
		cacheClient:       cacheClient,
		webSocketClient:   websocketClient,
		discordWebhookURL: webhookURL,
	}
	lambda.Start(c.Handler)
}

type ReceivedMessage struct {
	Message string `json:"Message"`
}

type WebsocketMessage struct {
	Action string                `json:"action"`
	Data   stdmessage.StdMessage `json:"message"`
}

var (
	errUnmarshalReceivedBody    = errors.New("error json unmarshal receive body")
	errUnmarshalReceivedMessage = errors.New("error json unmarshal receive message")
)

func (c *Config) Handler(ctx context.Context, sqsEvent events.SQSEvent) (events.APIGatewayProxyResponse, error) {
	var receiveBody ReceivedMessage
	var receiveMessage stdmessage.StdMessage
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &receiveBody)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), errUnmarshalReceivedBody
		}
		err = json.Unmarshal([]byte(receiveBody.Message), &receiveMessage)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), errUnmarshalReceivedMessage
		}

		webscoketMessage := WebsocketMessage{
			Action: "relay",
			Data:   receiveMessage,
		}

		jsonMessage, err := json.Marshal(webscoketMessage)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
		connections, err := c.cacheClient.GetShopConnections(ctx, receiveMessage.ShopID)
		if err != nil {
			return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
		}
		for _, connectionID := range connections {
			err = c.webSocketClient.Send(ctx, connectionID, string(jsonMessage))
			if err != nil {
				discord.Log(c.discordWebhookURL, fmt.Sprintf("Error send message to connectionID: %s", connectionID))
			}
		}

	}
	return apigateway.NewProxyResponse(200, "OK", "*"), nil
}
