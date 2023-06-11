package getconversationsresp

import "github.com/Nhuengzii/botio-livechat-backend/livechat"

type Response struct {
	Conversations []*livechat.StdConversation `json:"conversations"`
}
