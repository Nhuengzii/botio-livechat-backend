package main

import (
	"context"
	"errors"
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

func (dbc *dbClient) checkConversationExists(ctx context.Context, m *botioMessage) (bool, string, error) {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	filter := bson.D{{Key: "conversationID", Value: m.ConversationID}}
	c := &botioConversation{}
	if err := coll.FindOne(ctx, filter).Decode(c); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "", nil
		} else {
			return false, "", fmt.Errorf("dbClient.checkConversationExists: %w", err)
		}
	}
	return true, c.ConversationID, nil
}

func (dbc *dbClient) updateConversation(ctx context.Context, conversationID string, m *botioMessage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dbClient.updateConversation: %w", err)
		}
	}()
	if conversationID != m.ConversationID {
		return errors.New("conversationID mismatch")
	}
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	filter := bson.D{{Key: "conversationID", Value: conversationID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "lastActivity", Value: m.Message},
		{Key: "updatedTime", Value: m.Timestamp},
	}}}
	if _, err := coll.UpdateOne(ctx, filter, update); err != nil {
		return err
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
