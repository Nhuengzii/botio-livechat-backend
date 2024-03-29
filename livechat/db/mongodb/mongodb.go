// Package mongodb implements DBClient for manipulating mongodb database
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getall"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getshop"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/shopcfg"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/templates"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/shops"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNoDocuments = mongo.ErrNoDocuments

// A Client contains mongodb client and a Target struct.
type Client struct {
	client *mongo.Client // mongodb's client used to do mongo operation
	Target               // target database's  information
}

// A Target contains information about the target database.
type Target struct {
	URI                     string // connection URI of the db
	Database                string // Database name
	CollectionConversations string // Conversations collection name
	CollectionMessages      string // Messages collection name
	CollectionShops         string // Shops collection name
	CollectionShopConfig    string // ShopConfig collection name
	CollectionTemplates     string // Templates collection name
}

// NewClient returns a new Client which contains mongodb client inside.
// Return an error if it occurs.
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

// Close close a mongodb client connection, releasing resources.
// Return an error if it occurs.
func (c *Client) Close(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("mongodb.Client.Close: %w", err)
	}
	return nil
}

// InsertConversation insert a document containing new conversation to the target's CollectionConversations.
// Return an error if it occurs.
func (c *Client) InsertConversation(ctx context.Context, conversation *stdconversation.StdConversation) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	_, err := coll.InsertOne(ctx, conversation)
	if err != nil {
		return fmt.Errorf("mongodb.Client.InsertConversation: %w", err)
	}
	return nil
}

// InsertMessage insert a document containing new message to the target's CollectionMessages.
// Return an error if it occurs.
func (c *Client) InsertMessage(ctx context.Context, message *stdmessage.StdMessage) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	_, err := coll.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("mongodb.Client.InsertMessage: %w", err)
	}
	return nil
}

// UpdateConversationOnDeletedMessage update a conversation based on recieved unsend message events.
// Return an error if it occurs.
//
// update last activity if the unsent message was the last message
func (c *Client) UpdateConversationOnDeletedMessage(ctx context.Context, message stdmessage.StdMessage) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.UpdateConversationOnDeletedMessage: %w", err)
		}
	}()
	if message.IsDeleted {
		coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
		filter := bson.D{
			{Key: "shopID", Value: message.ShopID},
			{Key: "conversationID", Value: message.ConversationID},
			{Key: "pageID", Value: message.PageID},
		}
		fOpt := options.FindOneOptions{
			Sort: bson.D{{Key: "timestamp", Value: -1}},
		}
		var lastMessage stdmessage.StdMessage
		err = coll.FindOne(ctx, filter, &fOpt).Decode(&lastMessage)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return ErrNoDocuments
			}
			return err
		}
		if message.MessageID == lastMessage.MessageID {
			// change last activity
			lastActivity, err := message.ToLastActivityString()
			if err != nil {
				return err
			}

			coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
			filter := bson.D{
				{Key: "shopID", Value: message.ShopID},
				{Key: "conversationID", Value: message.ConversationID},
			}
			update := bson.M{
				"$set": bson.D{
					{Key: "lastActivity", Value: lastActivity},
				},
			}
			err = coll.FindOneAndUpdate(ctx, filter, update).Err()
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return ErrNoDocuments
				}
				return err
			}
		}
	}
	return nil
}

// UpdateConversationOnNewMessage update a conversation based on recieved new message events.
// Return an error if it occurs.
//
// update last activity, updated time and increment unread count if the sender wasn't an UserTypeAdmin
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
		{Key: "shopID", Value: message.ShopID},
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

	if message.Source.UserType == stdmessage.UserTypeAdmin {
		update := bson.M{
			"$set": bson.D{
				{Key: "lastActivity", Value: lastActivity},
				{Key: "updatedTime", Value: message.Timestamp},
			},
		}
		err = coll.FindOneAndUpdate(ctx, filter, update).Err()
	} else {
		currentUnread := conversation.Unread
		update := bson.M{
			"$set": bson.D{
				{Key: "lastActivity", Value: lastActivity},
				{Key: "updatedTime", Value: message.Timestamp},
				{Key: "lastUserActivityTime", Value: message.Timestamp},
				{Key: "unread", Value: currentUnread + 1},
			},
		}
		err = coll.FindOneAndUpdate(ctx, filter, update).Err()
	}
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoDocuments
		}
		return err
	}

	return nil
}

