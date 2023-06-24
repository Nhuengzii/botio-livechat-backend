// Package getconversations define response body for rest api endpoints that return a list of stdconversation.StdConversation.
//
// The getconversations package should only be use to return response to the api caller.
package getconversations

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
)

// A Response contains response that should be return to the api caller.
type Response struct {
	// Conversation store an slice of stdconversation.StdConversation struct define in package stdconversation.
	Conversations []stdconversation.StdConversation `json:"conversations"`
}

// A filter contains api request's query string parameters for filtering return result.
//
// *** only 1 field can exists at the same time. ***
type Filter struct {
	ParticipantsUsername string `json:"with_participants_username"` // filter conversations by participants's username
	Message              string `json:"with_message"`               // filter covnersations by messages that its contain
}
