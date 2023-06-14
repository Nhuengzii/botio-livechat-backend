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
	UpdateConversationIsRead(ctx context.Context, conversationID string) error
	CheckConversationExists(ctx context.Context, conversationID string) error
	UpdateConversationParticipants(ctx context.Context, conversationID string) error
	QueryConversation(ctx context.Context, shopID string, pageID string, conversationID string) (*stdconversation.StdConversation, error)
	QueryMessages(ctx context.Context, shopID string, pageID string, conversationID string) ([]stdmessage.StdMessage, error)
	QueryConversations(ctx context.Context, shopID string, pageID string) ([]stdconversation.StdConversation, error)
	QueryShop(ctx context.Context, pageID string) (*shops.Shop, error)
	QueryFacebookPage(ctx context.Context, pageID string) (*shops.FacebookPage, error)
	QueryLinePage(ctx context.Context, pageID string) (*shops.LinePage, error)
	QueryInstagramPage(ctx context.Context, pageID string) (*shops.InstagramPage, error)
}
