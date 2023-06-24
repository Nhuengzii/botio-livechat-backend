// Package getpage define response body for rest api endpoints that return page's specific information
//
// The getpage package should only be use to return response to the api caller.
package getpage

// A Response contains response body of the api request in case the call's a success.
type Response struct {
	UnreadConversations int64 `json:"unreadConversations"` // number of conversations that is not yet read by page admin
	AllConversations    int64 `json:"allConversations"`    // number of page's total conversations
}
