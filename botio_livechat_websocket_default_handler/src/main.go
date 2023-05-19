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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
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
	if err != nil {
		discordLog(fmt.Sprint("Error setting current connection: ", err))
	}
	endpoint := os.Getenv("WEBSOCKET_API_ENDPOINT")
	discordLog(fmt.Sprint("Endpoint: ", endpoint))
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
	input := &apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connectionID), Data: []byte("Hello this is from Default handler")}
	_, err = svc.PostToConnection(context.Background(), input)
	if err != nil {
		discordLog(fmt.Sprint("Error sending message: ", err))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
