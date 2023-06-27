package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/cache/redis"
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

type ReceiveMessage struct {
	Action  string                `json:"action"`
	Message stdmessage.StdMessage `json:"message"`
}

func (c *Config) Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	var receiveMessage ReceiveMessage
	err := json.Unmarshal([]byte(request.Body), &receiveMessage)
	if err != nil {
		fmt.Println(err)
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	shopID := receiveMessage.Message.ShopID
	connections, err := c.cacheClient.GetShopConnections(ctx, shopID)
	if err != nil {
		fmt.Println(err)
	}
	for _, connection := range connections {
		if connection != connectionID {
			err := c.webSocketClient.Send(ctx, connection, request.Body)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return apigateway.NewProxyResponse(200, "OK", "*"), nil
}
