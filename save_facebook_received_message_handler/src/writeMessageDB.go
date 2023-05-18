package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func WriteMessageDb(client *mongo.Client, record events.SQSMessage) error {
	var recieveBody RecievedAwsStruct
	var recieveMessage StandardMessage
	err := json.Unmarshal([]byte(record.Body), &recieveBody)
	if err != nil {
		return err
	}
	err = bson.UnmarshalExtJSON([]byte(recieveBody.Message), true, &recieveMessage)
	if err != nil {
		return err
	}
	log.Printf("%+v", recieveMessage)
	// check if need to create conversation
	err = ConversationCreate(client, recieveMessage)
	if err != nil {
		return err
	}
	coll := client.Database("BotioLivechat").Collection("facebook_messages")
	result, err := coll.InsertOne(context.TODO(), recieveMessage)
	if err != nil {
		return err
	}
	log.Printf("Inserted a document with _id: %v\n", result.InsertedID)
	return nil
}

type RecievedAwsStruct struct {
	Message string `json:"Message"`
}
