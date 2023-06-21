package getpage

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

type Response struct {
	UnreadConversations []stdconversation.StdConversation `json:"unreadConversations"`
	AllMessages         []stdmessage.StdMessage           `json:"allMessages"`
}
