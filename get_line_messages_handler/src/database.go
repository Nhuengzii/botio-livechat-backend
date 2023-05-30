package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbClient struct {
	client *mongo.Client
}

func newDBClient(ctx context.Context) (*dbClient, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodbURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("newDBClient: %w", err)
	}
	return &dbClient{client: client}, nil
}

func (dbc *dbClient) Close(ctx context.Context) error {
	err := dbc.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("dbClient.Close: %w", err)
	}
	return nil
}

func (dbc *dbClient) getMessagesInConversation(ctx context.Context, conversationID string) ([]botioMessage, error) {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineMessages)
	filter := bson.D{{Key: "conversationID", Value: conversationID}}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("dbClient.getMessagesInConversation: %w", err)
	}
	defer cur.Close(ctx)
	messages := []botioMessage{}
	if err := cur.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("dbClient.getMessagesInConversation: %w", err)
	}
	return messages, nil
}
