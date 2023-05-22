package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	discordLog(fmt.Sprint("Got disconnect events from: ", connectionID))
	my_ctx := context.Background()
	redis_addr := os.Getenv("REDIS_ACCESS_ADDR")
	redis_password := os.Getenv("REDIS_ACCESS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_password,
	})

	keys, err := rdb.Keys(my_ctx, "*:"+connectionID).Result()
	if err != nil {
		discordLog(fmt.Sprint("Error getting keys: ", err))
	}
	discordLog(fmt.Sprint("Keys: ", keys))

	if len(keys) > 0 {
		_, err := rdb.Del(my_ctx, keys[0]).Result()
		if err != nil {
			discordLog(fmt.Sprint("Error deleting connection: ", err))
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
