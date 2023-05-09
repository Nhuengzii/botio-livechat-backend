package main

type StandardMessage struct {
	ShopID         string       `json:"shopID"`
	PageID         string       `json:"pageID"`
	ConversationID string       `json:"conversationID"`
	MessageID      string       `json:"messageID"`
	Timestamp      int64        `json:"timestamp"`
	Source         Source       `json:"source"`
	Message        string       `json:"message"`
	Attachments    []Attachment `json:"attachments"`
	ReplyTo        ReplyMessage `json:"replyTo,omitempty"`
}

type Source struct {
	UserID   string `json:"userID"`
	UserType string `json:"userType"`
}
