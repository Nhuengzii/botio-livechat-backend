// Package getpage define response body for rest api endpoints that return page's specific information
//
// The getpage package should only be use to return response to the api caller.
package getpage

// A Response contains response body that should be return to the api caller.
type Response struct {
	UnreadConversations int64 `json:"unreadConversations"` // number of conversations that is not yet read by page admin
	AllConversations    int64 `json:"allConversations"`    // number of page's total conversations
}
