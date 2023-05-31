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

func main() {
	fmt.Println("Hello Me")
	lambda.Start(Handler)
}
func sendMessage(svc *apigatewaymanagementapi.Client, connectionID string, message string) {
	input := &apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connectionID), Data: []byte(message)}
	_, err := svc.PostToConnection(context.Background(), input)
	if err != nil {
		discordLog(fmt.Sprint("Error sending message: ", err))
	}
}

type LineReceivedMessage struct {
	ShopID         string `json:"shopID"`
	Platform       string `json:"platform"`
	PageID         string `json:"pageID"`
	ConversationID string `json:"conversationID"`
	MessageID      string `json:"messageID"`
	Timestamp      int64  `json:"timestamp"`
	Source         struct {
		UserID   string `json:"userID"`
		UserType string `json:"userType"`
	}
	Message     string `json:"message"`
	Attachments []struct {
		AttachmentType string `json:"attachmentType"`
		Payload        struct {
			Src string `json:"src"`
		}
	}
}

func getRedisClient() *redis.Client {
	redis_addr := os.Getenv("REDIS_ACCESS_ADDR")
	redis_password := os.Getenv("REDIS_ACCESS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_password,
	})
	return rdb
}

func getSVCClient() *apigatewaymanagementapi.Client {
	endpoint := os.Getenv("WEBSOCKET_API_ENDPOINT")
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
	return svc
}

type EventRecordBody struct {
	Message string `json:"Message"`
}
type WebsocketMessage struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1108750713758175293/U96dYkOWsQYSYrCx6rCFPGvrJ7TY_tMMVmm5IWIAdsCM7ffi_Fa-W9Dfxt7dAd8WNYR2"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error sending discord log: ", err)
	}
}
func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	rdb := getRedisClient()
	svc := getSVCClient()
	keys, err := rdb.Keys(ctx, "1:*").Result()
	if err != nil {
		discordLog(fmt.Sprint("Error getting keys from redis: ", err))
	}
	defer rdb.Close()
	for _, message := range sqsEvent.Records {
		var eventRec EventRecordBody
		err := json.Unmarshal([]byte(message.Body), &eventRec)
		if err != nil {
			discordLog(fmt.Sprint("Error unmarshalling message in sendLineReceivedMessage: ", err))
		}
		var lineRM LineReceivedMessage
		err = json.Unmarshal([]byte(eventRec.Message), &lineRM)
		if err != nil {
			discordLog(fmt.Sprint("Error unmarshalling message in sendLineReceivedMessage: ", err))
		}
		var wsMessage WebsocketMessage
		wsMessage.Action = "userMessage"
		wsMessage.Message = eventRec.Message
		wsMessageJSON, _ := json.Marshal(wsMessage)
		for _, key := range keys {
			sendMessage(svc, strings.Split(key, ":")[1], string(wsMessageJSON))
		}
	}
}
