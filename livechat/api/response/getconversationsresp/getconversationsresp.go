package getconversationsresp

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
)

type Resp struct {
	Conversations []*stdconversation.StdConversation `json:"conversations"`
}