// UpdateConversationUnread update a conversation's unread field to a specified integer.
// Return an error if it occurs.
func (c *Client) UpdateConversationUnread(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, conversationID string, unread int) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
		{Key: "conversationID", Value: conversationID},
	}
	update := bson.M{
		"$set": bson.D{
			{Key: "unread", Value: unread},
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

// CheckConversationExists return an ErrNoDocuments if the conversation with matching conversationID does not exist.
// If other errors occured CheckConversationExists will return that error.
//
// If the conversation was found CheckConversationExists return nil
func (c *Client) CheckConversationExists(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, conversationID string) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
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

// RemoveDeletedMessage update specific message's fields on unsend message events.
// Return an error if it occurs.
//
// update various fields
//   - isDeleted : true
//   - message : ""
//   - attachments : []
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

// ListMessages return a slice of stdmessage.StdMessage in a specific conversation.
// Only return messages in specific platform.
// Return an empty slice if none were found.
// Return an error if it occurs.
//
// The return values is sorted descending with the message's timestamp.
// This means that the queried slice will start with the latest message in the conversation.
//
// # pagination parameters (skip and limit should be input as nil if caller doesn't need any pagination)
//
//   - skip(integer): number of result messages to skip. Skip value should not be negative.
//   - limit(integer): number of maximum messages result. Limit value should not be negative.
func (c *Client) ListMessages(ctx context.Context, shopID string, platform stdmessage.Platform, pageID string, conversationID string, skip *int, limit *int) (_ []stdmessage.StdMessage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListMessages: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	filter := bson.M{
		"shopID":         shopID,
		"platform":       platform,
		"pageID":         pageID,
		"conversationID": conversationID,
	}

	var fOpt options.FindOptions
	fOpt.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // descending sort
	if limit != nil {
		fOpt.SetLimit(int64(*limit))
	}
	if skip != nil {
		fOpt.SetSkip(int64(*skip))
	}

	cur, err := coll.Find(ctx, filter, &fOpt)
	if err != nil {
		return nil, err
	}
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

// ListMessagesWithMessage return a slice of stdmessage.StdMessage in a specific conversation that has text message containing specified message string.
// Only return messages in specific platform.
// Return an empty slice if none were found.
// Return an error if it occurs.
//
// *** use case insensitive search ***
//
// The return values is sorted descending with the message's timestamp.
// This means that the queried slice will start with the latest message in the conversation.
//
// # pagination parameters (skip and limit should be input as nil if caller doesn't need any pagination)
//
//   - skip(integer): number of result messages to skip. Skip value should not be negative.
//   - limit(integer): number of maximum messages result. Limit value should not be negative.
func (c *Client) ListMessagesWithMessage(ctx context.Context, shopID string, platform stdmessage.Platform, pageID string, conversationID string, message string, skip *int, limit *int) (_ []stdmessage.StdMessage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListMessagesWithMessage: %w", err)
		}
	}()

	coll := c.client.Database(c.Database).Collection(c.CollectionMessages)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
		{Key: "conversationID", Value: conversationID},
		{Key: "message", Value: bson.M{
			"$regex": message, "$options": "i",
		}},
	}

	var fOpt options.FindOptions
	fOpt.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // descending sort
	if limit != nil {
		fOpt.SetLimit(int64(*limit))
	}
	if skip != nil {
		fOpt.SetSkip(int64(*skip))
	}

	cur, err := coll.Find(ctx, filter, &fOpt)
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

// GetConversation return a specific stdconversation.StdConversation that match the conversationID.
// Return an error if it occurs.
func (c *Client) GetConversation(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, conversationID string) (_ *stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.GetConversation: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
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

// ListConversations return a slice of stdconversation.StdConversation in a page.
// Only return conversations in specific platform.
// Return an empty slice if none were found.
// Return an error if it occurs.
//
// The return values is sorted descending with the conversation's lastActivity timestamp.
// This means that the queried slice will start with the latest conversation that an activity occured.
//
// # pagination parameters (skip and limit should be input as nil if caller doesn't need any pagination)
//
//   - skip(integer): number of result conversations to skip. Skip value should not be negative.
//   - limit(integer): number of maximum conversations result. Limit value should not be negative.
func (c *Client) ListConversations(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, skip *int, limit *int) (_ []stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListConversations: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
	}

	var fOpt options.FindOptions
	fOpt.SetSort(bson.D{{Key: "updatedTime", Value: -1}}) // descending sort
	if limit != nil {
		fOpt.SetLimit(int64(*limit))
	}
	if skip != nil {
		fOpt.SetSkip(int64(*skip))
	}

	cur, err := coll.Find(ctx, filter, &fOpt)
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

// ListConversationsWithParticipantsName return a slice of stdconversation.StdConversation in a specific page that has participants name containing input name string.
// Only return conversations in specific platform.
// Return an empty slice if none were found.
// Return an error if it occurs.
//
// *** use case insensitive search ***
//
// The return values is sorted descending with the conversation's lastActivity timestamp.
// This means that the queried slice will start with the latest conversation that an activity occured.
//
// # pagination parameters (skip and limit should be input as nil if caller doesn't need any pagination)
//
//   - skip(integer): number of result messages to skip. Skip value should not be negative.
//   - limit(integer): number of maximum messages result. Limit value should not be negative.
func (c *Client) ListConversationsWithParticipantsName(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, name string, skip *int, limit *int) (_ []stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListConversationsWithParticipantsName: %w", err)
		}
	}()

	name = strings.Trim(name, " ")
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)

	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
		{Key: "participants.username", Value: bson.M{"$regex": name, "$options": "i"}},
	}

	var fOpt options.FindOptions
	fOpt.SetSort(bson.D{{Key: "updatedTime", Value: -1}}) // descending sort
	if limit != nil {
		fOpt.SetLimit(int64(*limit))
	}
	if skip != nil {
		fOpt.SetSkip(int64(*skip))
	}

	cur, err := coll.Find(ctx, filter, &fOpt)
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

