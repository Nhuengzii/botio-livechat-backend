package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/mongo"
)

func WriteMessageDb(client *mongo.Client, record events.SQSMessage) error {
	var recieveBody RecievedAwsStruct
	var recieveMessage StandardMessage
	err := json.Unmarshal([]byte(record.Body), &recieveBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(recieveBody.Message), &recieveMessage)
	if err != nil {
		return err
	}
	log.Printf("%+v", recieveMessage)
	doc := StandardMessage{
		ShopID:         recieveMessage.ShopID,
		Platform:       "facebook",
		PageID:         recieveMessage.PageID,
		ConversationID: recieveMessage.ConversationID,
		MessageID:      recieveMessage.MessageID,
		Timestamp:      recieveMessage.Timestamp,
		Source:         recieveMessage.Source,
		Message:        recieveMessage.Message,
		Attachments:    recieveMessage.Attachments,
		ReplyTo:        recieveMessage.ReplyTo,
	}

	coll := client.Database("BotioLivechat").Collection("facebook_message")
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}
	log.Printf("Inserted a document with _id: %v\n", result.InsertedID)
	return nil
}

type RecievedAwsStruct struct {
	Message string `json:"Message"`
}
