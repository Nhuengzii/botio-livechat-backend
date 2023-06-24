// Package getconversation define response body for rest api endpoints that return a single stdconversation.StdConversation.
//
// The getconversation package should only be use to return response to the api caller.
package getconversation

import "github.com/Nhuengzii/botio-livechat-backend/livechat/stdconversation"

// A Response contains response body that should be return to the api caller.
type Response struct {
	// Conversation store a pointer to StdConversation struct define in package stdconversation.
	Conversation *stdconversation.StdConversation `json:"conversation"`
}
