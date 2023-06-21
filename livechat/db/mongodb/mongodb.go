package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	var conversation stdconversation.StdConversation
	err = coll.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoDocuments
		}
		return err
	}
	currentUnread := conversation.Unread
	update := bson.M{
		"$set": bson.D{
			{Key: "lastActivity", Value: lastActivity},
			{Key: "updatedTime", Value: message.Timestamp},
			{Key: "unread", Value: currentUnread + 1},
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
			{Key: "unread", Value: 0},
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

func (c *Client) RemoveDeletedMessage(ctx context.Context, shopID string, platform stdmessage.Platform, conversationID string, messageID string) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "conversationID", Value: conversationID},
		{Key: "messageID", Value: messageID},
	}
	update := bson.M{
		"$set": bson.D{
			{Key: "isDeleted", Value: true},
			{Key: "message", Value: ""},
			{Key: "attachments", Value: bson.A{}},
		},
	}
	err := coll.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoDocuments
		}
		return err
	}
	return nil
}

func (c *Client) QueryMessages(ctx context.Context, shopID string, pageID string, conversationID string) (_ []stdmessage.StdMessage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryMessages: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	filter := bson.M{
		"shopID":         shopID,
		"pageID":         pageID,
		"conversationID": conversationID,
	}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	messages := []stdmessage.StdMessage{}
	err = cur.All(ctx, &messages)
	if err != nil {
		return nil, err
	}
	if cur.Err() != nil {
		return nil, cur.Err()
	}
	return messages, nil
}

func (c *Client) QueryMessagesWithMessage(ctx context.Context, shopID string, platform stdmessage.Platform, pageID string, conversationID string, message string) (_ []stdmessage.StdMessage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryMessagesWithMessage: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
		{Key: "conversationID", Value: conversationID},
		{Key: "message", Value: bson.D{
			{Key: "$regex", Value: message},
		}},
	}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	messages := []stdmessage.StdMessage{}
	err = cur.All(ctx, &messages)
	if err != nil {
		return nil, err
	}
	if cur.Err() != nil {
		return nil, cur.Err()
	}
	return messages, nil
}

func (c *Client) QueryConversation(ctx context.Context, shopID string, pageID string, conversationID string) (_ *stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryConversation: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "pageID", Value: pageID},
		{Key: "conversationID", Value: conversationID},
	}
	var conversation stdconversation.StdConversation
	err = coll.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	return &conversation, nil
}

func (c *Client) QueryConversations(ctx context.Context, shopID string, pageID string) (_ []stdconversation.StdConversation, err error) {
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
		return nil, err
	}
	defer cur.Close(ctx)
	conversations := []stdconversation.StdConversation{}
	err = cur.All(ctx, &conversations)
	if err != nil {
		return nil, err
	}
	if cur.Err() != nil {
		return nil, cur.Err()
	}
	return conversations, nil
}

func (c *Client) QueryConversationsWithParticipantsName(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, name string) (_ []stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryConversationsWithParticipantsName: %w", err)
		}
	}()

	name = strings.Trim(name, " ")
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
		{Key: "participants.username", Value: bson.D{{Key: "$regex", Value: name}}},
	}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	conversations := []stdconversation.StdConversation{}
	err = cur.All(ctx, &conversations)
	if err != nil {
		return nil, err
	} else if cur.Err() != nil {
		return nil, cur.Err()
	}
	return conversations, nil
}

func (c *Client) QueryConversationsWithMessage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, message string) (_ []stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.QueryConversationsWithMessage: %w", err)
		}
	}()

	message = strings.Trim(message, " ")
	collMessage := c.client.Database(c.Database).Collection(c.CollectionMessages)
	filterMessage := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
		{Key: "message", Value: bson.D{{Key: "$regex", Value: message}}},
	}
	cur, err := collMessage.Find(ctx, filterMessage)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	conversations := []stdconversation.StdConversation{}

	uniqueConversationIDSet := map[string]struct{}{} // using map to implement set

	for cur.Next(ctx) {
		var message stdmessage.StdMessage
		err := cur.Decode(&message)
		if err != nil {
			return nil, err
		}
		uniqueConversationIDSet[message.ConversationID] = struct{}{} // add conversationID to set
	}

	var uniqueConversationIDFilter []string

	for conversationID := range uniqueConversationIDSet {
		uniqueConversationIDFilter = append(uniqueConversationIDFilter, conversationID)
	}

	if len(uniqueConversationIDFilter) != 0 {
		collConversation := c.client.Database(c.Database).Collection(c.CollectionConversations)
		filterConversation := bson.M{"conversationID": bson.M{"$in": uniqueConversationIDFilter}}
		cur, err = collConversation.Find(ctx, filterConversation)
		if err != nil {
			return nil, err
		}
		err = cur.All(ctx, &conversations)
		if err := cur.Err(); err != nil {
			return nil, err
		}
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
				{Key: "facebookPageID", Value: pageID},
			},
			{
				{Key: "linePageID", Value: pageID},
			},
			{
				{Key: "instagramPageID", Value: pageID},
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

func (c *Client) QueryFacebookAuthentication(ctx context.Context, pageID string) (_ *shops.FacebookAuthentication, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.QueryFacebookAuthentication: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "facebookPageID", Value: pageID},
	}
	var shop shops.Shop
	err = coll.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	return &shop.FacebookAuthentication, nil
}

func (c *Client) QueryLineAuthentication(ctx context.Context, pageID string) (_ *shops.LineAuthentication, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.QueryLineAuthentication: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "linePageID", Value: pageID},
	}
	var shop shops.Shop
	err = coll.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	return &shop.LineAuthentication, nil
}

func (c *Client) QueryInstagramAuthentication(ctx context.Context, pageID string) (_ *shops.InstagramAuthentication, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.QueryInstagramAuthentication: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "instagramPageID", Value: pageID},
	}
	var shop shops.Shop
	err = coll.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	return &shop.InstagramAuthentication, nil
}
