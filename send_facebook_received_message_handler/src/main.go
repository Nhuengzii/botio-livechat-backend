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
		log.Println("Error sending discord log: ", err)
	}
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	endpoint := os.Getenv("WEBSOCKET_API_ENDPOINT")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-15520.c252.ap-southeast-1-1.ec2.cloud.redislabs.com:15520",
		Password: "dcesPhFIPwWrItb2yaNe5UT0sbhv9FJk",
	})
	keys, _ := rdb.Keys(ctx, "1:*").Result()
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

	for _, message := range sqsEvent.Records {
		var standardMessage StandardMessage
		json.Unmarshal([]byte(message.Body), &standardMessage)
		for _, key := range keys {
			sendMessage(svc, strings.Split(key, ":")[1], standardMessage.Message)
		}
	}
}
