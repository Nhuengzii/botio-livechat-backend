package livechat

import (
	"context"
)

type DBClient interface {
	Close(ctx context.Context) error
	InsertConversation(ctx context.Context, conversation *StdConversation) error
	InsertMessage(ctx context.Context, message *StdMessage) error
	UpdateConversationOnNewMessage(ctx context.Context, message *StdMessage) error
	UpdateConversationIsRead(ctx context.Context, conversationID string) error
	CheckConversationExists(ctx context.Context, conversationID string) (bool, error)
}
