package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNoConversations = errors.New("mongodb: no conversations")

type Client struct {
	client *mongo.Client
	Target
}

type Target struct {
	URI                     string
	Database                string
	CollectionConversations string
	CollectionMessages      string
}

func NewClient(ctx context.Context, target *Target) (*Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(target.URI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("db.NewClient: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) Close(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("db.Client.Close: %w", err)
	}
	return nil
}

func (c *Client) InsertConversation(ctx context.Context, conversation *stdconversation.StdConversation) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	_, err := coll.InsertOne(ctx, conversation)
	if err != nil {
		return fmt.Errorf("db.Client.InsertConversation: %w", err)
	}
	return nil
}

func (c *Client) InsertMessage(ctx context.Context, message *stdmessage.StdMessage) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	_, err := coll.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("db.Client.InsertMessage: %w", err)
	}
	return nil
}

func (c *Client) UpdateConversationOnNewMessage(ctx context.Context, message *stdmessage.StdMessage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("db.Client.UpdateConversationOnNewMessage: %w", err)
		}
	}()
	lastActivity, err := message.ToLastActivityString()
	if err != nil {
		return err
	}
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{{Key: "conversationID", Value: message.ConversationID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "lastActivity", Value: lastActivity},
		{Key: "updatedTime", Value: message.Timestamp},
	}}}
	var updatedConversation bson.D
	err = coll.FindOneAndUpdate(ctx, filter, update).Decode(updatedConversation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoConversations
		}
		return err
	}
	return nil
}

func (c *Client) UpdateConversationIsRead(ctx context.Context, conversationID string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("db.Client.UpdateConversationIsRead %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{{Key: "conversationID", Value: conversationID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "isRead", Value: true},
	}}}
	var updatedConversation bson.D
	err = coll.FindOneAndUpdate(ctx, filter, update).Decode(updatedConversation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoConversations
		}
		return err
	}
	return nil
}
