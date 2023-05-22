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

const uri = "mongodb+srv://paff:thisispassword@botiolivechat.qsb7kv4.mongodb.net/?retryWrites=true&w=majority"

func QueryMessages(pageID string, conversationID string, outputMessage *OutputMessage) error {
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
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Println("Error Pinging DB : ", err)
		return err
	}
	log.Println("Successfully connect to MongoDB ", time.Since(start))
	discordLog(fmt.Sprintf("Successfully connect to mongo Elasped : %v", time.Since(start)))

	// start query
	coll := client.Database("BotioLivechat").Collection("facebook_messages")
	filter := bson.D{{Key: "pageID", Value: pageID}, {Key: "conversationID", Value: conversationID}}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		discordLog(fmt.Sprintf("Error query with filter pageID:%v conversationID:%v Error : %v", pageID, conversationID, err))
		return err
	}
	err = cur.All(ctx, &outputMessage.Messages)
	if err != nil {
		discordLog(fmt.Sprintf("Error retrieving doc in cur.ALL : %v", err))
		return err
	}

	return nil
}
