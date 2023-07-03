// Package getall contains structs for the response of cmd/lambda/all/get_all handler
// (the "get shop's all platforms' statuses" api)
package getall

import "github.com/Nhuengzii/botio-livechat-backend/livechat/shops"

// Response is the struct for marshalling the response of the api
type Response struct {
	Statuses []Status `json:"statuses"`
}

// Status contains conversation unread and all conversations counts of a platform
type Status struct {
	Platform            shops.Platform `json:"platform"`
	UnreadConversations int64          `json:"unreadConversations"`
	AllConversations    int64          `json:"allConversations"`
}
