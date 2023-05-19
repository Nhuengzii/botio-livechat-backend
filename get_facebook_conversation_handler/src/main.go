package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://paff:thisispassword@botiolivechat.qsb7kv4.mongodb.net/?retryWrites=true&w=majority"

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) {
	discordLog("get_facebook_conversation handler!!!")
	start := time.Now()

	pathParams := request.PathParameters
	shopID := pathParams["shop_id"]
	pageID := pathParams["page_id"]
	conversationID := pathParams["conversation_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	opts := options.Client().ApplyURI(uri)

	// create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Println("Error connecting to mongo atlas : ", err)
		return
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting from mongo atlas : ", err)
			return
		}
	}()

	// ping
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Println("Error Pinging DB : ", err)
		return
	}
	log.Println("Successfully connect to MongoDB ", time.Since(start))
}
