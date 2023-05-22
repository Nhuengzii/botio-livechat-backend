package main

type StandardMessage struct {
	ShopID         string             `json:"shopID"`
	Platform       string             `json:"platform"`
	PageID         string             `json:"pageID"`
	ConversationID string             `json:"conversationID"`
	MessageID      string             `json:"messageID"`
	Timestamp      int64              `json:"timestamp"`
	Source         Source             `json:"source"`
	Message        string             `json:"message"`
	Attachments    []AttachmentOutput `json:"attachments"`
	ReplyTo        ReplyMessage       `json:"replyTo"`
}

type Source struct {
	UserID   string `json:"userID"`
	UserType string `json:"userType"`
}

type AttachmentOutput struct {
	AttachmentType string            `json:"type"`
	Payload        PayloadTypeOutput `json:"payload"`
}

type PayloadTypeOutput struct {
	Src string `json:"src"`
}
