package main

import (
	"context"
	"fmt"
	"os"
	// "time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/cache/redis"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	fmt.Println("connect handler")
	// ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Second)
	// defer cancel()
	ADDR := os.Getenv("REDIS_ADDR")
	PASSWORD := os.Getenv("REDIS_PASSWORD")
	cacheClient := redis.NewClient(ADDR, PASSWORD)
	c := config{
		cacheClient: cacheClient,
	}
	lambda.Start(c.handler)
}

func (c *config) handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	shopID := request.QueryStringParameters["shopID"]
	connectionID := request.RequestContext.ConnectionID
	err = c.cacheClient.SetShopConnection(ctx, shopID, connectionID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       "Bad Boy",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
