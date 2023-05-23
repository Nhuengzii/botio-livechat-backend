package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://paff:thisispassword@botiolivechat.qsb7kv4.mongodb.net/?retryWrites=true&w=majority"

func AddDBMessage(pageID string, conversationID string, messageID string, message string, attachment Attachment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	timestamp, err := getMessageCreatedTime(messageID)
	if err != nil {
		return err
	}
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

func getMessageCreatedTime(messageID string) (int64, error) {
	access_token := os.Getenv("ACCESS_TOKEN")
	getMessageURI := fmt.Sprintf("https://graph.facebook.com/v16.0/%v?access_token=%v",
		messageID, access_token)

	response, err := http.Get(getMessageURI)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var messageDataResponse MessageDataResponse
	err = json.NewDecoder(response.Body).Decode(&messageDataResponse)

	timestampDatetime, err := time.Parse("2006-01-02T15:04:05-0700", messageDataResponse.Timestamp)
	if err != nil {
		return 0, err
	}
	return timestampDatetime.Unix(), nil
}

type MessageDataResponse struct {
	ID        string `json:"id"`
	Timestamp string `json:"created_time"`
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
