package main

import (
	"context"
	"fmt"
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
	lambda.Start(handle)
}

func handle(ctx context.Context, sqsEvent events.SQSEvent) {
	start := time.Now()
	log.Println("facebook database  handler")
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
			log.Println("Error  disconnecting from mongo atlas : ", err)
			return
		}
	}()

	// ping
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Println("Error Pinging DB : ", err)
		return
	}
	log.Println("Successfully connect to MongoDB ", time.Since(start))

	for _, record := range sqsEvent.Records {
		err := WriteMessageDb(client, record)
		if err != nil {
			discordLog(fmt.Sprintf("Error Inserting doc to DB : %v", err))
		}
	}

	log.Println("Elapsed End: ", time.Since(start))
	return
}
