package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var errAttachmentTypeNotSupport = errors.New("Attachment type not support")

func ConversationCreate(client *mongo.Client, recieveMessage StandardMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	// connect to DB
	coll := client.Database("BotioLivechat").Collection("facebook_conversations")
	filter := bson.D{{Key: "conversationID", Value: recieveMessage.ConversationID}}
	result := coll.FindOne(ctx, filter)

	if result.Err() == mongo.ErrNoDocuments {
		discordLog(fmt.Sprint("No Conversation need to create one"))
		// get userProfile
		userProfile, err := RequestFacebookUserProfile(recieveMessage.Source.UserID)
		if err != nil {
			return err
		}
		doc := Conversation{
			ShopID:         recieveMessage.ShopID,
			PageID:         recieveMessage.PageID,
			ConversationID: recieveMessage.ConversationID,
			ConversationPic: Payload{
				Src: userProfile.ProfilePic,
			},
			UpdatedTime: recieveMessage.Timestamp,
			Participants: []Participant{
				{
					UserID:   recieveMessage.Source.UserID,
					Username: userProfile.Name,
					ProfilePic: Payload{
						Src: userProfile.ProfilePic,
					},
				},
			},
			LastActivity: recieveMessage.Message,
			IsRead:       false,
		}
		log.Printf("Doc : %+v", doc)
		result, err := coll.InsertOne(ctx, doc)
		if err != nil {
			return err
		}
		log.Printf("Inserted a document with _id: %v\n", result.InsertedID)
	} else { // update conversation
		var update primitive.M
		lastActivity, err := LastActivityFormat(recieveMessage)
		if err != nil {
			return err
		}
		update = bson.M{"$set": bson.M{"updatedTime": recieveMessage.Timestamp, "lastActivity": lastActivity, "isRead": false}}
		updateFilter := bson.D{{Key: "conversationID", Value: recieveMessage.ConversationID}}
		result, err := coll.UpdateOne(ctx, updateFilter, update)
		if err != nil {
			return err
		}
		log.Printf("Updated a document; changed fields: %v\n", result.ModifiedCount)
		// discordLog(fmt.Sprintf("Updated a document; changed fields: %v\n", result.ModifiedCount))
	}

	return nil
}

func LastActivityFormat(recieveMessage StandardMessage) (string, error) {
	if recieveMessage.Message != "" {
		return fmt.Sprintf("%v", recieveMessage.Message), nil
	}
	attachment := recieveMessage.Attachments[0]
	if attachment.AttachmentType == "image" {
		return "ส่งรูปภาพ", nil
	} else if attachment.AttachmentType == "audio" {
		return "ส่งข้อความเสียง", nil
	} else if attachment.AttachmentType == "video" {
		return "ส่งวิดีโอ", nil
	} else if attachment.AttachmentType == "file" {
		return "ส่งไฟล์", nil
	} else if attachment.AttachmentType == "template" {
		return "ส่งเทมเพลท", nil
	}

	return "", errAttachmentTypeNotSupport
}

type Conversation struct {
	ShopID          string        `bson:"shopID"`
	PageID          string        `bson:"pageID"`
	ConversationID  string        `bson:"conversationID"`
	ConversationPic Payload       `bson:"conversationPic"`
	UpdatedTime     int64         `bson:"updatedTime"`
	Participants    []Participant `bson:"participants"`
	LastActivity    string        `bson:"lastActivity"`
	IsRead          bool          `bson:"isRead"`
}

type Participant struct {
	UserID     string  `bson:"userID"`
	Username   string  `bson:"username"`
	ProfilePic Payload `bson:"profilePic"`
}
type Payload struct {
	Src string `bson:"src"`
}
