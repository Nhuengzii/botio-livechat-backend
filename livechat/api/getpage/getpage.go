package getpage

type Response struct {
	UnreadConversations int64 `json:"unreadConversations"`
	AllConversations    int64 `json:"allConversations"`
}
