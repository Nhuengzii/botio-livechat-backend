package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/go-redis/redis/v8"
)

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1108750713758175293/U96dYkOWsQYSYrCx6rCFPGvrJ7TY_tMMVmm5IWIAdsCM7ffi_Fa-W9Dfxt7dAd8WNYR2"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
func main() {
	fmt.Println("Hello, World!")
	lambda.Start(Handler)
}

type Attachment struct {
	Type    string `json:"type"`
	Payload struct {
		Src string `json:"src"`
	} `json:"payload"`
}

type Source struct {
	SourceID   string `json:"sourceID"`
	SourceType string `json:"sourceType"`
}

type Message struct {
	ConversationID string       `json:"conversationID"`
	MessageID      string       `json:"messageID"`
	Timestamp      int64        `json:"timeStamp"`
	Source         Source       `json:"source"`
	Message        string       `json:"message"`
	Attachments    []Attachment `json:"attachments"`
}

type BroadcastMessage struct {
	Action  string  `json:"action"`
	Message Message `json:"message"`
}

type IncommingMessage struct {
	Action  string  `json:"action"`
	Message Message `json:"message"`
}

func Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	// discordLog("Got broadcast")
	// Unmarshal the message
	endpoint := os.Getenv("WEBSOCKET_API_ENDPOINT")
	var message IncommingMessage
	// discordLog("Raw message: " + request.Body)
	err := json.Unmarshal([]byte(request.Body), &message)
	if err != nil {
		discordLog(fmt.Sprint("Error unmarshalling message: ", err))
	}
	// discordLog(fmt.Sprintf("Message: %+v", message))
	shopId := "1"
	my_ctx := context.Background()
	redis_addr := os.Getenv("REDIS_ACCESS_ADDR")
	redis_password := os.Getenv("REDIS_ACCESS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_password,
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		discordLog(fmt.Sprint("Error loading config: ", err))
	}
	svc := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.EndpointResolver = apigatewaymanagementapi.EndpointResolverFunc(func(region string, options apigatewaymanagementapi.EndpointResolverOptions) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})
	})
	keys, err := rdb.Keys(my_ctx, shopId+":*").Result()
	// discordLog(fmt.Sprintf("Keys: %+v", keys))
	if err != nil {
		discordLog(fmt.Sprint("Error getting keys: ", err))
	}
	var broadcastMessage BroadcastMessage
	broadcastMessage.Action = "broadcast"
	broadcastMessage.Message = message.Message
	for _, key := range keys {
		if key == shopId+":"+connectionID {
			continue
		}

		json_message, err := json.Marshal(broadcastMessage)
		if err != nil {
			discordLog(fmt.Sprint("Error marshalling message: ", err))
		}
		sendMessage(svc, strings.Split(key, ":")[1], string(json_message))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
