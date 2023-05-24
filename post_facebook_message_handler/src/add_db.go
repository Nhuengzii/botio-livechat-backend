package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://paff:thisispassword@botiolivechat.qsb7kv4.mongodb.net/?retryWrites=true&w=majority"

var errAttachmentTypeNotSupport = errors.New("Attachment type not support")

func UpdateDB(pageID string, conversationID string, facebookResponse FacebookResponse, requestMessage RequestMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	// connect mongo
	client, err := ConnectMongo(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting from mongo atlas : ", err)
			return
		}
	}()

	err = AddDBMessage(ctx, client, pageID, conversationID, facebookResponse.MessageID, facebookResponse.Timestamp, requestMessage.Message, requestMessage.Attachment)
	err = AddDbConversation(ctx, client, conversationID, facebookResponse, requestMessage)
	if err != nil {
		return err
	}

	return nil
}

func AddDbConversation(ctx context.Context, client *mongo.Client, conversationID string, facebookResponse FacebookResponse, requestMessage RequestMessage) error {
	coll := client.Database("BotioLivechat").Collection("facebook_conversations")

	lastActivity, err := adminLastActivityFormat(requestMessage)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"updatedTime": facebookResponse.Timestamp, "lastActivity": lastActivity}}
	updateFilter := bson.D{{Key: "conversationID", Value: conversationID}}
	result, err := coll.UpdateOne(ctx, updateFilter, update)
	if err != nil {
		return err
	}
	discordLog(fmt.Sprintf("Updated a document; changed fields: %v\n", result.ModifiedCount))
	return nil
}

func adminLastActivityFormat(requestMessage RequestMessage) (string, error) {
	if requestMessage.Message != "" {
		return fmt.Sprintf("คุณ: %v", requestMessage.Message), nil
	} else if requestMessage.Attachment.AttachmentType == "image" {
		return "คุณส่งรูปภาพ", nil
	} else if requestMessage.Attachment.AttachmentType == "audio" {
		return "คุณส่งข้อความเสียง", nil
	} else if requestMessage.Attachment.AttachmentType == "video" {
		return "คุณส่งวิดีโอ", nil
	} else if requestMessage.Attachment.AttachmentType == "file" {
		return "คุณส่งไฟล์", nil
	} else if requestMessage.Attachment.AttachmentType == "template" {
		return "คุณส่งเทมเพลท", nil
	}
	return "", errAttachmentTypeNotSupport
}

func AddDBMessage(ctx context.Context, client *mongo.Client, pageID string, conversationID string, messageID string, timestamp int64, message string, attachment Attachment) error {
	coll := client.Database("BotioLivechat").Collection("facebook_messages")

	doc := StandardMessage{
		ShopID:         "1",
		Platform:       "Facebook",
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      messageID,
		Timestamp:      timestamp,
		Source: Source{
			UserID:   pageID, // botio user id?
			UserType: "Admin",
		},
		Message: message,
		Attachments: []Attachment{
			attachment,
		},
		ReplyTo: ReplyMessage{
			MessageId: "",
		},
	}

	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	log.Printf("Inserted a document with _id: %v\n", result.InsertedID)
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

	// ping
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Println("Error Pinging DB : ", err)
		return nil, err
	}
	log.Println("Successfully connect to MongoDB ", time.Since(start))

	return client, nil
}