// ListConversationsWithMessage return a slice of stdconversation.StdConversation in a specific page that has text message containing input message string.
// Only return conversations in specific platform.
// Return an empty slice if none were found.
// Return an error if it occurs.
//
// *** use case insensitive search ***
//
// The return values is sorted descending with the conversation's lastActivity timestamp.
// This means that the queried slice will start with the latest conversation that an activity occured.
//
// # pagination parameters (skip and limit should be input as nil if caller doesn't need any pagination)
//
//   - skip(integer): number of result messages to skip. Skip value should not be negative.
//   - limit(integer): number of maximum messages result. Limit value should not be negative.
func (c *Client) ListConversationsWithMessage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, message string, skip *int, limit *int) (_ []stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListConversationsWithMessage: %w", err)
		}
	}()

	message = strings.Trim(message, " ")
	collMessage := c.client.Database(c.Database).Collection(c.CollectionMessages)

	filterMessage := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
		{Key: "message", Value: bson.M{"$regex": message, "$options": "i"}},
	}

	var fOpt options.FindOptions
	fOpt.SetSort(bson.D{{Key: "updatedTime", Value: -1}}) // descending sort
	if limit != nil {
		fOpt.SetLimit(int64(*limit))
	}
	if skip != nil {
		fOpt.SetSkip(int64(*skip))
	}

	cur, err := collMessage.Find(ctx, filterMessage, &fOpt)
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

// ListConversationsOfAllPlatforms return a slice of stdconversation.StdConversation in a page.
// Return conversations in all platform.
// Return an empty slice if none were found.
// Return an error if it occurs.
//
// The return values is sorted descending with the conversation's lastActivity timestamp.
// This means that the queried slice will start with the latest conversation that an activity occured.
//
// # pagination parameters (skip and limit should be input as nil if caller doesn't need any pagination)
//
//   - skip(integer): number of result conversations to skip. Skip value should not be negative.
//   - limit(integer): number of maximum conversations result. Limit value should not be negative.
func (c *Client) ListConversationsOfAllPlatforms(ctx context.Context, shopID string, skip *int, limit *int) ([]stdconversation.StdConversation, error) {
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
	}

	var fOpt options.FindOptions
	fOpt.SetSort(bson.D{{Key: "updatedTime", Value: -1}}) // descending sort
	if limit != nil {
		fOpt.SetLimit(int64(*limit))
	}
	if skip != nil {
		fOpt.SetSkip(int64(*skip))
	}

	cur, err := coll.Find(ctx, filter, &fOpt)
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

