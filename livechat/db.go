package livechat

import (
	"context"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/shops"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

type DBClient interface {
	Close(ctx context.Context) error
	InsertConversation(ctx context.Context, conversation *stdconversation.StdConversation) error
	InsertMessage(ctx context.Context, message *stdmessage.StdMessage) error
	UpdateConversationOnNewMessage(ctx context.Context, message *stdmessage.StdMessage) error
	UpdateConversationUnread(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, conversationID string, unread int) error
	CheckConversationExists(ctx context.Context, conversationID string) error
	UpdateConversationParticipants(ctx context.Context, conversationID string) error
	RemoveDeletedMessage(ctx context.Context, shopID string, platform stdmessage.Platform, conversationID string, messageID string) error
	QueryConversation(ctx context.Context, shopID string, pageID string, conversationID string) (*stdconversation.StdConversation, error)
	QueryMessages(ctx context.Context, shopID string, pageID string, conversationID string, offset *int, limit *int) ([]stdmessage.StdMessage, error)
	QueryMessagesWithMessage(ctx context.Context, shopID string, platform stdmessage.Platform, pageID string, conversationID string, message string) ([]stdmessage.StdMessage, error)
	QueryConversations(ctx context.Context, shopID string, pageID string, offset *int, limit *int) ([]stdconversation.StdConversation, error)
	QueryConversationsWithParticipantsName(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, name string, offset *int, limit *int) ([]stdconversation.StdConversation, error)
	QueryConversationsWithMessage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string, message string, offset *int, limit *int) ([]stdconversation.StdConversation, error)
	QueryShop(ctx context.Context, pageID string) (*shops.Shop, error)
	QueryFacebookAuthentication(ctx context.Context, pageID string) (*shops.FacebookAuthentication, error)
	QueryLineAuthentication(ctx context.Context, pageID string) (*shops.LineAuthentication, error)
	QueryInstagramAuthentication(ctx context.Context, pageID string) (*shops.InstagramAuthentication, error)
	GetPage(ctx context.Context, shopID string, platform stdconversation.Platform, pageID string) (unreadConversations int64, allConversations int64, err error)
}
