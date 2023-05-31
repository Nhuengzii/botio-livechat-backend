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

type TypingEvent struct {
	Action  string `json:"action"`
	Message struct {
		ConversationID string `json:"conversationID"`
		Platform       string `json:"platform"`
		Typing         bool   `json:"typing"`
	} `json:"message"`
}

func main() {
	fmt.Println("Hello")
	lambda.Start(Handler)
}

func sendMessage(svc *apigatewaymanagementapi.Client, connectionID string, message string) {
	input := &apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connectionID), Data: []byte(message)}
	_, err := svc.PostToConnection(context.Background(), input)
	if err != nil {
		discordLog(fmt.Sprint("Error sending message: ", err))
	}
}

func Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	endpoint := os.Getenv("WEBSOCKET_API_ENDPOINT")
	var typingEvent TypingEvent
	err := json.Unmarshal([]byte(request.Body), &typingEvent)
	if err != nil {
		discordLog(fmt.Sprint("Error Unmarshal message: ", err))
	}

	shopId := "1"
	my_ctx := context.Background()
	redis_addr := os.Getenv("REDIS_ACCESS_ADDR")
	redis_password := os.Getenv("REDIS_ACCESS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_password,
	})
	defer rdb.Close()
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
	for _, key := range keys {
		if key == shopId+":"+connectionID {
			continue
		}
		sendMessage(svc, strings.Split(key, ":")[1], request.Body)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
