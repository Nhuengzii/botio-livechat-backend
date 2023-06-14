package getconversation

import "github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"

type Response struct {
	Conversation stdconversation.StdConversation `json:"conversation"`
}
