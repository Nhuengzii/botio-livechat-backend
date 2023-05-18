package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConversationCreate(client *mongo.Client, recieveMessage StandardMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	coll := client.Database("BotioLivechat").Collection("facebook_conversations")
	filter := bson.D{{Key: "conversationID", Value: recieveMessage.ConversationID}}
	result := coll.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		discordLog(fmt.Sprint("No Conversation need to create one"))
		doc := Conversation{
			ShopID:         recieveMessage.ShopID,
			PageID:         recieveMessage.PageID,
			ConversationID: recieveMessage.ConversationID,
			ConversationPic: Payload{
				Src: recieveMessage.ConversationID,
			},
		}
		result, err := coll.InsertOne(ctx)
	}
	return nil
}

type Conversation struct {
	ShopID          string        `bson:"shopID"`
	PageID          string        `bson:"pageID"`
	ConversationID  string        `bson:"conversationID"`
	ConversationPic Payload       `bson:"conversationPic"`
	UpdatedTime     int64         `bson:"updatedTime"`
	Participants    []Participant `bson:"participants"`
	LastActivity    string        `bson:"lastActivity"`
}

type Participant struct {
	UserID     string  `bson:"userID"`
	Username   string  `bson:"username"`
	ProfilePic Payload `bson:"profilePic"`
}
type Payload struct {
	Src string `bson:"src"`
}
