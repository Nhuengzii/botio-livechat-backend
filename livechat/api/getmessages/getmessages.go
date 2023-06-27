// Package getmessages define response body and filter query string parameter for rest api endpoints that return a list of stdmessage.StdMessage.
//
// The getmessages package should only be use to read query string parameter and return response to the api caller.
package getmessages

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

// A Response contains response body of the api request in case the call's a success.
type Response struct {
	// Messages store a slice of stdmessage.StdMessage struct define in package stdmessage.
	Messages []stdmessage.StdMessage `json:"messages"`
}

// A Filter contains api request's query string parameters for filtering return result.
type Filter struct {
	Message string `json:"with_message"` // filter messages that contains the text message
}