// ListConversationsOfAllPlatformsWithParticipantsName return a slice of stdconversation.StdConversation in a specific page that has participants name containing input name string.
// Return conversations in all platform.
// Return an empty slice if none were found.
// Return an error if it occurs.
//
// *** use case insensitive search ***
//
// The return values is sorted descending with the conversation's lastActivity timestamp.
// This means that the queried slice will start with the latest conversation that an activity occured.
//
// # pagination parameters (skip and limit should be input as nil if caller doesn't need any pagination)
//
//   - skip(integer): number of result messages to skip. Skip value should not be negative.
//   - limit(integer): number of maximum messages result. Limit value should not be negative.
func (c *Client) ListConversationsOfAllPlatformsWithParticipantsName(ctx context.Context, shopID string, name string, skip *int, limit *int) (_ []stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListConversationsWithParticipantsName: %w", err)
		}
	}()

	name = strings.Trim(name, " ")
	coll := c.client.Database(c.Database).Collection(c.CollectionConversations)

	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "participants.username", Value: bson.M{"$regex": name, "$options": "i"}},
	}

	var fOpt options.FindOptions
	fOpt.SetSort(bson.D{{Key: "updatedTime", Value: -1}}) // descending sort
	if limit != nil {
		fOpt.SetLimit(int64(*limit))
	}
	if skip != nil {
		fOpt.SetSkip(int64(*skip))
	}

	cur, err := coll.Find(ctx, filter, &fOpt)
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

// ListConversationsOfAllPlatformsWithMessage return a slice of stdconversation.StdConversation in a specific page that has text message containing input message string.
// Return conversations in all platform.
// Return an empty slice if none were found.
// Return an error if it occurs.
//
// *** use case insensitive search ***
//
// The return values is sorted descending with the conversation's lastActivity timestamp.
// This means that the queried slice will start with the latest conversation that an activity occured.
//
// # pagination parameters (skip and limit should be input as nil if caller doesn't need any pagination)
//
//   - skip(integer): number of result messages to skip. Skip value should not be negative.
//   - limit(integer): number of maximum messages result. Limit value should not be negative.
func (c *Client) ListConversationsOfAllPlatformsWithMessage(ctx context.Context, shopID string, message string, skip *int, limit *int) (_ []stdconversation.StdConversation, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListConversationsWithMessage: %w", err)
		}
	}()

	message = strings.Trim(message, " ")
	collMessage := c.client.Database(c.Database).Collection(c.CollectionMessages)

	filterMessage := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "message", Value: bson.M{"$regex": message, "$options": "i"}},
	}

	var fOpt options.FindOptions
	fOpt.SetSort(bson.D{{Key: "updatedTime", Value: -1}}) // descending sort
	if limit != nil {
		fOpt.SetLimit(int64(*limit))
	}
	if skip != nil {
		fOpt.SetSkip(int64(*skip))
	}

	cur, err := collMessage.Find(ctx, filterMessage, &fOpt)
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

