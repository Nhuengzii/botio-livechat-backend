package livechat

import (
	"context"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getshop"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/shops"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

// DBClient is an interface for database client, for example, MongoDB
type DBClient interface {
	// Close close a db client connection, releasing resources.
	// Return an error if it occurs.
	Close(ctx context.Context) error

	// InsertConversation insert a document containing new conversation to the conversations collection.
	// Return an error if it occurs.
	InsertConversation(ctx context.Context, conversation *stdconversation.StdConversation) error
	// InsertMessage insert a document containing new message to the messages collection.
	// Return an error if it occurs.
	InsertMessage(ctx context.Context, message *stdmessage.StdMessage) error
	// UpdateConversationOnDeletedMessage update a conversation based on recieved unsend message events.
	// Return an error if it occurs.
	//
	// update last activity if the unsent message was the last message
	UpdateConversationOnDeletedMessage(ctx context.Context, message *stdmessage.StdMessage) error
	// UpdateConversationOnNewMessage update a conversation based on recieved new message events.
	// Return an error if it occurs.
	//
	// update last activity, updated time and increment unread count if the sender wasn't an admin user type.
	UpdateConversationOnNewMessage(ctx context.Context, message *stdmessage.StdMessage) error
	// UpdateConversationUnread update a conversation's unread field to a specified integer.
	// Return an error if it occurs.
	UpdateConversationUnread(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, conversationID string, unread int) error
	// CheckConversationExists return an no document error if the conversation with matching conversationID does not exist.
	// If other errors occured CheckConversationExists will return that error.
	//
	// If the conversation was found CheckConversationExists return nil
	CheckConversationExists(ctx context.Context, conversationID string) error
	// RemoveDeletedMessage update specific message's fields on unsend message events.
	// Return an error if it occurs.
	//
	// update various fields
	//   - isDeleted : true
	//   - message : ""
	//   - attachments : []
	RemoveDeletedMessage(ctx context.Context, shopID string, platform stdmessage.Platform, conversationID string, messageID string) error
	// QueryConversation return a specific stdconversation.StdConversation that match the conversationID.
	// Return an error if it occurs.
	QueryConversation(ctx context.Context, shopID string, pageID string, conversationID string) (*stdconversation.StdConversation, error)
	// QueryMessages return a slice of stdmessage.StdMessage in a specific conversation.
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
	QueryMessages(ctx context.Context, shopID string, pageID string, conversationID string, skip *int, limit *int) ([]stdmessage.StdMessage, error)
	// QueryMessages return a slice of stdmessage.StdMessage in a specific conversation that has text message containing specified message string.
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
	QueryMessagesWithMessage(ctx context.Context, shopID string, platform stdmessage.Platform, pageID string, conversationID string, message string, skip *int, limit *int) ([]stdmessage.StdMessage, error)
	// QueryConversations return a slice of stdconversation.StdConversation in a page.
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
	QueryConversations(ctx context.Context, shopID string, pageID string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
	// QueryConversationsWithParticipantsName return a slice of stdconversation.StdConversation in a specific page that has participants name containing input name string.
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
	QueryConversationsWithParticipantsName(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, name string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
	// QueryConversationsWithMessage return a slice of stdconversation.StdConversation in a specific page that has text message containing input message string.
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
	QueryConversationsWithMessage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, message string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
	// QueryShop return shops.Shop that contains a matching pageID of any platform.
	// Return an error if it occurs.
	QueryShop(ctx context.Context, pageID string) (*shops.Shop, error)
	// QueryConversationsOfAllPlatforms return a slice of stdconversation.StdConversation in a page.
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
	QueryConversationsOfAllPlatforms(ctx context.Context, shopID string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
	// QueryConversationsOfAllPlatformsWithParticipantsName return a slice of stdconversation.StdConversation in a specific page that has participants name containing input name string.
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
	QueryConversationsOfAllPlatformsWithParticipantsName(ctx context.Context, shopID string, name string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
	// QueryConversationsWithMessage return a slice of stdconversation.StdConversation in a specific page that has text message containing input message string.
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
	QueryConversationsOfAllPlatformsWithMessage(ctx context.Context, shopID string, message string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
	// QueryFacebookAuthentication return shops.FacebookAuthentication that contains a matching pageID of facebook platform.
	// Return an error if it occurs.
	//
	// Can be use to get access token
	QueryFacebookAuthentication(ctx context.Context, pageID string) (*shops.FacebookAuthentication, error)
	// QueryLineAuthentication return shops.QueryLineAuthentication that contains a matching pageID of line platform.
	// Return an error if it occurs.
	//
	// Can be use to get access token
	QueryLineAuthentication(ctx context.Context, pageID string) (*shops.LineAuthentication, error)
	// QueryInstagramAuthentication return shops.QueryInstagramAuthentication that contains a matching pageID of instagram platform.
	// Return an error if it occurs.
	//
	// Can be use to get access token
	QueryInstagramAuthentication(ctx context.Context, pageID string) (*shops.InstagramAuthentication, error)
	// GetPage return number of unread conversations and total conversations of the specified page.
	// Return an error if it occurs.
	GetPage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string) (unreadConversations int64, allConversations int64, err error)

	// InsertShop creates a document in the mongodb "shops" collection with the information provided .
	// It returns nil if the operation is successful, otherwise returns error.
	InsertShop(ctx context.Context, shop shops.Shop) error

	// CheckShopExists returns nil if a shop with shopID already exists, if not returns error wrapping mongodb.ErrorNoDocuments,
	// otherwise returns error.
	CheckShopExists(ctx context.Context, shopID string) error

	// ListShopPlatforms returns a slice of all shop's platforms and corresponding pageID.
	// If the operation is successful, a slice will be returned and err will be nil,
	// otherwise nil, nil are returned
	ListShopPlatforms(ctx context.Context, shopID string) (_ []getshop.Platform, err error)
}
