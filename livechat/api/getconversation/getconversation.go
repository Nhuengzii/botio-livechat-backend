// Package getconversation define response body for rest api endpoints that return a single stdconversation.StdConversation.
//
// The getconversation package should only be use to return response to the api caller.
package getconversation

import "github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"

// A Response contains response body of the api request in case the call's a success.
type Response struct {
	// Conversation store a pointer to StdConversation struct define in package stdconversation.
	Conversation *stdconversation.StdConversation `json:"conversation"`
}
