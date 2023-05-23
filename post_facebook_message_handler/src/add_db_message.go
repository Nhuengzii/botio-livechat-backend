package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://paff:thisispassword@botiolivechat.qsb7kv4.mongodb.net/?retryWrites=true&w=majority"

func AddDBMessage() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	// connect mongo
	client, err := ConnectMongo(ctx)
	if err != nil {
		return err
	}

	return nil
}

func ConnectMongo(ctx context.Context) (*mongo.Client, error) {
	start := time.Now()
	opts := options.Client().ApplyURI(uri)

	// create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Println("Error connecting to mongo atlas : ", err)
		return nil, err
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
		return nil, err
	}
	log.Println("Successfully connect to MongoDB ", time.Since(start))

	return client, nil
}
