package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const uri = "mongodb+srv://paff:thisispassword@botiolivechat.qsb7kv4.mongodb.net/?retryWrites=true&w=majority"

func QueryMessages(client *mongo.Client, pageID string, conversationID string, outputMessage *OutputMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

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