// GetShop return shops.Shop that contains a matching pageID of any platform.
// Return an error if it occurs.
func (c *Client) GetShop(ctx context.Context, pageID string) (_ *shops.Shop, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.GetShop: %w", err)
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

// GetFacebookAuthentication return shops.FacebookAuthentication that contains a matching pageID of facebook platform.
// Return an error if it occurs.
//
// Can be use to get access token
func (c *Client) GetFacebookAuthentication(ctx context.Context, pageID string) (_ *shops.FacebookAuthentication, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.GetFacebookAuthentication: %w", err)
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

// GetLineAuthentication return shops.GetLineAuthentication that contains a matching pageID of line platform.
// Return an error if it occurs.
//
// Can be use to get access token
func (c *Client) GetLineAuthentication(ctx context.Context, pageID string) (_ *shops.LineAuthentication, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.GetLineAuthentication: %w", err)
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

// GetInstagramAuthentication return shops.GetInstagramAuthentication that contains a matching pageID of instagram platform.
// Return an error if it occurs.
//
// Can be use to get access token
func (c *Client) GetInstagramAuthentication(ctx context.Context, pageID string) (_ *shops.InstagramAuthentication, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.GetInstagramAuthentication: %w", err)
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

// GetPage return number of unread conversations and total conversations of the specified page.
// Return an error if it occurs.
func (c *Client) GetPage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string) (_ int64, _ int64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.GetPage: %w", err)
		}
	}()

	collConversations := c.client.Database(c.Database).Collection(c.CollectionConversations)

	filterUnreadConversations := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
		{Key: "unread", Value: bson.D{
			{Key: "$gt", Value: 0},
		}},
	}
	unreadConversations, err := collConversations.CountDocuments(ctx, filterUnreadConversations)
	if err != nil {
		return 0, 0, err
	}

	filterPageMessages := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "platform", Value: platform},
		{Key: "pageID", Value: pageID},
	}
	allConversations, err := collConversations.CountDocuments(ctx, filterPageMessages)
	if err != nil {
		return 0, 0, err
	}

	return unreadConversations, allConversations, nil
}

// InsertShop creates a document in the mongodb "shops" collection with the information provided .
// It returns nil if the operation is successful, otherwise returns error.
func (c *Client) InsertShop(ctx context.Context, shop shops.Shop) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	_, err := coll.InsertOne(ctx, shop)
	if err != nil {
		return fmt.Errorf("mongodb.Client.InsertShop: %w", err)
	}
	return nil
}

// UpdateShop updates a shop with the given shopID with the information provided in the given shop shops.Shop.
func (c *Client) UpdateShop(ctx context.Context, shopID string, shop shops.Shop) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.UpdateShop: %w", err)
		}
	}()

	setElements := bson.D{}
	if shop.FacebookPageID != "" {
		setElements = append(setElements, bson.E{Key: "facebookPageID", Value: shop.FacebookPageID})
	}
	if shop.FacebookAuthentication.AccessToken != "" {
		setElements = append(setElements, bson.E{Key: "facebookAuthentication.accessToken", Value: shop.FacebookAuthentication.AccessToken})
	}
	if shop.InstagramPageID != "" {
		setElements = append(setElements, bson.E{Key: "instagramPageID", Value: shop.InstagramPageID})
	}
	if shop.InstagramAuthentication.AccessToken != "" {
		setElements = append(setElements, bson.E{Key: "instagramAuthentication.accessToken", Value: shop.InstagramAuthentication.AccessToken})
	}
	if shop.LinePageID != "" {
		setElements = append(setElements, bson.E{Key: "linePageID", Value: shop.LinePageID})
	}
	if shop.LineAuthentication.AccessToken != "" {
		setElements = append(setElements, bson.E{Key: "lineAuthentication.accessToken", Value: shop.LineAuthentication.AccessToken})
	}
	if shop.LineAuthentication.Secret != "" {
		setElements = append(setElements, bson.E{Key: "lineAuthentication.secret", Value: shop.LineAuthentication.Secret})
	}

	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
	}
	update := bson.D{
		{Key: "$set", Value: setElements},
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

// ListShopPlatforms returns a slice of a shop's platforms and corresponding pageIDs.
// If the operation is successful, a slice will be returned and err will be nil,
// otherwise nil, nil are returned
func (c *Client) ListShopPlatforms(ctx context.Context, shopID string) (_ []getshop.PlatformPageID, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListShopPlatforms: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
	}
	shop := shops.Shop{}
	err = coll.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	result := []getshop.PlatformPageID{}
	if shop.FacebookPageID != "" {
		result = append(result, getshop.PlatformPageID{
			PlatformName: shops.PlatformFacebook,
			PageID:       shop.FacebookPageID,
		})
	}
	if shop.InstagramPageID != "" {
		result = append(result, getshop.PlatformPageID{
			PlatformName: shops.PlatformInstagram,
			PageID:       shop.InstagramPageID,
		})
	}
	if shop.LinePageID != "" {
		result = append(result, getshop.PlatformPageID{
			PlatformName: shops.PlatformLine,
			PageID:       shop.LinePageID,
		})
	}
	return result, nil
}

