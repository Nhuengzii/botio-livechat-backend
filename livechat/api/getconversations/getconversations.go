package getconversations

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
)

type Response struct {
	Conversations []*stdconversation.StdConversation `json:"conversations"`
}
