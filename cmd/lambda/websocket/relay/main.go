package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/cache/redis"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/websocketwrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type receivedMessage struct {
	Message string `json:"Message"`
}

var (
	errUnmarshalReceivedBody    = errors.New("Error json unmarshal recieve body")
	errUnmarshalReceivedMessage = errors.New("Error json unmarshal recieve message")
)

func main() {
	fmt.Println("relay handler")

	ADDR := os.Getenv("REDIS_ADDR")
	PASSWORD := os.Getenv("REDIS_PASSWORD")
	WEBSOCKET_API_ENDPOINT := os.Getenv("WEBSOCKET_API_ENDPOINT")
	cacheClient := redis.NewClient(ADDR, PASSWORD)
	websocketClient := websocketwrapper.NewClient(WEBSOCKET_API_ENDPOINT)
	c := config{
		cacheClient:     cacheClient,
		websocketClient: websocketClient,
	}
	lambda.Start(c.handler)
}

type output struct {
	Action  string                `json:"action"`
	Message stdmessage.StdMessage `json:"message"`
}

func (c *config) handler(ctx context.Context, sqsEvent events.SQSEvent) {
	fmt.Println("Start relay handler")
	var receiveMessage stdmessage.StdMessage
	var shopID string
	var keys []string
	var connectionID string
	var jsonMessage []byte
	for _, record := range sqsEvent.Records {
		err := json.Unmarshal([]byte(record.Body), &receiveMessage)
		if err != nil {
			fmt.Println(errUnmarshalReceivedBody)
			return
		}
		fmt.Printf("Receive message: %v\n", receiveMessage)
		shopID = receiveMessage.ShopID
		keys, err = c.cacheClient.GetShopConnections(ctx, shopID)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, key := range keys {
			connectionID = strings.Split(key, ":")[1]
			jsonMessage, err = json.Marshal(output{Action: "relay", Message: receiveMessage})
			if err != nil {
				fmt.Println(err)
			}
			err = c.websocketClient.Send(ctx, connectionID, string(jsonMessage))
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}
