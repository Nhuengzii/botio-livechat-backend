package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	discordLog("facebook get messages handler!!!!")
	// context
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	start := time.Now()

	// check path params
	pathParams := request.PathParameters
	// shopID := pathParams["shop_id"]
	pageID := pathParams["page_id"]
	conversationID := pathParams["conversation_id"]

	// connect to mongo
	var client *mongo.Client
	err := ConnectMongo(client, ctx)
	if err != nil {
		discordLog(fmt.Sprintf("error connect mongo : %v", err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
		}, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			discordLog(fmt.Sprintln("Error disconnecting from mongo atlas : ", err))
			return
		}
	}()

	var outputMessage OutputMessage
	err = QueryMessages(client, pageID, conversationID, &outputMessage)
	if err != nil {
		discordLog(fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
		}, err
	}

	jsonBodyByte, err := json.Marshal(outputMessage)
	jsonString := string(jsonBodyByte)
	if err != nil {
		discordLog(fmt.Sprint(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
		}, err
	}

	err = UpdateConversationIsRead(client, conversationID)
	if err != nil {
		discordLog(fmt.Sprintf("Error updating conversationDB isRead field : %v", err))
	}
	discordLog(fmt.Sprintf("Total Elasped time: %v", time.Since(start)))

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       jsonString,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func ConnectMongo(client *mongo.Client, ctx context.Context) error {
	start := time.Now()

	opts := options.Client().ApplyURI(uri)

	// create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		discordLog(fmt.Sprintln("Error connecting to mongo atlas : ", err))
		return err
	}

	// ping
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		discordLog(fmt.Sprintln("Error Pinging DB : ", err))
		return err
	}
	log.Println("Successfully connect to MongoDB ", time.Since(start))
	discordLog(fmt.Sprintf("Successfully connect to mongo Elasped : %v", time.Since(start)))

	return nil
}
