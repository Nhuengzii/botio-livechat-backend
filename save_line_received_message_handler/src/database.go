package main

import (
	"context"
	"errors"

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
		return nil, &newDBclientError{
			message: "couldn't connect to mongodb",
			err:     err,
		}
	}
	return &dbClient{client: client}, nil
}

func (dbc *dbClient) close(ctx context.Context) error {
	err := dbc.client.Disconnect(ctx)
	if err != nil {
		return &closeDBclientError{
			message: "couldn't disconnect from mongodb",
			err:     err,
		}
	}
	return nil
}

func (dbc *dbClient) insertConversation(ctx context.Context, c *botioConversation) error {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	_, err := coll.InsertOne(ctx, c)
	if err != nil {
		return &insertConversationError{
			message: "couldn't inserst conversation into mongodb",
			err:     err,
		}
	}
	return nil
}

func (dbc *dbClient) checkConversationExists(ctx context.Context, m *botioMessage) (string, error) {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	filter := bson.D{{Key: "conversationID", Value: m.ConversationID}}
	c := &botioConversation{}
	err := coll.FindOne(ctx, filter).Decode(c)
	if err != nil {
		return "", &checkConversationExistsError{
			message: "couldn't check if conversation exists in mongodb",
			err:     err,
		}
	}
	return c.ConversationID, nil
}

func (dbc *dbClient) updateConversation(ctx context.Context, conversationID string, m *botioMessage) error {
	if conversationID != m.ConversationID {
		return &updateConversationError{
			message: "couldn't update conversation in mongodb",
			err:     errors.New("conversationID mismatch"),
		}
	}
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineConversations)
	filter := bson.D{{Key: "conversationID", Value: conversationID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "lastActivity", Value: m.Message},
		{Key: "updatedTime", Value: m.Timestamp},
	}}}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return &updateConversationError{
			message: "couldn't update conversation in mongodb",
			err:     err,
		}
	}
	return nil
}

func (dbc *dbClient) insertMessage(ctx context.Context, m *botioMessage) error {
	coll := dbc.client.Database(mongodbDatabase).Collection(mongodbCollectionLineMessages)
	_, err := coll.InsertOne(ctx, m)
	if err != nil {
		return &insertMessageError{
			message: "couldn't insert message into mongodb",
			err:     err,
		}
	}
	return nil
}

type newDBclientError struct {
	message string
	err     error
}
type closeDBclientError struct {
	message string
	err     error
}
type insertConversationError struct {
	message string
	err     error
}
type checkConversationExistsError struct {
	message string
	err     error
}
type updateConversationError struct {
	message string
	err     error
}
type insertMessageError struct {
	message string
	err     error
}

func (e *newDBclientError) Error() string {
	return e.message + ": " + e.err.Error()
}
func (e *closeDBclientError) Error() string {
	return e.message + ": " + e.err.Error()
}
func (e *insertConversationError) Error() string {
	return e.message + ": " + e.err.Error()
}
func (e *checkConversationExistsError) Error() string {
	return e.message + ": " + e.err.Error()
}
func (e *updateConversationError) Error() string {
	return e.message + ": " + e.err.Error()
}
func (e *insertMessageError) Error() string {
	return e.message + ": " + e.err.Error()
}
