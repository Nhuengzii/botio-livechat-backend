// Package getconversations define response body and filter query string parameter for rest api endpoints that return a list of stdconversation.StdConversation.
//
// The getconversations package should only be use to read query string parameter and return response to the api caller.
package getconversations

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"
)

// A Response contains response body of the api request in case the call's a success.
type Response struct {
	// Conversations store a slice of stdconversation.StdConversation struct define in package stdconversation.
	Conversations []stdconversation.StdConversation `json:"conversations"`
}

// A Filter contains api request's query string parameters for filtering return result.
//
// *** only 1 field can exists at the same time. ***
type Filter struct {
	ParticipantsUsername string `json:"with_participants_username"` // filter conversations by participants's username
	Message              string `json:"with_message"`               // filter conversations by messages that its contain
}
