package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
				Src: "PlaceHolder",
			},
			UpdatedTime: recieveMessage.Timestamp,
			Participants: []Participant{
				{
					UserID:   recieveMessage.Source.UserID,
					Username: "PlaceHolder",
					ProfilePic: Payload{
						Src: "PlaceHolder",
					},
				},
			},
			LastActivity: recieveMessage.Message,
		}
		log.Printf("Doc : %+v", doc)
		result, err := coll.InsertOne(ctx, doc)
		if err != nil {
			return err
		}
		log.Printf("Inserted a document with _id: %v\n", result.InsertedID)
	} else {
		var update primitive.M
		if recieveMessage.Message != "" {
			update = bson.M{"$set": bson.M{"updatedTime": recieveMessage.Timestamp, "lastActivity": recieveMessage.Message}}
		} else {
			attachType := recieveMessage.Attachments[0].AttachmentType
			update = bson.M{"$set": bson.M{"updatedTime": recieveMessage.Timestamp, "lastActivity": fmt.Sprintf("send a %v", attachType)}}
		}
		updateFilter := bson.D{{Key: "conversationID", Value: recieveMessage.ConversationID}}
		result, err := coll.UpdateOne(ctx, updateFilter, update)
		if err != nil {
			return err
		}
		discordLog(fmt.Sprintf("Updated a document; changed fields: %v\n", result.ModifiedCount))
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
