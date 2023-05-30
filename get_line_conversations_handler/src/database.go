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

func (dbc *dbClient) getConversationsOfPage(ctx context.Context, pageID string) ([]botioConversation, error) {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	filter := bson.D{{Key: "pageID", Value: pageID}}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("dbClient.getConversationsOfPage: %w", err)
	}
	defer cur.Close(ctx)
	conversations := []botioConversation{}
	if err := cur.All(ctx, &conversations); err != nil {
		return nil, fmt.Errorf("dbClient.getConversationsOfPage: %w", err)
	}
	return conversations, nil
}
