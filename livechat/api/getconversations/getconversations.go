package getconversations

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
)

type Response struct {
	Conversations []stdconversation.StdConversation `json:"conversations"`
}

type Filter struct { // only 1 fields can exist at the same time
	ParticipantsUsername string `json:"with_participants_username"`
	Message              string `json:"with_message"`
}
