package getall

import "github.com/Nhuengzii/botio-livechat-backend/livechat/shops"

type Response struct {
	Statuses []Status `json:"statuses"`
}

type Status struct {
	Platform            shops.Platform `json:"platform"`
	UnreadConversations int64          `json:"unreadConversations"`
	AllConversations    int64          `json:"allConversations"`
}
