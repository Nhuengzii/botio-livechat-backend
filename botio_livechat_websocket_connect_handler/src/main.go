package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("Got connect")
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

func Handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	shopId := request.QueryStringParameters["shopId"]
	discordLog(fmt.Sprint("Got connect: ", connectionID))
	my_ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-15520.c252.ap-southeast-1-1.ec2.cloud.redislabs.com:15520",
		Password: "dcesPhFIPwWrItb2yaNe5UT0sbhv9FJk",
	})

	err := rdb.Set(my_ctx, shopId+":"+connectionID, 1, 0).Err()
	if err != nil {
		discordLog(fmt.Sprint("Error setting current connection: ", err))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
