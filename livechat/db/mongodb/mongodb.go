package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoConversations = errors.New("mongodb: no conversations")
	ErrNoMessages      = errors.New("mongodb: no messages")
)

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

func (c *Client) InsertConversation(ctx context.Context, conversation *livechat.StdConversation) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	_, err := coll.InsertOne(ctx, conversation)
	if err != nil {
		return fmt.Errorf("db.Client.InsertConversation: %w", err)
	}
	return nil
}

func (c *Client) InsertMessage(ctx context.Context, message *livechat.StdMessage) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	_, err := coll.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("mongodb.Client.InsertMessage: %w", err)
	}
	return nil
}

func (c *Client) UpdateConversationOnNewMessage(ctx context.Context, message *livechat.StdMessage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.UpdateConversationOnNewMessage: %w", err)
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
	err = coll.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoConversations
		}
		return err
	}
	return nil
}

func (c *Client) UpdateConversationIsRead(ctx context.Context, conversationID string) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{{Key: "conversationID", Value: conversationID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "isRead", Value: true},
	}}}
	err := coll.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("mongodb.Client.UpdateConversationIsRead %w", ErrNoConversations)
		}
		return err
	}
	return nil
}

func (c *Client) CheckConversationExists(ctx context.Context, conversationID string) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{{Key: "conversationID", Value: conversationID}}
	err := coll.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("mongodb.Client.CheckConversationExists %w", ErrNoConversations)
		}
		return err
	}
	return nil
}

func (c *Client) QueryMessages(ctx context.Context, pageID string, conversationID string) (_ []*livechat.StdMessage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryMessages: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	filter := bson.D{{Key: "pageID", Value: pageID}, {Key: "conversationID", Value: conversationID}}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoMessages
		}
		return nil, err
	}
	defer cur.Close(ctx)
	var messages []*livechat.StdMessage
	err = cur.All(ctx, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *Client) QueryConversations(ctx context.Context, pageID string) (_ []*livechat.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryConversations: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{{Key: "pageID", Value: pageID}}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoConversations
		}
		return nil, err
	}
	defer cur.Close(ctx)
	var conversations []*livechat.StdConversation
	err = cur.All(ctx, &conversations)
	if err != nil {
		return nil, err
	}
	return conversations, nil
}
