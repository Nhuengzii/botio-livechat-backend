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

func newDBclient(ctx context.Context) (*dbClient, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodbURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("newDBclient: %w", err)
	}
	return &dbClient{client: client}, nil
}

func (dbc *dbClient) close(ctx context.Context) error {
	err := dbc.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("dbClient.close: %w", err)
	}
	return nil
}

func (dbc *dbClient) insertConversation(ctx context.Context, c *botioConversation) error {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	_, err := coll.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("dbClient.insertConversation: %w", err)
	}
	return nil
}

func (dbc *dbClient) checkConversationExists(ctx context.Context, m *botioMessage) (bool, error) {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	filter := bson.D{{Key: "conversationID", Value: m.ConversationID}}
	c := &botioConversation{}
	if err := coll.FindOne(ctx, filter).Decode(c); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		} else {
			return false, fmt.Errorf("dbClient.checkConversationExists: %w", err)
		}
	}
	return true, nil
}

func (dbc *dbClient) updateConversationOfMessage(ctx context.Context, m *botioMessage) (err error) {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	filter := bson.D{{Key: "conversationID", Value: m.ConversationID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "lastActivity", Value: m.Message},
		{Key: "updatedTime", Value: m.Timestamp},
	}}}
	if _, err := coll.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("dbClient.updateConversation: %w", err)
	}
	return nil
}

func (dbc *dbClient) insertMessage(ctx context.Context, m *botioMessage) error {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineMessages)
	if _, err := coll.InsertOne(ctx, m); err != nil {
		return fmt.Errorf("dbClient.insertMessage: %w", err)
	}
	return nil
}
