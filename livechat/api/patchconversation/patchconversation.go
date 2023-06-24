// Package getpage define request body for rest api endpoints that return page's specific information
//
// The patchconversation package should only be use to read api request.
package patchconversation

// A Request contains request body of the api request.
type Request struct {
	Unread *int `json:"unread"` // number of unread messages in the conversation
}
