package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/shops"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNoDocuments = mongo.ErrNoDocuments

type Client struct {
	client *mongo.Client
	Target
}

type Target struct {
	URI                     string
	Database                string
	CollectionConversations string
	CollectionMessages      string
	CollectionShops         string
}

func NewClient(ctx context.Context, target Target) (*Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(target.URI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("mongodb.NewClient: %w", err)
	}
	return &Client{
		client: client,
		Target: target,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("mongodb.Client.Close: %w", err)
	}
	return nil
}

func (c *Client) InsertConversation(ctx context.Context, conversation *stdconversation.StdConversation) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	_, err := coll.InsertOne(ctx, conversation)
	if err != nil {
		return fmt.Errorf("mongodb.Client.InsertConversation: %w", err)
	}
	return nil
}

func (c *Client) InsertMessage(ctx context.Context, message *stdmessage.StdMessage) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	_, err := coll.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("mongodb.Client.InsertMessage: %w", err)
	}
	return nil
}

func (c *Client) UpdateConversationOnNewMessage(ctx context.Context, message *stdmessage.StdMessage) (err error) {
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
	filter := bson.D{
		{Key: "conversationID", Value: message.ConversationID},
	}
	update := bson.M{
		"$set": bson.D{
			{Key: "lastActivity", Value: lastActivity},
			{Key: "updatedTime", Value: message.Timestamp},
			{Key: "isRead", Value: false},
		},
	}
	err = coll.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoDocuments
		}
		return err
	}
	return nil
}

func (c *Client) UpdateConversationIsRead(ctx context.Context, conversationID string) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "conversationID", Value: conversationID},
	}
	update := bson.M{
		"$set": bson.D{
			{Key: "isRead", Value: true},
		},
	}
	err := coll.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("mongodb.Client.UpdateConversationIsRead %w", ErrNoDocuments)
		}
		return err
	}
	return nil
}

func (c *Client) CheckConversationExists(ctx context.Context, conversationID string) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "conversationID", Value: conversationID},
	}
	err := coll.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("mongodb.Client.CheckConversationExists %w", ErrNoDocuments)
		}
		return err
	}
	return nil
}

func (c *Client) UpdateConversationParticipants(ctx context.Context, conversationID string) error {
	// TODO implement
	return nil
}

func (c *Client) QueryMessages(ctx context.Context, shopID string, pageID string, conversationID string) (_ []*stdmessage.StdMessage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryMessages: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "pageID", Value: pageID},
		{Key: "conversationID", Value: conversationID},
	}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	defer cur.Close(ctx)
	var messages []*stdmessage.StdMessage
	err = cur.All(ctx, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *Client) QueryConversations(ctx context.Context, shopID string, pageID string) (_ []*stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryConversations: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "pageID", Value: pageID},
	}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	defer cur.Close(ctx)
	var conversations []*stdconversation.StdConversation
	err = cur.All(ctx, &conversations)
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (c *Client) QueryShop(ctx context.Context, pageID string) (_ *shops.Shop, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryShop: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.M{
		"$or": []bson.D{
			{
				{Key: "facebookPages", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "pageID", Value: pageID},
					}},
				}},
			},
			{
				{Key: "linePages", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "pageID", Value: pageID},
					}},
				}},
			},
			{
				{Key: "instagramPages", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "pageID", Value: pageID},
					}},
				}},
			},
		},
	}
	var shop shops.Shop
	err = coll.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	return &shop, nil
}

func (c *Client) QueryFacebookPage(ctx context.Context, pageID string) (_ *shops.FacebookPage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.QueryFacebookPage: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "facebookPages", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "pageID", Value: pageID},
			}},
		}},
	}
	var shop shops.Shop
	err = coll.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	var facebookPage *shops.FacebookPage
	for _, page := range shop.FacebookPages {
		if pageID == page.PageID {
			facebookPage = page
			break
		}
	}
	return facebookPage, nil
}

func (c *Client) QueryLinePage(ctx context.Context, pageID string) (_ *shops.LinePage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.QueryLinePage: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "linePages", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "pageID", Value: pageID},
			}},
		}},
	}
	var shop shops.Shop
	err = coll.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	var linePage *shops.LinePage
	for _, page := range shop.LinePages {
		if pageID == page.PageID {
			linePage = page
			break
		}
	}
	return linePage, nil
}

func (c *Client) QueryInstagramPage(ctx context.Context, pageID string) (_ *shops.InstagramPage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.QueryInstagramPage: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "instagramPages", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "pageID", Value: pageID},
			}},
		}},
	}
	var shop shops.Shop
	err = coll.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	var instagramPage *shops.InstagramPage
	for _, page := range shop.InstagramPages {
		if pageID == page.PageID {
			instagramPage = page
			break
		}
	}
	return instagramPage, nil
}
