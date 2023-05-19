package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func QueryConversations() error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	opts := options.Client().ApplyURI(uri)

	// create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Println("Error connecting to mongo atlas : ", err)
		return err
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
		return err
	}
	log.Println("Successfully connect to MongoDB ", time.Since(start))
	discordLog(fmt.Sprint("Successfully connect to mongo Elasped : %v", time.Since(start)))

	return nil
}
