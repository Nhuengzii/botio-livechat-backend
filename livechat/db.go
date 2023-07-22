package livechat

import (
	"context"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getall"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getshop"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/shopcfg"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/templates"

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
	UpdateConversationOnDeletedMessage(ctx context.Context, message stdmessage.StdMessage) error
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
	CheckConversationExists(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, conversationID string) error
	// RemoveDeletedMessage update specific message's fields on unsend message events.
	// Return an error if it occurs.
	//
	// update various fields
	//   - isDeleted : true
	//   - message : ""
	//   - attachments : []
	RemoveDeletedMessage(ctx context.Context, shopID string, platform stdmessage.Platform, conversationID string, messageID string) error
	// GetConversation return a specific stdconversation.StdConversation that match the conversationID.
	// Return an error if it occurs.
	GetConversation(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, conversationID string) (*stdconversation.StdConversation, error)
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
	ListMessages(ctx context.Context, shopID string, pageID string, conversationID string, skip *int, limit *int) ([]stdmessage.StdMessage, error)
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
	ListMessagesWithMessage(ctx context.Context, shopID string, platform stdmessage.Platform, pageID string, conversationID string, message string, skip *int, limit *int) ([]stdmessage.StdMessage, error)
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
	ListConversations(ctx context.Context, shopID string, pageID string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
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
	ListConversationsWithParticipantsName(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, name string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
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
	ListConversationsWithMessage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, message string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
	// GetShop return shops.Shop that contains a matching pageID of any platform.
	// Return an error if it occurs.
	GetShop(ctx context.Context, pageID string) (*shops.Shop, error)
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
	ListConversationsOfAllPlatforms(ctx context.Context, shopID string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
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
	ListConversationsOfAllPlatformsWithParticipantsName(ctx context.Context, shopID string, name string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
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
	ListConversationsOfAllPlatformsWithMessage(ctx context.Context, shopID string, message string, skip *int, limit *int) ([]stdconversation.StdConversation, error)
	// GetFacebookAuthentication return shops.FacebookAuthentication that contains a matching pageID of facebook platform.
	// Return an error if it occurs.
	//
	// Can be use to get access token
	GetFacebookAuthentication(ctx context.Context, pageID string) (*shops.FacebookAuthentication, error)
	// GetLineAuthentication return shops.GetLineAuthentication that contains a matching pageID of line platform.
	// Return an error if it occurs.
	//
	// Can be use to get access token
	GetLineAuthentication(ctx context.Context, pageID string) (*shops.LineAuthentication, error)
	// GetInstagramAuthentication return shops.GetInstagramAuthentication that contains a matching pageID of instagram platform.
	// Return an error if it occurs.
	//
	// Can be use to get access token
	GetInstagramAuthentication(ctx context.Context, pageID string) (*shops.InstagramAuthentication, error)
	// GetPage return number of unread conversations and total conversations of the specified page.
	// Return an error if it occurs.
	GetPage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string) (unreadConversations int64, allConversations int64, err error)

	// InsertShop creates a document in the mongodb "shops" collection with the information provided .
	// It returns nil if the operation is successful, otherwise returns error.
	InsertShop(ctx context.Context, shop shops.Shop) error

	// UpdateShop updates a shop with the given shopID with the information provided in the given shop shops.Shop.
	UpdateShop(ctx context.Context, shopID string, shop shops.Shop) (err error)

	// ListShopPlatforms returns a slice of a shop's platforms and corresponding pageIDs.
	// If the operation is successful, a slice will be returned and err will be nil,
	// otherwise nil, nil are returned
	ListShopPlatforms(ctx context.Context, shopID string) (_ []getshop.PlatformPageID, err error)

	// ListShopPlatformsStatuses returns a slice of a shop's platforms statuses (unread and all conversations counts)
	// Returns a slice and error = nil if successful.
	// Otherwise,  returns a nil slice and an error.
	ListShopPlatformsStatuses(ctx context.Context, shopID string) (_ []getall.Status, err error)

	// InsertShopConfig inserts a shop's config into the database.
	InsertShopConfig(ctx context.Context, config shopcfg.Config) error

	// GetShopConfig returns a shop's config.
	GetShopConfig(ctx context.Context, shopID string) (_ *shopcfg.Config, err error)

	// InsertShopTemplateMessage adds a new template message to a shop's config.
	InsertShopTemplateMessage(ctx context.Context, template templates.Template) (err error)

	// GetShopTemplateMessages returns array of template messages of specific shop.
	GetShopTemplateMessages(ctx context.Context, shopID string) (_ []templates.Template, err error)

	// DeleteShopTemplateMessage removes a template from a shop_config's templates
	DeleteShopTemplateMessage(ctx context.Context, shopID string, templateID string) (deletedCount int, err error)
}