// ListShopPlatformsStatuses returns a slice of a shop's platforms statuses (unread and all conversations counts)
// Returns a slice and error = nil if successful.
// Otherwise,  returns a nil slice and an error.
func (c *Client) ListShopPlatformsStatuses(ctx context.Context, shopID string) (_ []getall.Status, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.ListShopPlatformsStatuses: %w", err)
		}
	}()
	collShops := c.client.Database(c.Database).Collection(c.CollectionShops)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
	}
	shop := shops.Shop{}
	err = collShops.FindOne(ctx, filter).Decode(&shop)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	result := []getall.Status{}

	if shop.FacebookPageID != "" {
		fbUnread, fbAll, err := c.GetPage(ctx, shopID, stdconversation.PlatformFacebook, shop.FacebookPageID)
		if err != nil {
			return nil, err
		}
		result = append(result, getall.Status{
			Platform:            shops.PlatformFacebook,
			UnreadConversations: fbUnread,
			AllConversations:    fbAll,
		})
	}
	if shop.InstagramPageID != "" {
		igUnread, igAll, err := c.GetPage(ctx, shopID, stdconversation.PlatformInstagram, shop.InstagramPageID)
		if err != nil {
			return nil, err
		}
		result = append(result, getall.Status{
			Platform:            shops.PlatformInstagram,
			UnreadConversations: igUnread,
			AllConversations:    igAll,
		})
	}
	if shop.LinePageID != "" {
		lineUnread, lineAll, err := c.GetPage(ctx, shopID, stdconversation.PlatformLine, shop.LinePageID)
		if err != nil {
			return nil, err
		}
		result = append(result, getall.Status{
			Platform:            shops.PlatformLine,
			UnreadConversations: lineUnread,
			AllConversations:    lineAll,
		})
	}

	return result, nil
}

// InsertShopConfig inserts a shop's config into the database.
func (c *Client) InsertShopConfig(ctx context.Context, config shopcfg.Config) error {
	coll := c.client.Database(c.Database).Collection(c.CollectionShopConfig)
	_, err := coll.InsertOne(ctx, config)
	if err != nil {
		return fmt.Errorf("mongodb.Client.InsertShopConfig: %w", err)
	}
	return nil
}

// GetShopConfig returns a shop's config.
func (c *Client) GetShopConfig(ctx context.Context, shopID string) (_ *shopcfg.Config, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.GetShopConfig: %w", err)
		}
	}()
	coll := c.client.Database(c.Database).Collection(c.CollectionShopConfig)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
	}
	config := shopcfg.Config{}
	err = coll.FindOne(ctx, filter).Decode(&config)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	return &config, nil
}

// GetShopTemplateMessages returns array of template messages of specific shop.
func (c *Client) GetShopTemplateMessages(ctx context.Context, shopID string) (_ []templates.Template, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.GetShopTemplateMessage: %w", err)
		}
	}()

	coll := c.client.Database(c.Database).Collection(c.CollectionTemplates)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
	}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	template := []templates.Template{}
	err = cur.All(ctx, &template)
	if err != nil {
		return nil, err
	}
	if cur.Err() != nil {
		return nil, cur.Err()
	}
	return template, nil
}

// InsertShopTemplateMessage adds a new template message to a shop's config.
func (c *Client) InsertShopTemplateMessage(ctx context.Context, template templates.Template) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.InsertShopTemplateMessage: %w", err)
		}
	}()

	coll := c.client.Database(c.Database).Collection(c.CollectionTemplates)
	_, err = coll.InsertOne(ctx, template)
	if err != nil {
		return fmt.Errorf("mongodb.Client.InsertShopTemplateMessage: %w", err)
	}
	return nil
}

// DeleteShopTemplateMessage removes a template from a shop_config's templates and return number of deleted document.
//
// return error if it occurs.
func (c *Client) DeleteShopTemplateMessage(ctx context.Context, shopID string, templateID string) (deleteCount int, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mongodb.Client.DeleteShopTemplateMessage: %w", err)
		}
	}()

	coll := c.client.Database(c.Database).Collection(c.CollectionTemplates)
	filter := bson.D{
		{Key: "shopID", Value: shopID},
		{Key: "templateID", Value: templateID},
	}

	delResult, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(delResult.DeletedCount), err
}
