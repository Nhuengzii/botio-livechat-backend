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

func discordLog(content string) {
	webhookURL := "https://discord.com/api/webhooks/1108750713758175293/U96dYkOWsQYSYrCx6rCFPGvrJ7TY_tMMVmm5IWIAdsCM7ffi_Fa-W9Dfxt7dAd8WNYR2"
	payload := map[string]string{"content": content}
	json_payload, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(json_payload))
	if err != nil {
		log.Println("Error  sending discord log: ", err)
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

func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	rdb := getRedisClient()
	svc := getSVCClient()
	keys, _ := rdb.Keys(ctx, "1:*").Result()

	for _, message := range sqsEvent.Records {
		// discordLog(fmt.Sprint("Received message: ", message.Body))
		var standardMessage StandardMessage

		json.Unmarshal([]byte(message.Body), &standardMessage)
		// discordLog(fmt.Sprint("Unmarshalled message: ", standardMessage.Message))
		websocketMessage := WebsocketMessage{
			Action:  "userMessage",
			Message: standardMessage.Message,
		}
		json_websocketMessage, _ := json.Marshal(websocketMessage)
		// discordLog(fmt.Sprint("Sending message: ", string(json_websocketMessage)))
		for _, key := range keys {
			sendMessage(svc, strings.Split(key, ":")[1], string(json_websocketMessage))
		}
	}
}
