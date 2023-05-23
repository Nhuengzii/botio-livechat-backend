package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateConversationIsRead(client *mongo.Client, conversationID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	coll := client.Database("BotioLivechat").Collection("facebook_conversations")

	update := bson.M{"isRead": true}
	updateFilter := bson.D{{Key: "conversationID", Value: conversationID}}
	result, err := coll.UpdateOne(ctx, updateFilter, update)
	if err != nil {
		return err
	}
	discordLog(fmt.Sprintf("Updated a document; changed fields: %v\n", result.ModifiedCount))

	return nil
}
